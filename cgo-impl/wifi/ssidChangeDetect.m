#import <Foundation/Foundation.h>
#import <CoreWLAN/CoreWLAN.h>

extern void ssidChangedCallback(const char *ssid); // defined in Go

@interface SSIDMonitor : NSObject <CWEventDelegate>
@end

@implementation SSIDMonitor

- (void)ssidDidChangeForWiFiClient:(CWWiFiClient *)client {
    NSString *ssid = client.interface.ssid;
    if (ssid) {
        ssidChangedCallback([ssid UTF8String]);
    }
}

@end

static SSIDMonitor *monitor;


void startMonitoringSSID() {
    CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    monitor = [[SSIDMonitor alloc] init]; // global static retains it
    client.delegate = monitor;
    [client startMonitoringEventWithType:CWEventTypeSSIDDidChange error:nil];
    [[NSRunLoop mainRunLoop] run];
}
