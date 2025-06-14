package main

/*
#cgo LDFLAGS: -framework CoreWLAN -framework Foundation
void startMonitoringSSID();
*/
import "C"
import "fmt"

//export ssidChangedCallback
func ssidChangedCallback(ssid *C.char) {
    fmt.Println("SSID changed to:", C.GoString(ssid))
}

func main() {
    fmt.Println("Starting SSID monitor...")
    C.startMonitoringSSID()
    fmt.Println("works")
}


