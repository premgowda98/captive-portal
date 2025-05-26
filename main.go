package main

import (
	"log"
	"time"

	"github.com/premgowda98/captive-portal/captive"
)

const DetectionOn = true

func main() {
	log.Println("Starting Captive Portal Monitor...")

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
