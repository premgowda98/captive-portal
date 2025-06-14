//go:build linux

package captive

import (
	"fmt"
	"os/exec"
)

func OpenCaptivePortalLogin(redirectUrl string) error {
	browsers := []string{"xdg-open", "firefox", "google-chrome", "chromium", "mozilla"}

	if redirectUrl == "" {
		redirectUrl = "http://clients3.google.com/generate_204"
	}

	fmt.Println("Opening in", redirectUrl)

	for _, browser := range browsers {
		cmd := exec.Command(browser, redirectUrl)
		if err := cmd.Start(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to open browser on Linux")
}
