//go:build darwin

package platform

/*
#cgo LDFLAGS: -framework Foundation -framework SystemConfiguration
void startMonitoringNetworkChanges();
*/
import "C"

func StartNetworkMonitor() {
	C.startMonitoringNetworkChanges()
}
