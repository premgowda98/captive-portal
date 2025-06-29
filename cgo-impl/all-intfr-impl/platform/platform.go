package platform

/*
#include "platform/platform.h"
*/
import (
	"C"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)
import "strings"

var (
	CaptiveCheckURLs = map[string][]string{
		"darwin": {
			"http://captive.apple.com/hotspot-detect.html",
			"http://clients3.google.com/generate_204",
		},
		"linux": {
			"http://clients3.google.com/generate_204",
		},
		"windows": {
			"http://clients3.google.com/generate_204",
		},
	}
)

const (
	retryCount  = 5
	retryDelay  = 1 * time.Second
	httpTimeout = 3 * time.Second
)

//export networkChangedCallback
func networkChangedCallback() {
	slog.Info("Network change detected\n")

	if isCaptive, captiveUrl := behindCaptivePortal(); isCaptive {
		slog.Info("Captive portal detected. Opening browser...\n")
		openBrowser(captiveUrl)
	}
}

func isNetworkReachable() bool {
	conn, err := net.DialTimeout("udp", "8.8.8.8:53", 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

// behindCaptivePortal checks if the device is behind a captive portal by making an HTTP request to a known URL.
func behindCaptivePortal() (bool, string) {

	client := &http.Client{
		Timeout: httpTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects so we can inspect the portal behavior
			return http.ErrUseLastResponse
		},
	}

	urls := CaptiveCheckURLs[runtime.GOOS]
	if len(urls) == 0 {
		slog.Warn(fmt.Sprintf("Captive: No check URLs defined for platform: %s\n", runtime.GOOS))
		return false, ""
	}

	for attempt := 1; attempt <= retryCount; attempt++ {
		if !isNetworkReachable() {
			slog.Info(fmt.Sprintf("Captive: Attempt %d: Network not reachable, waiting...\n", attempt))
			time.Sleep(retryDelay)
			continue
		}

		url := urls[(attempt-1)%len(urls)]
		slog.Info(fmt.Sprintf("Captive: Attempt %d: Sending request to %s\n", attempt, url))

		resp, err := client.Get(url)
		if err != nil {
			if isTimeoutError(err) {
				slog.Info(fmt.Sprintf("Captive: Attempt %d: Request timed out, likely captive portal or slow net\n", attempt))
				time.Sleep(retryDelay)
				continue
			}
			slog.Warn(fmt.Sprintf("Captive: Attempt %d: Unexpected error: %v\n", attempt, err))
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		slog.Info(fmt.Sprintf("Captive: Response Status=%d\n", resp.StatusCode))
		headersMap := make(map[string][]string)
		for key, vals := range resp.Header {
			headersMap[key] = vals
		}
		slog.Info(fmt.Sprintf("Captive: Response Headers: %+v\n", headersMap))

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Warn(fmt.Sprintf("Captive: Failed to read response body: %v\n", err))
		} else {
			slog.Info(fmt.Sprintf("Captive: Response Body:\n%s\n", string(bodyBytes)))
		}

		// captive portal redirect
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			redirectURL := resp.Header.Get("Location")
			if redirectURL == "" {
				time.Sleep(retryDelay)
				continue
			}
			slog.Info(fmt.Sprintf(
				"Captive: Redirect detected (status %d): Captive portal likely at %s\n",
				resp.StatusCode,
				redirectURL,
			))
			return true, redirectURL
		}

		if resp.StatusCode == http.StatusOK {
			// In some cases even if we get 200 OK, the body might indicate a captive portal
			// Here captive portal hijacks the native os by sending javascript based redirects instead of HTTP redirects
			// captive.apple.com is a standard where it will give success if the captive portal is not detected
			// clients3.google.com is a standard where it will give 204 No Content if the captive portal is not detected

			if strings.Contains(url, "captive.apple.com") {
				if !strings.Contains(strings.ToLower(string(bodyBytes)), "success") {
					slog.Info("Captive: captive.apple.com body does not contain 'Success', captive portal detected")
					return true, url
				}
			}

			if strings.Contains(url, "clients3.google.com") {
				if len(bodyBytes) != 0 {
					slog.Info("Captive: clients3.google.com body is not empty, captive portal detected")
					return true, url
				}
			}
		}

		slog.Info("Captive: No captive portal detected\n")

		return false, url
	}

	slog.Info("Captive: Max attempts reached, no captive portal detected\n")
	return false, ""
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	slog.Info(fmt.Sprintf("Opening browser to: %s\n", url))

	chromePath := detectChromePath()

	if chromePath != "" {
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", "-a", "Google Chrome", url)
		default:
			cmd = exec.Command(chromePath, url)
		}
	} else {
		// Fallback to system default browser
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		default:
			slog.Warn("Unsupported platform\n")
			return
		}
	}

	if err := cmd.Start(); err != nil {
		slog.Warn(fmt.Sprintf("Failed to open browser: %v\n", err))
	}
}

func detectChromePath() string {
	var paths []string

	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	case "linux":
		paths = []string{
			"google-chrome", "chrome", "chromium", "chromium-browser",
		}
	case "windows":
		paths = []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		}
	default:
		return ""
	}

	for _, path := range paths {
		switch runtime.GOOS {
		case "linux":
			// Check if Chrome-like command exists in PATH
			if fullPath, err := exec.LookPath(path); err == nil {
				return fullPath
			}
		case "windows":
			// Check is Chrome executable exists at the specified path
			_, err := os.Stat(path)
			if err == nil || !os.IsNotExist(err) {
				return path
			}

		default: // macOS or others
			// Use 'test -f' to check if file exists
			checkCmd := exec.Command("test", "-f", path)
			if err := checkCmd.Run(); err == nil {
				return path
			}
		}
	}

	return ""
}
