package main

import (
	"log"

	"github.com/premgowda98/captive-portal/captive"
)

func main() {
	log.Println("Starting Captive Portal Monitor...")

	go func() {
		captive.MonitorCaptivePortal()
	}()

	select {}
}
