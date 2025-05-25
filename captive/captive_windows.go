//go:build windows

package captive

import (
	"fmt"
	"os/exec"
)

func OpenCaptivePortalLogin() error {
	url := "http://clients3.google.com/generate_204"
	
	cmd := exec.Command("cmd", "/c", "start", url)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser on Windows: %v", err)
	}
	
	return nil
}