//go:build windows

package platform

/*
#cgo LDFLAGS: -liphlpapi -lws2_32
*/
import "C"
