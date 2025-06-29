package main

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/premgowda/cgo-impl/iall-ntfr-impl/platform"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.Info(fmt.Sprintf("Starting network monitor for %s\n", runtime.GOOS))
	// to test captive at first run
	go platform.StartNetworkMonitor()

	// Keep app running
	select {}
}
