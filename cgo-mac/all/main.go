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

	"github.com/premgowda/cgo-mac/all/platform"
)

//export networkChangedCallback
func networkChangedCallback() {
	fmt.Println("🔄 Network change detected")

	if behindCaptivePortal() {
		fmt.Println("🌐 Captive portal detected. Opening browser...")
		openBrowser("http://captive.apple.com")
	}
}

func main() {
	fmt.Println("📡 Starting network monitor for", runtime.GOOS)
	go platform.StartNetworkMonitor()

	// Keep app running
	select {}
}

func behindCaptivePortal() bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var url string
	switch runtime.GOOS {
	case "darwin":
		url = "http://clients3.google.com/generate_204"
	case "linux":
		url = "http://detectportal.firefox.com/canonical.html"
	case "windows":
		url = "http://www.msftconnecttest.com/connecttest.txt"
	default:
		url = "http://clients3.google.com/generate_204"
	}

	fmt.Println("📶 Checking captive portal at:", url)

	timeoutCount := 0
	for i := 1; i <= 3; i++ {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("❌ Attempt %d: %v\n", i, err)
			if isTimeoutError(err) {
				timeoutCount++
			}
			time.Sleep(3 * time.Second)
			continue
		}
		defer resp.Body.Close()

		// Captive portal usually redirects or doesn't return expected status
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			fmt.Printf("🔁 Redirect detected (status: %d)\n", resp.StatusCode)
			return true
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("⚠️ Unexpected status code: %d\n", resp.StatusCode)
			return true
		}

		// No issues found
		return false
	}

	// If all retries timed out, assume we're behind a captive portal
	if timeoutCount == 3 {
		fmt.Println("⏱️ All attempts timed out — assuming captive portal")
		return true
	}

	return false
}

func isTimeoutError(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	_ = cmd.Start()
}
