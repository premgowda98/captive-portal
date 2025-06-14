# Captive Portal Detector

A cross-platform Go program that automatically detects captive portals (hotel/cafe WiFi login pages) and opens them in your browser for authentication.

## What it does

- **Monitors internet connectivity** every 5 seconds
- **Detects captive portals** by checking for redirects/blocked traffic
- **Automatically opens login pages** in your default browser
- **Handles network changes** without repeatedly opening browsers
- **Provides detailed DNS debugging** to troubleshoot connectivity issues
- **Works across platforms**: Linux, macOS, and Windows

## How to run

1. **Build the program:**
   ```bash
   go build -o captive-portal
   or
   go run main.go
   ```

2. **Run it:**
   ```bash
   ./captive-portal
   ```

3. **Connect to a captive portal network** (like hotel/cafe WiFi)

4. **The program will:**
   - Detect the captive portal automatically
   - Open the login page in your browser
   - Monitor until you complete authentication
   - Continue monitoring for network changes