//go:build linux

package platform

/*
#cgo LDFLAGS: -lresolv
#include "platform.h"

extern void networkChangedCallback(); // implemented in Go
*/
import "C"

func StartNetworkMonitor() {
	C.startListening()
}
