# Captive Portal Detector

A cross-platform network monitoring tool that detects captive portals by monitoring network changes and checking internet connectivity.

## Overview

This application monitors network interface changes (WiFi connections, IP address changes, etc.) and automatically checks for captive portals when network changes are detected. It works on macOS, Windows, and Linux using native system APIs for efficient, event-driven network monitoring.

## Features

- **Cross-platform support**: macOS, Windows, and Linux
- **Event-driven monitoring**: Low CPU usage, immediate response to network changes
- **Captive portal detection**: Automatically checks for captive portals when network changes
- **Native system integration**: Uses platform-specific APIs for optimal performance

## Prerequisites

- **Go 1.16+** installed
- **C compiler**:
  - **macOS**: Xcode Command Line Tools (`xcode-select --install`)
  - **Windows**: MinGW-w64 or Microsoft Visual C++
  - **Linux**: GCC (`sudo apt-get install build-essential` on Ubuntu/Debian)
- **CGO enabled** (default in most Go installations)

## Installation & Usage

### Clone the repository
```bash
git clone https://github.com/premgowda/captive-portal.git
cd captive-portal/cgo-mac/all
```

### Run on macOS
```bash
go run main.go
```

### Run on Windows
```cmd
# Make sure CGO is enabled
set CGO_ENABLED=1
go run main.go
```

### Run on Linux
```bash
go run main.go
```

### Build executable
```bash
# Build for current platform
go build -o captive-portal main.go

# Run the executable
./captive-portal        # macOS/Linux
captive-portal.exe      # Windows
```

## How it Works

### Network Monitoring
The application uses platform-specific APIs for efficient network monitoring:

- **macOS**: System Configuration framework (`SCDynamicStore`)
- **Windows**: IP Helper API (`NotifyAddrChange`)
- **Linux**: Netlink sockets (`NETLINK_ROUTE`)

### Captive Portal Detection
1. Monitors network interface changes
2. When a change is detected, makes HTTP requests to well-known endpoints
3. Checks if responses are redirected (indicating a captive portal)
4. Logs results and can trigger custom actions

### Architecture
```
main.go
├── platform/
│   ├── detect_darwin.m     # macOS network monitoring
│   ├── detect_windows.c    # Windows network monitoring  
│   ├── detect_linux.c      # Linux network monitoring
│   ├── network_darwin.go   # macOS Go wrapper
│   ├── network_windows.go  # Windows Go wrapper
│   └── network_linux.go    # Linux Go wrapper
└── captive.go              # Captive portal detection logic
```

## Expected Output

```
Started monitoring network changes on [Platform]...
Network configuration changed
Checking for captive portal...
No captive portal detected - Internet access is available
```

When a captive portal is detected:
```
Network configuration changed
Checking for captive portal...
Captive portal detected! Please open your browser to authenticate.
```

## Troubleshooting

### Windows Issues
- **Linker errors**: Make sure MinGW-w64 is properly installed and in PATH
- **CGO disabled**: Run `set CGO_ENABLED=1` before building
- **Missing libraries**: Ensure you have the Windows SDK installed

### macOS Issues
- **Xcode tools**: Install with `xcode-select --install`
- **Permissions**: Some network monitoring may require administrator privileges

### Linux Issues
- **Permissions**: May need to run with `sudo` for netlink socket access
- **Missing headers**: Install development packages (`build-essential`, `libc6-dev`)

### General Issues
- **CGO not enabled**: Ensure `echo $CGO_ENABLED` returns `1`
- **Go version**: Make sure you're using Go 1.16 or later
- **Network access**: Ensure the machine has internet connectivity for testing

## Testing

To test the captive portal detection:

1. **Run the application**
2. **Trigger network changes**:
   - Connect/disconnect from WiFi
   - Switch between different networks
   - Enable/disable network adapters
   - Connect/disconnect Ethernet cable

The application should detect these changes and check for captive portals automatically.

## Use Cases

- **Corporate networks**: Detect when connecting to networks with authentication portals
- **Public WiFi**: Automatically detect hotel/cafe captive portals
- **Network diagnostics**: Monitor network connectivity changes
- **Automated workflows**: Trigger custom actions when network changes occur

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test on your target platform(s)
5. Submit a pull request

## License

This project is open source. Please check the license file for details.

## Platform-Specific Notes

### macOS
- Uses Objective-C for System Configuration framework integration
- Requires Cocoa framework for proper event loop handling

### Windows  
- Uses Windows IP Helper API for network change notifications
- Supports Windows Vista and later versions

### Linux
- Uses netlink sockets for kernel network event notifications  
- Should work on most Linux distributions with kernel 2.6+