//go:build linux

package captive

import (
	"fmt"
	"os/exec"
)

func OpenCaptivePortalLogin() error {
	browsers := []string{"xdg-open", "firefox", "google-chrome", "chromium", "mozilla"}
	url := "http://clients3.google.com/generate_204"

	for _, browser := range browsers {
		cmd := exec.Command(browser, url)
		if err := cmd.Start(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to open browser on Linux")
}
