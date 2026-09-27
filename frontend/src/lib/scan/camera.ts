// Live camera scanning, shared by the desktop "Prüfen" tab and the mobile
// scan screen. Frames go to the same decoder as files; the
// pipeline is described in docs/decoder.md.
// Frames go through DecodeCameraFrame: like DecodeImage, but the slow search
// for codes on bottles is cut short — the next frame follows anyway.
import { DecodeCameraFrame } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
import type { DecodeImageResult } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
import type { Lang } from '../i18n/errors';
import { UI_TEXT } from '../i18n/ui';

export const SCAN_INTERVAL_MS = 300;
export const FRAME_MAX_SIDE = 1600;

export const cameraSupported = typeof navigator !== 'undefined' && !!navigator.mediaDevices?.getUserMedia;

export type Rect = { x: number; y: number; w: number; h: number };

export interface ScanHit {
  dataUrl: string; // the frame (or guide area) the codes were found in
  result: DecodeImageResult;
  guide: boolean; // found in the guide area only
}

// Copies a rectangle of an image or video frame into a JPEG data URL, at
// most FRAME_MAX_SIDE on the long side; mirror flips it left to right.
export function grab(source: CanvasImageSource, r: Rect, mirror = false): string {
  const scale = Math.min(1, FRAME_MAX_SIDE / Math.max(r.w, r.h));
  const canvas = document.createElement('canvas');
  canvas.width = Math.max(1, Math.round(r.w * scale));
  canvas.height = Math.max(1, Math.round(r.h * scale));
  const ctx = canvas.getContext('2d')!;
  if (mirror) {
    ctx.translate(canvas.width, 0);
    ctx.scale(-1, 1);
  }
  ctx.drawImage(source, r.x, r.y, r.w, r.h, 0, 0, canvas.width, canvas.height);
  return canvas.toDataURL('image/jpeg', 0.92);
}

// cameraErrorText words an error thrown by CameraScanner.start.
export function cameraErrorText(e: any, lang: Lang): string {
  const t = UI_TEXT[lang];
  switch (e?.name) {
    case 'NotAllowedError':
      return t.camDenied;
    case 'NotFoundError':
    case 'OverconstrainedError':
      return t.camNone;
    case 'NotReadableError':
      return t.camBusy;
  }
  return t.camUnavailable.replace('{name}', String(e?.name ?? e));
}

// CameraScanner runs the camera into a <video> and decodes one frame at a
// time: the next is taken only after the decoder answered, so a slow device
// never builds up a queue. Every other frame only the guide area is sent,
// at full camera resolution — small codes (a DataMatrix on a pack) get more
// pixels there than in the reduced whole frame. The first frame with a code
// is reported to onHit and scanning stops.
export class CameraScanner {
  stream: MediaStream | null = null;
  cameras: MediaDeviceInfo[] = [];
  facingUser = false; // front camera: the preview should be mirrored
  torchOn = false;
  private video: HTMLVideoElement | null = null;
  private timer: ReturnType<typeof setTimeout> | undefined;
  private tick = 0;
  private scanning = false;

  constructor(
    private onHit: (hit: ScanHit) => void,
    // The guide area in video pixels, or null for whole frames only.
    private guide: (video: HTMLVideoElement) => Rect | null = () => null,
    // Called for every frame sent to the decoder, e.g. to show that
    // something is happening.
    private onAttempt: () => void = () => {},
    // What the code sits on ("auto", "flat", "bottle", "tube"): narrows the
    // decoder's search for codes on curved surfaces.
    private surface: () => string = () => 'auto',
    // What kind of code is scanned ("auto", "square", "wide"): only those
    // readers run — 2D alone is about ten times faster per frame.
    private family: () => string = () => 'auto'
  ) {}

  get running() {
    return this.stream !== null;
  }

  // Starts the camera (the rear one unless deviceId says otherwise) and
  // scanning. Throws getUserMedia's error; cameraErrorText words it.
  async start(video: HTMLVideoElement, choice: { deviceId?: string; facing?: 'user' | 'environment' } = {}) {
    // On a switch the old stream ends only after the new one runs, and the
    // video element never goes without a source: iOS falls back to its
    // fullscreen player when play() meets an element in between.
    const old = this.stream;
    this.pause();
    const which = choice.deviceId
      ? { deviceId: { exact: choice.deviceId } }
      : { facingMode: choice.facing ? { exact: choice.facing } : { ideal: 'environment' } };
    old?.getTracks().forEach((t) => t.stop()); // iOS allows one camera at a time
    try {
      this.stream = await navigator.mediaDevices.getUserMedia({
        video: { ...which, width: { ideal: 1920 }, height: { ideal: 1080 } },
        audio: false,
      });
    } catch (e) {
      this.stream = null;
      throw e;
    }
    // Labels and the full list are only available once access was granted.
    this.cameras = (await navigator.mediaDevices.enumerateDevices()).filter((d) => d.kind === 'videoinput');
    this.facingUser = this.track()?.getSettings().facingMode === 'user';
    this.video = video;
    // Inline only, never the system player (iOS); set here as well as in
    // the markup, since the element may come from either UI.
    video.muted = true;
    video.playsInline = true;
    video.setAttribute('playsinline', '');
    video.setAttribute('webkit-playsinline', '');
    video.disablePictureInPicture = true;
    video.srcObject = this.stream;
    video.play().catch(() => {});
    this.resume();
  }

