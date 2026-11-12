//go:build darwin && !ios && cgo

// macOS side of photokit: PHPickerViewController for picking, PhotoKit
// for saving. Every request ends in exactly one pkDone(handle, …) call,
// on whatever thread finishes it; Go waits on that.

#import <Cocoa/Cocoa.h>
#import <Photos/Photos.h>
#import <PhotosUI/PhotosUI.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>
#include <stdint.h>
#include "_cgo_export.h"

enum { PK_OK = 0, PK_CANCELLED = 1, PK_ERROR = 2, PK_DENIED = 3, PK_UNSUPPORTED = 4 };

static void done(uintptr_t h, int status, NSString *data, NSString *msg) {
    @autoreleasepool {
        pkDone(h, status, (char *)(data ? data.UTF8String : ""), (char *)(msg ? msg.UTF8String : ""));
    }
}

int pkPickAvailable(void) {
    if (@available(macOS 13.0, *)) return 1;
    return 0;
}

int pkAppRunning(void) {
    return NSApp != nil && [NSApp isRunning];
}

// First registered type that is an image, in the provider's own order —
// that is the original (public.heic, public.jpeg, com.adobe.raw-image …),
// not a converted copy.
static NSString *pkImageType(NSItemProvider *p) {
    for (NSString *t in p.registeredTypeIdentifiers) {
        UTType *u = [UTType typeWithIdentifier:t];
        if (u && [u conformsToType:UTTypeImage]) return t;
    }
    return nil;
}

// A file name in dir that is neither on disk nor already handed out.
static NSString *pkReserve(NSString *dir, NSString *name, NSMutableSet *taken) {
    NSString *base = name.stringByDeletingPathExtension, *ext = name.pathExtension;
    NSString *cand = name;
    for (int i = 2; [taken containsObject:cand] ||
                    [[NSFileManager defaultManager] fileExistsAtPath:[dir stringByAppendingPathComponent:cand]]; i++) {
        cand = ext.length ? [NSString stringWithFormat:@"%@-%d.%@", base, i, ext]
                          : [NSString stringWithFormat:@"%@-%d", base, i];
    }
    [taken addObject:cand];
    return cand;
}

static NSMutableSet *pkLive; // PHPicker's delegate is weak: keep ours alive

API_AVAILABLE(macos(13.0))
@interface PKPick : NSObject <PHPickerViewControllerDelegate>
@property (nonatomic) uintptr_t handle;
@property (nonatomic, copy) NSString *dir;
@property (nonatomic, strong) NSWindow *window;
@property (nonatomic, weak) NSWindow *parent;
@end

@implementation PKPick

- (void)finish:(int)status data:(NSString *)data msg:(NSString *)msg {
    done(self.handle, status, data, msg);
    dispatch_async(dispatch_get_main_queue(), ^{ [pkLive removeObject:self]; });
}

- (void)picker:(PHPickerViewController *)picker didFinishPicking:(NSArray<PHPickerResult *> *)results {
    if (self.parent) [self.parent endSheet:self.window];
    else [self.window orderOut:nil];
    self.window = nil;

    if (results.count == 0) { [self finish:PK_CANCELLED data:nil msg:nil]; return; }

    // Names are fixed here, on the main thread, so parallel loads cannot
    // pick the same one.
    NSUInteger n = results.count;
    NSMutableArray *lines = [NSMutableArray arrayWithCapacity:n];
    NSMutableArray *types = [NSMutableArray arrayWithCapacity:n];
    NSMutableArray *names = [NSMutableArray arrayWithCapacity:n];
    NSMutableSet *taken = [NSMutableSet set];
    for (NSUInteger i = 0; i < n; i++) {
        NSItemProvider *p = results[i].itemProvider;
        NSString *t = pkImageType(p);
        NSString *ext = t ? [UTType typeWithIdentifier:t].preferredFilenameExtension : nil;
        NSString *stem = p.suggestedName.length ? p.suggestedName : [NSString stringWithFormat:@"Photo-%lu", (unsigned long)i + 1];
        NSString *name = ext.length ? [stem stringByAppendingPathExtension:ext] : stem;
        [lines addObject:[NSNull null]];
        [types addObject:t ?: (id)[NSNull null]];
        [names addObject:pkReserve(self.dir, name, taken)];
    }

    __block NSString *firstErr = nil;
    NSObject *lock = [NSObject new];
    dispatch_group_t g = dispatch_group_create();
    for (NSUInteger i = 0; i < n; i++) {
        if (types[i] == [NSNull null]) { firstErr = firstErr ?: @"a chosen item has no image data"; continue; }
        NSString *type = types[i], *name = names[i];
        NSString *dst = [self.dir stringByAppendingPathComponent:name];
        dispatch_group_enter(g);
        // The URL is a temporary copy, valid only inside the block.
        [results[i].itemProvider loadFileRepresentationForTypeIdentifier:type
                                                     completionHandler:^(NSURL *url, NSError *err) {
            NSString *line = nil, *e = nil;
            NSError *ce = nil;
            if (!url) e = err.localizedDescription ?: @"no file";
            else if ([[NSFileManager defaultManager] copyItemAtPath:url.path toPath:dst error:&ce])
                line = [NSString stringWithFormat:@"%@\t%@\t%@", dst, name, type];
            else e = ce.localizedDescription;
            @synchronized (lock) {
                if (line) lines[i] = line;
                else if (!firstErr) firstErr = e;
            }
            dispatch_group_leave(g);
        }];
    }
    dispatch_group_notify(g, dispatch_get_global_queue(QOS_CLASS_USER_INITIATED, 0), ^{
        NSMutableArray *ok = [NSMutableArray array];
        for (id l in lines) if (l != [NSNull null]) [ok addObject:l];
        if (ok.count == 0) [self finish:PK_ERROR data:nil msg:firstErr ?: @"nothing could be loaded"];
        else [self finish:PK_OK data:[ok componentsJoinedByString:@"\n"] msg:firstErr];
    });
}

