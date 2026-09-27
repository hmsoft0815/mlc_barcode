//go:build ios

package main

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework UIKit -framework WebKit

#import <UIKit/UIKit.h>
#import <WebKit/WebKit.h>

// Wails sets no WKUIDelegate on its web view, so WebKit asks "allow
// camera?" for the app's own page on every launch. This delegate answers
// for our own content only (the wails:// app page); the system permission
// (NSCameraUsageDescription) is untouched and still asked once by iOS.
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
    decisionHandler(own ? WKPermissionDecisionGrant : WKPermissionDecisionPrompt);
}
@end

static id mlcGrant; // UIDelegate is weak: keep the delegate alive

static WKWebView *mlcFindWebView(UIView *v) {
    if ([v isKindOfClass:[WKWebView class]]) return (WKWebView *)v;
    for (UIView *s in v.subviews) {
        WKWebView *w = mlcFindWebView(s);
        if (w) return w;
    }
    return nil;
}

// The web view appears some moments after launch: look for it on the main
// queue, twice a second, for up to 20 seconds.
static void mlcTryInstall(int attempt) {
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(0.5 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
        if (@available(iOS 15.0, *)) {
            for (UIScene *scene in UIApplication.sharedApplication.connectedScenes) {
                if (![scene isKindOfClass:[UIWindowScene class]]) continue;
                for (UIWindow *win in ((UIWindowScene *)scene).windows) {
                    WKWebView *wv = mlcFindWebView(win);
                    if (wv) {
                        if (!mlcGrant) mlcGrant = [MLCMediaGrant new];
                        if (!wv.UIDelegate) wv.UIDelegate = mlcGrant;
                        return;
                    }
                }
            }
            if (attempt < 40) mlcTryInstall(attempt + 1);
        }
    });
}

static void mlcInstallMediaGrant(void) { mlcTryInstall(0); }
*/
import "C"

// installCameraGrant lets the app's own page use the camera without
// WebKit asking again on every launch (see the Objective-C above).
func installCameraGrant() { C.mlcInstallMediaGrant() }
