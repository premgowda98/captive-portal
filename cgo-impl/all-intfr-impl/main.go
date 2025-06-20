package main

/*
#include "platform/platform.h"
*/
import (
	"C"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/premgowda/cgo-impl/iall-ntfr-impl/platform"
)
import "os"

const (
	checkURL    = "http://clients3.google.com/generate_204"
	retryCount  = 5
	retryDelay  = 1 * time.Second
	httpTimeout = 3 * time.Second
)

//export networkChangedCallback
func networkChangedCallback() {
	fmt.Println("Network change detected")

	if isCaptive, captiveUrl := behindCaptivePortal(); isCaptive {
		fmt.Println("Captive portal detected. Opening browser...")
		openBrowser(captiveUrl)
	}
}

func main() {
	fmt.Println("Starting network monitor for", runtime.GOOS)
	// to test captive at first run
	networkChangedCallback()
	go platform.StartNetworkMonitor()

	// Keep app running
	select {}
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

func behindCaptivePortal() (bool, string) {

	client := &http.Client{
		Timeout: httpTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects so we can inspect the portal behavior
			return http.ErrUseLastResponse
		},
	}

	fmt.Printf("Checking for captive portal at: %s", checkURL)

	for attempt := 1; attempt <= retryCount; attempt++ {
		if !isNetworkReachable() {
			fmt.Printf("Attempt %d: Network not reachable, waiting...", attempt)
			time.Sleep(retryDelay)
			continue
		}

		resp, err := client.Get(checkURL)
		if err != nil {
			if isTimeoutError(err) {
				fmt.Printf("Attempt %d: Request timed out, likely captive portal or slow net", attempt)
				time.Sleep(retryDelay)
				continue
			}
			fmt.Printf("Attempt %d: Unexpected error: %v", attempt, err)
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		// captive portal redirect
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			redirectURL := resp.Header.Get("Location")
			if redirectURL == "" {
				redirectURL = checkURL
			}
			fmt.Printf("Redirect detected (status %d): Captive portal likely at %s", resp.StatusCode, redirectURL)
			return true, redirectURL
		}

		fmt.Println("No captive portal detected")
		return false, checkURL
	}

	fmt.Println("Max attempts reached — assuming captive portal or unreachable network")
	return true, checkURL
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	fmt.Println("Opening browser to:", url)

	chromePath := detectChromePath()

	if chromePath != "" {
		cmd = exec.Command(chromePath, url)
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
			fmt.Println("Unsupported platform")
			return
		}
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to open browser:", err)
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
