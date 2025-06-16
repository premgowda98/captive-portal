package main

/*
#include "platform/platform.h"
*/
import (
	"C"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/premgowda/cgo-impl/iall-ntfr-impl/platform"
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

func behindCaptivePortal() (bool, string) {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent following redirects so we can inspect the redirect URL
			return http.ErrUseLastResponse
		},
	}

	url := "http://clients3.google.com/generate_204"

	fmt.Println("Checking captive portal at:", url)

	timeoutCount := 0
	for i := 1; i <= 3; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("Attempt %d: %v\n", i, err)
			if isTimeoutError(err) {
				timeoutCount++
			}
			time.Sleep(3 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			redirectURL := resp.Header.Get("Location")
			if redirectURL == "" {
				redirectURL = url // fallback if location not present
			}
			fmt.Printf("Redirect detected (status: %d, location: %s)\n", resp.StatusCode, redirectURL)
			return true, redirectURL
		}

		fmt.Println("Captive portal not found")
		return false, url
	}

	if timeoutCount == 3 {
		fmt.Println("All attempts timed out — assuming captive portal")
		return true, url
	}

	fmt.Println("Captive portal not found")
	return false, url
}

func isTimeoutError(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	fmt.Println("Opening browser to:", url)

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return
	}
	_ = cmd.Start()
}
