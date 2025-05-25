package main

import (
	"fmt"
	"log"

	"github.com/premgowda98/captive-portal/captive"
)

func main() {
	log.Println("Starting Captive Portal Monitor...")

	state := captive.GetCurrentState()

	displayStatus(state)

	go func() {
		captive.MonitorCaptivePortal()
	}()

	select {}
}

func displayStatus(state captive.NetworkState) {
	fmt.Printf("\n=== Network Status ===\n")
	fmt.Printf("Outbound IP: %s\n", state.OutboundIP)
	fmt.Printf("Internet Connected: %t\n", state.IsConnected)
	fmt.Printf("Captive Portal: %t\n", state.HasCaptivePortal)
	fmt.Printf("Browser Opened: %t\n", state.BrowserOpened)
	fmt.Printf("Last Check: %s\n", state.LastCheck.Format("15:04:05"))
	fmt.Println("====================")
}
