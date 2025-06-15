//go:build windows

package platform

/*
// #cgo windows LDFLAGS: -liphlpapi -lws2_32
#include "platform.h"

extern void networkChangedCallback(); // implemented in Go
*/
import "C"

func StartNetworkMonitor() {
	C.startListening()
}
