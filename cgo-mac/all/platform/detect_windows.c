// +build windows

#define _WIN32_WINNT 0x0600 // Vista or later
#define WIN32_LEAN_AND_MEAN

#include <windows.h>
#include <winsock2.h>
#include <ws2tcpip.h>
#include <iphlpapi.h>
#include <netioapi.h>  // This is crucial for NotifyIpInterfaceChange()
#include <stdio.h> 

#pragma comment(lib, "iphlpapi.lib")
#pragma comment(lib, "ws2_32.lib")

extern void networkChangedCallback(); // Defined in Go

// Callback invoked when IP interface changes
VOID CALLBACK onNetworkChange(
    PVOID CallerContext,
    PMIB_IPINTERFACE_ROW Row,
    MIB_NOTIFICATION_TYPE NotificationType
) {
    networkChangedCallback(); // Call into Go
}

void startListening() {
    HANDLE handle;
    DWORD result;

    result = NotifyIpInterfaceChange(
        AF_UNSPEC,
        (PIPINTERFACE_CHANGE_CALLBACK)onNetworkChange,
        NULL,
        FALSE,
        &handle
    );

    if (result != NO_ERROR) {
        fprintf(stderr, "NotifyIpInterfaceChange failed: %lu\n", result);
        return;
    }

    // Keep the thread alive to receive callbacks
    while (1) {
        Sleep(INFINITE);
    }
}
