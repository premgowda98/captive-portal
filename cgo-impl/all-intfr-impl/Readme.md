# Captive Portal Detection

This project implements a cross-platform captive portal detection mechanism using Go and C. It utilizes CGO to interface with platform-specific network APIs.

To build and run the code

```bash
make build
make run
```

To run the process, run the following command:

```bash
make run-process
```

Expected Behavior

1. When the process starts it first looks for the captive portal.
2. Then it starts the network monitor through CGO.
3. When a network change is detected, it checks for captive portal again.
4. If a captive portal is detected, it opens the browser to the specified URL.