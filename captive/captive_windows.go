//go:build windows

package captive

import (
	"fmt"
	"os/exec"
)

// OpenCaptivePortalLogin opens the captive portal login page in the default browser on Windows
func OpenCaptivePortalLogin() error {
	url := "http://clients3.google.com/generate_204"
	
	// Use the 'start' command on Windows to open the URL in the default browser
	cmd := exec.Command("cmd", "/c", "start", url)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser on Windows: %v", err)
	}
	
	return nil
}