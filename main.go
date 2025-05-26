package main

import (
	"time"

	"github.com/premgowda98/captive-portal/captive"
)

const DetectionOn = true

func main() {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				if DetectionOn {
					captive.MonitorCaptivePortal()
				}
			}()
		}
	}
}
