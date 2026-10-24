package com.wails.app;

import android.app.Activity;
import android.content.ContentResolver;
import android.content.ContentValues;
import android.content.Intent;
import android.net.Uri;
import android.os.Build;
import android.os.Environment;
import android.provider.MediaStore;
import android.util.Base64;
import android.util.Log;
import android.webkit.JavascriptInterface;

import androidx.core.content.FileProvider;

import java.io.File;
import java.io.FileOutputStream;
import java.io.OutputStream;

/**
 * MLC Barcode's own calls from the page (window.mlcNative), for what the
 * Wails bridge does not offer on Android: saving a created code to the
 * gallery and sharing it as a picture rather than as text.
 *
 * The page passes the PNG as base64 (with or without "data:" prefix) and a
 * file name without extension. Both calls report false when they could not
 * do it, so the page can fall back.
 */
public class MlcNative {
    private static final String TAG = "MlcNative";
    private static final String ALBUM = "MLC Barcode";

    private final Activity activity;

    MlcNative(Activity activity) {
        this.activity = activity;
    }

    /**
     * Saves the PNG to Pictures/MLC Barcode via MediaStore. Android 10 and
     * newer need no storage permission for that; older versions report
     * false and the page shares the picture instead.
     */
    @JavascriptInterface
    public boolean saveImage(String pngBase64, String name) {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.Q) return false;
        byte[] png = decode(pngBase64);
        if (png == null) return false;
        ContentResolver resolver = activity.getContentResolver();
        ContentValues values = new ContentValues();
        values.put(MediaStore.Images.Media.DISPLAY_NAME, fileName(name));
        values.put(MediaStore.Images.Media.MIME_TYPE, "image/png");
        values.put(MediaStore.Images.Media.RELATIVE_PATH, Environment.DIRECTORY_PICTURES + "/" + ALBUM);
        values.put(MediaStore.Images.Media.IS_PENDING, 1);
        Uri uri = resolver.insert(MediaStore.Images.Media.EXTERNAL_CONTENT_URI, values);
        if (uri == null) return false;
        try (OutputStream out = resolver.openOutputStream(uri)) {
            if (out == null) throw new java.io.IOException("no output stream");
            out.write(png);
        } catch (Exception e) {
            Log.e(TAG, "saveImage failed", e);
            resolver.delete(uri, null, null);
            return false;
        }
        values.clear();
        values.put(MediaStore.Images.Media.IS_PENDING, 0);
        resolver.update(uri, values, null, null);
        return true;
    }

    /**
     * Opens the share sheet with the PNG as a file (cache/shared/, handed
     * out through the FileProvider) — messengers and mail get the picture.
     */
    @JavascriptInterface
    public boolean shareImage(String pngBase64, String name) {
        byte[] png = decode(pngBase64);
        if (png == null) return false;
        try {
            File dir = new File(activity.getCacheDir(), "shared");
            if (!dir.exists() && !dir.mkdirs()) return false;
            File file = new File(dir, fileName(name));
            try (FileOutputStream out = new FileOutputStream(file)) {
                out.write(png);
            }
            Uri uri = FileProvider.getUriForFile(activity, activity.getPackageName() + ".fileprovider", file);
            Intent send = new Intent(Intent.ACTION_SEND)
                    .setType("image/png")
                    .putExtra(Intent.EXTRA_STREAM, uri)
                    .addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION);
            activity.runOnUiThread(() -> activity.startActivity(Intent.createChooser(send, null)));
            return true;
        } catch (Exception e) {
            Log.e(TAG, "shareImage failed", e);
            return false;
        }
    }

    private static byte[] decode(String pngBase64) {
        if (pngBase64 == null) return null;
        int comma = pngBase64.indexOf(',');
        if (pngBase64.startsWith("data:") && comma > 0) pngBase64 = pngBase64.substring(comma + 1);
        try {
            return Base64.decode(pngBase64, Base64.DEFAULT);
        } catch (IllegalArgumentException e) {
            return null;
        }
    }

    // A safe file name: letters, digits, "-" and "_", at most 60 characters.
    private static String fileName(String name) {
        String safe = name == null ? "" : name.replaceAll("[^A-Za-z0-9_-]+", "-").replaceAll("^-+|-+$", "");
        if (safe.isEmpty()) safe = "qr-code";
        if (safe.length() > 60) safe = safe.substring(0, 60);
        return safe + ".png";
    }
}
