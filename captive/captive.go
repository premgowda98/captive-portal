package captive

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	CaptivePortalURL = "http://clients3.google.com/generate_204"
	DialPingAddress  = "8.8.8.8:80"
	InternetCheckURL = "https://www.google.com"
)

type NetworkState struct {
	OutboundIP       string
	LastCheck        time.Time
	IsConnected      bool
	HasCaptivePortal bool
	BrowserOpened    bool
}

var currentState NetworkState

func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", DialPingAddress)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func CheckInternetConnectivity() bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(InternetCheckURL)
	if err != nil {
		log.Printf("Internet connectivity check failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func CheckCaptivePortal() bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects - we want to detect them
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(CaptivePortalURL)
	if err != nil {
		log.Printf("Captive portal check failed: %v", err)
		// If the request fails entirely, it could be a captive portal blocking traffic
		// We'll return true (captive portal likely) in this case
		return true
	}
	defer resp.Body.Close()

	// Status 204 (No Content) means no captive portal
	// Any other status (redirect, blocked page, etc.) indicates captive portal
	isCaptivePortal := resp.StatusCode != 204

	if isCaptivePortal {
		log.Printf("Captive portal detected - received status %d instead of 204", resp.StatusCode)
	}

	return isCaptivePortal
}

func HasNetworkChanged() bool {
	newIP, err := GetOutboundIP()
	if err != nil {
		log.Printf("Failed to get outbound IP: %v", err)
		return false
	}

	changed := currentState.OutboundIP != newIP

	if changed {
		log.Printf("Network change detected: IP %s->%s",
			currentState.OutboundIP, newIP)
	}

	return changed
}

func UpdateNetworkState() error {
	ip, err := GetOutboundIP()
	if err != nil {
		return fmt.Errorf("failed to get outbound IP: %v", err)
	}

	currentState.OutboundIP = ip
	currentState.LastCheck = time.Now()

	log.Printf("Network state updated: IP=%s", ip)
	return nil
}

func MonitorCaptivePortal() {
	log.Println("Starting captive portal monitoring...")

	if err := UpdateNetworkState(); err != nil {
		log.Printf("Failed to initialize network state: %v", err)
		return
	}

	if HasNetworkChanged() {
		log.Println("Network change detected, updating state...")
		if err := UpdateNetworkState(); err != nil {
			log.Printf("Failed to update network state: %v", err)
		}
		currentState.BrowserOpened = false
		log.Println("Browser opened flag reset due to network change")
	}

	performConnectivityCheck()
}

func performConnectivityCheck() {
	log.Println("Checking connectivity...")

	isConnected := CheckInternetConnectivity()
	currentState.IsConnected = isConnected

	if isConnected {
		log.Println("Internet connectivity confirmed - full access available")
		currentState.HasCaptivePortal = false
		currentState.BrowserOpened = false
		return
	}

	log.Println("No internet connectivity detected - checking for captive portal...")

	hasCaptivePortal := CheckCaptivePortal()
	currentState.HasCaptivePortal = hasCaptivePortal

	if hasCaptivePortal {
		if !currentState.BrowserOpened {
			log.Println("Captive portal detected Opening login page...")
			if err := OpenCaptivePortalLogin(); err != nil {
				log.Printf("Failed to open captive portal login: %v", err)
			} else {
				currentState.BrowserOpened = true
				log.Println("Browser opened - will not open again until network changes")
			}
		} else {
			log.Println("Captive portal still detected - browser already opened, waiting for user authentication")
		}
	} else {
		log.Println("No captive portal detected - network issue or no connectivity")
	}
}
