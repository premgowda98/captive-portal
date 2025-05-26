package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/premgowda98/captive-portal/captive"
)

var DetectionOn = true
var mu sync.RWMutex

func main() {
	log.Println("Starting Captive Portal Monitor...")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				client := &http.Client{
					Timeout: 5 * time.Second,
				}

				resp, err := client.Get("http://localhost:8909")
				if resp != nil {
                    defer resp.Body.Close()
                }
				mu.Lock()
				if err != nil || resp.StatusCode != http.StatusOK {
					fmt.Println("setting detection to off")
					DetectionOn = false
				} else {
					DetectionOn = true
				}
				mu.Unlock()

				if DetectionOn {
					fmt.Println("detection is on")
					captive.MonitorCaptivePortal()
				} else {
					fmt.Println("detection is off")
				}
			}()
		}
	}
}
