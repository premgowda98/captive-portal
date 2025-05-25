package captive

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// NetworkState holds the current network state
type NetworkState struct {
	OutboundIP       string
	LastCheck        time.Time
	IsConnected      bool
	HasCaptivePortal bool
	BrowserOpened    bool  // Track if we've already opened browser for current network
}

var currentState NetworkState

// GetOutboundIP determines the current active outbound IP address
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// CheckInternetConnectivity tests connectivity to Google
func CheckInternetConnectivity() bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://www.google.com")
	if err != nil {
		log.Printf("Internet connectivity check failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// CheckCaptivePortal detects if a captive portal is intercepting requests
func CheckCaptivePortal() bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects - we want to detect them
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get("http://clients3.google.com/generate_204")
	if err != nil {
		log.Printf("Captive portal check failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Status 204 (No Content) means no captive portal
	// Any other status or redirect indicates captive portal
	return resp.StatusCode != 204
}

// HasNetworkChanged checks if the IP address has changed
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

// UpdateNetworkState updates the current network state
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

// MonitorCaptivePortal is the main monitoring function
func MonitorCaptivePortal() {
	log.Println("Starting captive portal monitoring...")

	// Initial network state setup
	if err := UpdateNetworkState(); err != nil {
		log.Printf("Failed to initialize network state: %v", err)
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Initial check
	performConnectivityCheck()

	for {
		select {
		case <-ticker.C:
			// Check if network has changed
			if HasNetworkChanged() {
				log.Println("Network change detected, updating state...")
				if err := UpdateNetworkState(); err != nil {
					log.Printf("Failed to update network state: %v", err)
					continue
				}
				// Reset browser opened flag when network changes
				currentState.BrowserOpened = false
				log.Println("Browser opened flag reset due to network change")
			}

			performConnectivityCheck()
		}
	}
}

// performConnectivityCheck runs the connectivity and captive portal checks
func performConnectivityCheck() {
	log.Println("Checking connectivity...")

	// Check internet connectivity
	isConnected := CheckInternetConnectivity()
	currentState.IsConnected = isConnected

	if isConnected {
		log.Println("Internet connectivity confirmed - full access available")
		currentState.HasCaptivePortal = false
		// Reset browser opened flag when internet is working
		currentState.BrowserOpened = false
		return
	}

	log.Println("No internet connectivity detected - checking for captive portal...")

	// Check for captive portal only when internet is not accessible
	hasCaptivePortal := CheckCaptivePortal()
	currentState.HasCaptivePortal = hasCaptivePortal

	if hasCaptivePortal {
		if !currentState.BrowserOpened {
			log.Println("Captive portal detected! Opening login page...")
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

// GetCurrentState returns the current network state
func GetCurrentState() NetworkState {
	return currentState
}