@end

void pkPick(uintptr_t h, int limit, const char *cdir) {
    NSString *dir = [NSString stringWithUTF8String:cdir]; // copy now: Go frees cdir on return
    if ([NSThread isMainThread]) {
        // We would wait for the main thread on the main thread.
        done(h, PK_ERROR, nil, @"Pick must be called from a goroutine, not the main thread");
        return;
    }
    dispatch_async(dispatch_get_main_queue(), ^{
        if (@available(macOS 13.0, *)) {
            PHPickerConfiguration *c = [[PHPickerConfiguration alloc] init];
            c.selectionLimit = limit;
            c.filter = [PHPickerFilter imagesFilter];
            // The photo as the library holds it (HEIC stays HEIC), not a
            // copy converted to JPEG.
            c.preferredAssetRepresentationMode = PHPickerConfigurationAssetRepresentationModeCurrent;
            PHPickerViewController *vc = [[PHPickerViewController alloc] initWithConfiguration:c];

            PKPick *d = [PKPick new];
            d.handle = h;
            d.dir = dir;
            vc.delegate = d;
            if (!pkLive) pkLive = [NSMutableSet set];
            [pkLive addObject:d];

            NSWindow *w = [NSWindow windowWithContentViewController:vc];
            w.styleMask &= ~NSWindowStyleMaskClosable; // closing must go through the picker's Cancel
            [w setContentSize:NSMakeSize(880, 620)];
            d.window = w;
            NSWindow *parent = NSApp.keyWindow ?: NSApp.mainWindow;
            if (parent) {
                d.parent = parent;
                [parent beginSheet:w completionHandler:nil];
            } else {
                [w center];
                [w makeKeyAndOrderFront:nil];
                [NSApp activateIgnoringOtherApps:YES];
            }
        } else {
            done(h, PK_UNSUPPORTED, nil, nil);
        }
    });
}

void pkSave(uintptr_t h, const char *cpath, const char *cname) {
    NSString *path = [NSString stringWithUTF8String:cpath];
    NSString *name = [NSString stringWithUTF8String:cname];
    [PHPhotoLibrary requestAuthorizationForAccessLevel:PHAccessLevelAddOnly
                                               handler:^(PHAuthorizationStatus s) {
        if (s != PHAuthorizationStatusAuthorized && s != PHAuthorizationStatusLimited) {
            done(h, PK_DENIED, nil, nil);
            return;
        }
        __block NSString *ident = nil;
        [[PHPhotoLibrary sharedPhotoLibrary] performChanges:^{
            PHAssetCreationRequest *r = [PHAssetCreationRequest creationRequestForAsset];
            PHAssetResourceCreationOptions *o = [PHAssetResourceCreationOptions new];
            if (name.length) o.originalFilename = name;
            [r addResourceWithType:PHAssetResourceTypePhoto fileURL:[NSURL fileURLWithPath:path] options:o];
            ident = r.placeholderForCreatedAsset.localIdentifier;
        } completionHandler:^(BOOL ok, NSError *e) {
            if (ok) done(h, PK_OK, ident, nil);
            else done(h, PK_ERROR, nil, e.localizedDescription ?: @"saving failed");
        }];
    }];
}
