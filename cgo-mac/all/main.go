package main

/*
#cgo LDFLAGS: -framework Foundation -framework SystemConfiguration
void startMonitoringNetworkChanges();
*/
import "C"
import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

//export networkChangedCallback
func networkChangedCallback(msg *C.char) {
	fmt.Println("[Go] Network change detected:", C.GoString(msg))
	checkCaptivePortal()
}

func checkCaptivePortal() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://captive.apple.com")
	if err != nil {
		fmt.Println("[Go] Error checking captive portal:", err)
		return
	}
	defer resp.Body.Close()

	// Captive.apple.com returns a page with the title "Success"
	// If redirected or blocked, this likely means a captive portal is present
	if resp.StatusCode != http.StatusOK {
		fmt.Println("[Go] Captive portal suspected: non-200 response")
		openBrowser("http://captive.apple.com")
		return
	}

	// Try to detect if it was redirected
	finalURL := resp.Request.URL.String()
	if finalURL != "http://captive.apple.com" {
		fmt.Printf("[Go] Redirected to: %s → Possible captive portal\n", finalURL)
		openBrowser("http://captive.apple.com")
		return
	}

	fmt.Println("[Go] No captive portal detected.")
}

func openBrowser(url string) {
	fmt.Println("[Go] Opening browser to:", url)
	err := exec.Command("open", url).Start() // macOS-specific
	if err != nil {
		fmt.Println("[Go] Failed to open browser:", err)
	}
}

func main() {
	fmt.Println("[Go] Starting network change monitor...")
	checkCaptivePortal()

	go C.startMonitoringNetworkChanges()

	// Wait for termination
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("[Go] Shutting down...")
}