  // Scans again after a hit, with the camera still running.
  resume() {
    if (!this.stream || this.scanning) return;
    this.scanning = true;
    this.timer = setTimeout(() => this.scanFrame(), SCAN_INTERVAL_MS);
  }

  pause() {
    this.scanning = false;
    clearTimeout(this.timer);
  }

  stop() {
    this.pause();
    this.stream?.getTracks().forEach((t) => t.stop());
    this.stream = null;
    this.torchOn = false;
    if (this.video) this.video.srcObject = null;
    this.video = null;
  }

  // Switches between front and rear camera. Going through the device list
  // instead lands on the second rear camera (ultra-wide) of phones and
  // tablets; that is only the fallback for cameras without a facing mode,
  // e.g. two webcams on a desktop.
  async switchCamera() {
    const video = this.video;
    if (!video || this.cameras.length < 2) return;
    const settings = this.track()?.getSettings();
    if (settings?.facingMode) {
      try {
        await this.start(video, { facing: settings.facingMode === 'user' ? 'environment' : 'user' });
        return;
      } catch {
        // no camera on the other side: fall back to the device list
      }
    }
    const i = this.cameras.findIndex((c) => c.deviceId === settings?.deviceId);
    await this.start(video, { deviceId: this.cameras[(i + 1) % this.cameras.length]?.deviceId });
  }

  // The camera light through the track, where the browser supports it
  // (Android). Returns false when it is not available this way.
  async setTorch(on: boolean): Promise<boolean> {
    const track = this.track();
    const caps = (track?.getCapabilities?.() ?? {}) as { torch?: boolean };
    if (!track || !caps.torch) return false;
    try {
      await track.applyConstraints({ advanced: [{ torch: on } as MediaTrackConstraintSet] });
      this.torchOn = on;
      return true;
    } catch {
      return false;
    }
  }

  private track() {
    return this.stream?.getVideoTracks()[0];
  }

  private async scanFrame() {
    const video = this.video;
    if (!this.scanning || !video) return;
    const w = video.videoWidth, h = video.videoHeight;
    if (w > 0 && h > 0) {
      const tick = this.tick++;
      const g = tick % 2 === 0 ? this.guide(video) : null;
      // Front camera: every third frame also mirrored. Whether a browser
      // hands out the front camera's picture mirrored is not specified;
      // a mirrored DataMatrix or QR code is not read.
      const mirror = this.facingUser && tick % 3 === 2;
      const dataUrl = grab(video, g ?? { x: 0, y: 0, w, h }, mirror);
      this.onAttempt();
      try {
        const result = await DecodeCameraFrame(dataUrl, this.surface(), this.family(), this.facingUser);
        if (this.scanning && result.success && (result.codes?.length ?? 0) > 0) {
          this.pause();
          this.onHit({ dataUrl, result, guide: g !== null });
          return;
        }
      } catch {
        // A single bad frame is no reason to stop scanning.
      }
    }
    if (this.scanning) this.timer = setTimeout(() => this.scanFrame(), SCAN_INTERVAL_MS);
  }
}

// The guide area for a preview that shows the whole frame: an inset of
// the given fraction on every side.
export function insetGuide(inset: number) {
  return (video: HTMLVideoElement): Rect => {
    const w = video.videoWidth, h = video.videoHeight;
    return { x: w * inset, y: h * inset, w: w * (1 - 2 * inset), h: h * (1 - 2 * inset) };
  };
}

// The guide area for a preview with object-fit: cover, where the frame
// drawn over the video is `frame` (in CSS pixels of the video element).
export function coverGuide(frame: () => DOMRect | null) {
  return (video: HTMLVideoElement): Rect | null => {
    const f = frame();
    const box = video.getBoundingClientRect();
    if (!f || box.width === 0) return null;
    const s = Math.max(box.width / video.videoWidth, box.height / video.videoHeight);
    const offX = (video.videoWidth - box.width / s) / 2;
    const offY = (video.videoHeight - box.height / s) / 2;
    return { x: offX + (f.left - box.left) / s, y: offY + (f.top - box.top) / s, w: f.width / s, h: f.height / s };
  };
}
