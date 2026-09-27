//go:build ios

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework UIKit -framework WebKit

#import <UIKit/UIKit.h>
#import <WebKit/WebKit.h>
#import <objc/runtime.h>

// Wails sets no WKUIDelegate on its web view, so WebKit asks "allow
// camera?" for the app's own page on every launch. This delegate answers
// for our own content only (the wails:// app page); the system permission
// (NSCameraUsageDescription) is untouched and still asked once by iOS.
//
// It must be in place before the page asks: the scan view starts the camera
// as soon as it shows. So it is set when the web view is created — the
// category below extends -[WKWebView initWithFrame:configuration:] at load
// time — not looked up afterwards.

// Diagnostics go to the system log and to stderr (seen with
// `devicectl device process launch --console`).
#define MLCLog(fmt, ...) do { NSString *m = [NSString stringWithFormat:fmt, ##__VA_ARGS__]; NSLog(@"%@", m); fprintf(stderr, "%s\n", m.UTF8String); } while (0)

API_AVAILABLE(ios(15.0))
@interface MLCMediaGrant : NSObject <WKUIDelegate>
@end

@implementation MLCMediaGrant
- (void)webView:(WKWebView *)webView
    requestMediaCapturePermissionForOrigin:(WKSecurityOrigin *)origin
                          initiatedByFrame:(WKFrameInfo *)frame
                                      type:(WKMediaCaptureType)type
                           decisionHandler:(void (^)(WKPermissionDecision))decisionHandler {
    BOOL own = [origin.protocol isEqualToString:@"wails"] || [origin.host isEqualToString:@"localhost"];
    MLCLog(@"[mlc-camera] media request from %@://%@ type %ld -> %@", origin.protocol, origin.host, (long)type, own ? @"grant" : @"prompt");
    decisionHandler(own ? WKPermissionDecisionGrant : WKPermissionDecisionPrompt);
}
@end

static id mlcGrant; // UIDelegate is weak: keep the delegate alive

@implementation WKWebView (MLCMediaGrant)
+ (void)load {
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        Method original = class_getInstanceMethod(self, @selector(initWithFrame:configuration:));
        Method extended = class_getInstanceMethod(self, @selector(mlc_initWithFrame:configuration:));
        method_exchangeImplementations(original, extended);
    });
}

// After the exchange this calls the original initializer.
- (instancetype)mlc_initWithFrame:(CGRect)frame configuration:(WKWebViewConfiguration *)configuration {
    WKWebView *wv = [self mlc_initWithFrame:frame configuration:configuration];
    if (@available(iOS 15.0, *)) {
        if (wv && !wv.UIDelegate) {
            if (!mlcGrant) mlcGrant = [MLCMediaGrant new];
            wv.UIDelegate = mlcGrant;
            MLCLog(@"[mlc-camera] media delegate set on new web view %@", NSStringFromClass([wv class]));
        }
    }
    return wv;
}
@end

static void mlcInstallMediaGrant(void) {}
*/
import "C"

// installCameraGrant: the Objective-C above does its work when the app
// image loads (+load); calling it keeps the code linked in.
func installCameraGrant() { C.mlcInstallMediaGrant() }
