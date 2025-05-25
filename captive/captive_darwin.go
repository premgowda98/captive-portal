//go:build darwin

package captive

import (
	"fmt"
	"os/exec"
)

// OpenCaptivePortalLogin opens the captive portal login page in the default browser on macOS
func OpenCaptivePortalLogin() error {
	url := "http://clients3.google.com/generate_204"

	// Use the 'open' command on macOS to open the URL in the default browser
	cmd := exec.Command("open", url)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser on macOS: %v", err)
	}

	return nil
}
