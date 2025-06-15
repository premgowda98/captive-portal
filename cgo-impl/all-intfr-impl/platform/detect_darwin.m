#import <Foundation/Foundation.h>
#import <SystemConfiguration/SystemConfiguration.h>

extern void networkChangedCallback(const char *msg); // Go callback

void networkChanged(SCDynamicStoreRef store, CFArrayRef changedKeys, void *info) {
    networkChangedCallback("Network configuration changed");
}

void startMonitoringNetworkChanges() {
    SCDynamicStoreContext context = {0, NULL, NULL, NULL, NULL};
    SCDynamicStoreRef store = SCDynamicStoreCreate(NULL, CFSTR("NetworkMonitor"), networkChanged, &context);

    if (!store) {
        NSLog(@"Failed to create SCDynamicStore");
        return;
    }

    CFStringRef patterns[2] = {
        CFSTR("State:/Network/Interface/.*/IPv4"),
        CFSTR("State:/Network/Global/IPv4")
    };
    CFArrayRef patternArray = CFArrayCreate(NULL, (const void **)patterns, 2, &kCFTypeArrayCallBacks);

    SCDynamicStoreSetNotificationKeys(store, NULL, patternArray);

    CFRunLoopSourceRef runLoopSource = SCDynamicStoreCreateRunLoopSource(NULL, store, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
    CFRelease(runLoopSource);
    CFRelease(patternArray);

    NSLog(@"Started monitoring network changes...");
    CFRunLoopRun();
}