//go:build darwin

package captive

import (
	"fmt"
	"os/exec"
)

func OpenCaptivePortalLogin(redirectUrl string) error {
	if redirectUrl == "" {
		redirectUrl = "http://clients3.google.com/generate_204"
	}

	cmd := exec.Command("open", redirectUrl)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser on macOS: %v", err)
	}

	return nil
}
