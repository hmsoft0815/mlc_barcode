<script lang="ts">
  import { onMount, onDestroy, tick } from 'svelte';
  import { DecodeImage, CopyToClipboard } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import type { DecodeImageResult } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import { BARCODE_TYPES } from '../types';
  import { formatError } from '../i18n/errors';
  import { KIND_LABELS, fieldLabel, formatFieldValue, isCheckKey, sortedFieldKeys } from '../i18n/fields';

  // Tabs stay mounted; paste (Ctrl+V) is only taken while this one is shown.
  export let active = false;

  let imageUrl = '';
  let fileName = '';
  let result: DecodeImageResult | null = null;
  let busy = false;
  let dragOver = false;
  let copied = -1;

  // Camera: live preview, frames go to the same decoder as files. The first
  // frame with a code stops the camera and stays as the checked image.
  const cameraSupported = typeof navigator !== 'undefined' && !!navigator.mediaDevices?.getUserMedia;
  let video: HTMLVideoElement;
  let stream: MediaStream | null = null;
  let cameraOn = false;
  let cameraError = '';
  let scanTimer: ReturnType<typeof setTimeout> | undefined;
  let cameras: MediaDeviceInfo[] = [];
  let facingUser = false; // front camera: preview is mirrored like a mirror
  let scanTick = 0;
  const SCAN_INTERVAL_MS = 300;
  const FRAME_MAX_SIDE = 1600;
  const GUIDE_INSET = 0.18; // matches .camera-frame

  // Copies a rectangle of an image or video frame into a JPEG data URL,
  // at most FRAME_MAX_SIDE on the long side.
  function grab(source: CanvasImageSource, r: { x: number; y: number; w: number; h: number }) {
    const scale = Math.min(1, FRAME_MAX_SIDE / Math.max(r.w, r.h));
    const canvas = document.createElement('canvas');
    canvas.width = Math.max(1, Math.round(r.w * scale));
    canvas.height = Math.max(1, Math.round(r.h * scale));
    canvas.getContext('2d')!.drawImage(source, r.x, r.y, r.w, r.h, 0, 0, canvas.width, canvas.height);
    return canvas.toDataURL('image/jpeg', 0.92);
  }

  // Region selection on a still image: drag a rectangle, only that part is
  // decoded; the marks are shifted back onto the whole image.
  let imgEl: HTMLImageElement;
  let cropMode = false;
  let sel: { x0: number; y0: number; x1: number; y1: number } | null = null;

  function relPos(e: PointerEvent) {
    const b = imgEl.getBoundingClientRect();
    return { x: Math.min(1, Math.max(0, (e.clientX - b.left) / b.width)), y: Math.min(1, Math.max(0, (e.clientY - b.top) / b.height)) };
  }

  function selStart(e: PointerEvent) {
    if (!cropMode) return;
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
    const p = relPos(e);
    sel = { x0: p.x, y0: p.y, x1: p.x, y1: p.y };
  }

  function selMove(e: PointerEvent) {
    if (!cropMode || !sel) return;
    const p = relPos(e);
    sel = { ...sel, x1: p.x, y1: p.y };
  }

  async function selEnd() {
    if (!cropMode || !sel) return;
    const nw = imgEl.naturalWidth, nh = imgEl.naturalHeight;
    const x = Math.min(sel.x0, sel.x1) * nw, y = Math.min(sel.y0, sel.y1) * nh;
    const w = Math.abs(sel.x1 - sel.x0) * nw, h = Math.abs(sel.y1 - sel.y0) * nh;
    if (w < 16 || h < 16) {
      sel = null; // a tap, not a drag
      return;
    }
    const crop = { x, y, w, h };
    const scale = Math.min(1, FRAME_MAX_SIDE / Math.max(w, h));
    busy = true;
    result = null;
    try {
      const r = await DecodeImage(grab(imgEl, crop));
      for (const c of r.codes ?? []) {
        c.points = (c.points ?? []).map((pt) => ({ x: crop.x + pt.x / scale, y: crop.y + pt.y / scale }));
      }
      result = { ...r, width: nw, height: nh };
    } catch (e: any) {
      result = { success: false, codes: [], width: 0, height: 0, error: e?.message ?? String(e) };
    } finally {
      busy = false;
      cropMode = false;
      sel = null;
    }
  }

  $: if (!active && cameraOn) stopCamera();

  const typeName = (id: string) => BARCODE_TYPES.find((t) => t.id === id)?.name ?? id.toUpperCase();

  async function check(dataUrl: string, name: string) {
    cropMode = false;
    sel = null;
    imageUrl = dataUrl;
    fileName = name;
    busy = true;
    result = null;
    try {
      result = await DecodeImage(dataUrl);
    } catch (e: any) {
      result = { success: false, codes: [], width: 0, height: 0, error: e?.message ?? String(e) };
    } finally {
      busy = false;
    }
  }

  async function startCamera(deviceId?: string) {
    cameraError = '';
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        video: {
          ...(deviceId ? { deviceId: { exact: deviceId } } : { facingMode: { ideal: 'environment' } }),
          width: { ideal: 1920 },
          height: { ideal: 1080 },
        },
        audio: false,
      });
      // Labels and the full list are only available once access was granted.
      cameras = (await navigator.mediaDevices.enumerateDevices()).filter((d) => d.kind === 'videoinput');
      facingUser = stream.getVideoTracks()[0]?.getSettings().facingMode === 'user';
    } catch (e: any) {
      cameraError =
        e?.name === 'NotAllowedError' ? 'Der Zugriff auf die Kamera wurde nicht erlaubt.'
        : e?.name === 'NotFoundError' || e?.name === 'OverconstrainedError' ? 'Keine Kamera gefunden.'
        : e?.name === 'NotReadableError' ? 'Die Kamera wird gerade von einem anderen Programm verwendet.'
        : `Kamera nicht verfügbar (${e?.name ?? e}).`;
      return;
    }
    cameraOn = true;
    cropMode = false;
    imageUrl = '';
    result = null;
    await tick(); // the video element exists only after this render
    if (!video || !stream) return;
    video.srcObject = stream;
    video.play().catch(() => {});
    scanTimer = setTimeout(scanFrame, SCAN_INTERVAL_MS);
  }

  async function switchCamera() {
    const current = stream?.getVideoTracks()[0]?.getSettings().deviceId;
    const i = cameras.findIndex((c) => c.deviceId === current);
    const next = cameras[(i + 1) % cameras.length];
    stopCamera();
    await startCamera(next?.deviceId);
  }

  function stopCamera() {
    clearTimeout(scanTimer);
    stream?.getTracks().forEach((t) => t.stop());
    stream = null;
    cameraOn = false;
    if (video) video.srcObject = null;
  }

  // One frame at a time: the next is taken only after the decoder answered,
  // so a slow device never builds up a queue.
  async function scanFrame() {
    if (!cameraOn || !video) return;
    const w = video.videoWidth, h = video.videoHeight;
    if (w > 0 && h > 0) {
      // Every other frame only the area inside the guide frame, at full
      // camera resolution: small codes (DataMatrix on a pack) get more
      // pixels there than in the downscaled whole frame.
      const guide = scanTick++ % 2 === 0;
      const src = guide ? { x: w * GUIDE_INSET, y: h * GUIDE_INSET, w: w * (1 - 2 * GUIDE_INSET), h: h * (1 - 2 * GUIDE_INSET) } : { x: 0, y: 0, w, h };
      const dataUrl = grab(video, src);
      try {
        const r = await DecodeImage(dataUrl);
        if (cameraOn && r.success && (r.codes?.length ?? 0) > 0) {
          stopCamera();
          imageUrl = dataUrl;
          fileName = guide ? 'Kamera (Zielrahmen)' : 'Kamera';
          result = r;
          return;
        }
      } catch {
        // A single bad frame is no reason to stop scanning.
      }
    }
    if (cameraOn) scanTimer = setTimeout(scanFrame, SCAN_INTERVAL_MS);
  }

  onDestroy(stopCamera);

  function readFile(file: File | null | undefined) {
    if (!file) return;
    if (cameraOn) stopCamera();
    const reader = new FileReader();
    reader.onload = () => check(reader.result as string, file.name);
    reader.readAsDataURL(file);
  }

  function onDrop(e: DragEvent) {
    e.preventDefault();
    dragOver = false;
    readFile(e.dataTransfer?.files?.[0]);
  }

  function onPaste(e: ClipboardEvent) {
    if (!active) return;
    const item = Array.from(e.clipboardData?.items ?? []).find((i) => i.type.startsWith('image/'));
    if (item) {
      e.preventDefault();
      readFile(item.getAsFile());
    }
  }

  // Frame around the reported points. QR points are the centres of the
  // finder patterns (plus the alignment pattern), 1D codes give two points on
  // a line — a padded bounding box fits both.
  function markBox(points: { x: number; y: number }[], w: number, h: number) {
    if (points.length === 0) return null;
    const xs = points.map((p) => p.x);
    const ys = points.map((p) => p.y);
    const minX = Math.min(...xs), maxX = Math.max(...xs), minY = Math.min(...ys), maxY = Math.max(...ys);
    const pad = Math.max((maxX - minX) * 0.2, (maxY - minY) * 0.2, Math.max(w, h) * 0.02);
    const x = Math.max(0, minX - pad), y = Math.max(0, minY - pad);
    return { x, y, w: Math.min(w, maxX + pad) - x, h: Math.min(h, maxY + pad) - y, pad };
  }

  async function copyText(text: string, i: number) {
    if (!(await CopyToClipboard(text))) await navigator.clipboard.writeText(text);
    copied = i;
    setTimeout(() => (copied = -1), 1500);
  }

  onMount(() => {
    window.addEventListener('paste', onPaste);
    return () => window.removeEventListener('paste', onPaste);
  });
</script>

<div class="container-fluid py-3">
  <div class="row g-3">
    <!-- Image input and preview -->
    <div class="col-lg-6">
      <div class="card shadow-sm border h-100">
        <div class="card-header bg-body border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-qr-code-scan me-1 text-primary"></i> Barcode prüfen
          </h6>
        </div>
        <div class="card-body">
          <div
            class="drop-zone rounded border p-3 text-center {dragOver ? 'drag-over' : ''}"
            role="region"
            aria-label="Bild ablegen"
            on:dragover|preventDefault={() => (dragOver = true)}
            on:dragleave={() => (dragOver = false)}
            on:drop={onDrop}
          >
            {#if cameraOn}
              <div class="camera">
                <!-- svelte-ignore a11y-media-has-caption -->
                <video bind:this={video} class:mirrored={facingUser} autoplay playsinline muted></video>
                <div class="camera-frame" aria-hidden="true"></div>
              </div>
              <div class="small text-body-secondary mt-2">
                <span class="spinner-grow spinner-grow-sm text-primary me-1"></span> Code ins Bild halten …
              </div>
            {:else if imageUrl}
              <div class="preview" class:cropping={cropMode}>
                <img bind:this={imgEl} src={imageUrl} alt={fileName} />
                {#if result?.success && result.width > 0}
                  <svg viewBox="0 0 {result.width} {result.height}" preserveAspectRatio="xMidYMid meet" aria-hidden="true">
                    {#each result.codes ?? [] as code, i}
                      {@const box = markBox(code.points ?? [], result.width, result.height)}
                      {#if box}
                        <rect class="mark" x={box.x} y={box.y} width={box.w} height={box.h} rx={box.pad / 3} />
                        <text class="mark-label" x={box.x + box.pad / 2} y={box.y - box.pad / 3} font-size={Math.max(result.width, result.height) / 30}>{i + 1}</text>
                      {/if}
                    {/each}
                  </svg>
                {/if}
                {#if cropMode}
                  <div
                    class="crop-layer"
                    role="presentation"
                    on:pointerdown={selStart}
                    on:pointermove={selMove}
                    on:pointerup={selEnd}
                    on:pointercancel={() => (sel = null)}
                  >
                    {#if sel}
                      <div
                        class="crop-rect"
                        style="left:{Math.min(sel.x0, sel.x1) * 100}%;top:{Math.min(sel.y0, sel.y1) * 100}%;width:{Math.abs(sel.x1 - sel.x0) * 100}%;height:{Math.abs(sel.y1 - sel.y0) * 100}%"
                      ></div>
                    {/if}
                  </div>
                {/if}
              </div>
              <div class="small text-body-secondary mt-2 text-truncate">
                {cropMode ? 'Rahmen um den Code ziehen …' : fileName}
              </div>
            {:else}
              <i class="bi bi-image fs-1 d-block mb-2 opacity-50"></i>
              <p class="mb-1 text-body">Bild hierher ziehen, auswählen oder mit <kbd>Strg</kbd>+<kbd>V</kbd> einfügen</p>
              <p class="small text-body-secondary mb-0">PNG, JPEG, GIF oder WebP · Fotos, Scans, Screenshots</p>
            {/if}
          </div>
          <div class="d-flex flex-wrap gap-2 mt-3">
            <label class="btn btn-outline-primary btn-sm mb-0">
              <i class="bi bi-folder2-open me-1"></i> Bild auswählen …
              <input
                type="file"
                accept="image/png,image/jpeg,image/gif,image/webp"
                class="d-none"
                on:change={(e) => readFile(e.currentTarget.files?.[0])}
              />
            </label>
            {#if cameraSupported}
              {#if cameraOn}
                {#if cameras.length > 1}
                  <button type="button" class="btn btn-outline-primary btn-sm" on:click={switchCamera}>
                    <i class="bi bi-arrow-repeat me-1"></i> Kamera wechseln
                  </button>
                {/if}
                <button type="button" class="btn btn-outline-secondary btn-sm" on:click={stopCamera}>
                  <i class="bi bi-camera-video-off me-1"></i> Kamera beenden
                </button>
              {:else}
                <button type="button" class="btn btn-outline-primary btn-sm" on:click={() => startCamera()}>
                  <i class="bi bi-camera me-1"></i> Mit Kamera scannen
                </button>
              {/if}
            {/if}
            {#if imageUrl && !cameraOn}
              <button
                type="button"
                class="btn btn-sm {cropMode ? 'btn-primary' : 'btn-outline-primary'}"
                on:click={() => ((cropMode = !cropMode), (sel = null))}
              >
                <i class="bi bi-crop me-1"></i> {cropMode ? 'Auswahl abbrechen' : 'Ausschnitt wählen'}
              </button>
            {/if}
          </div>
          {#if cameraError}
            <div class="alert alert-warning small mt-2 mb-0">
              <i class="bi bi-camera-video-off me-1"></i> {cameraError}
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Results -->
    <div class="col-lg-6">
      <div class="card shadow-sm border h-100">
        <div class="card-header bg-body border-bottom py-2 d-flex justify-content-between align-items-center">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-list-check me-1 text-primary"></i> Ergebnis
          </h6>
          {#if busy}
            <span class="badge bg-primary-subtle text-primary small">
              <span class="spinner-border spinner-border-sm me-1"></span> Lese …
            </span>
          {:else if result?.success}
            <span class="badge bg-success-subtle text-success small">{result.codes?.length ?? 0} Code{result.codes?.length === 1 ? '' : 's'} gefunden</span>
          {/if}
        </div>
        <div class="card-body">
          {#if result && !result.success}
            <div class="alert alert-warning mb-0">
              <i class="bi bi-exclamation-triangle me-1"></i> {formatError(result)}
            </div>
          {:else if result?.success}
            {#each result.codes ?? [] as code, i}
              <div class="border rounded p-3 mb-3">
                <div class="d-flex justify-content-between align-items-center mb-2">
                  <div>
                    <span class="badge bg-primary me-1">{i + 1}</span>
                    <span class="fw-semibold text-body">{typeName(code.type)}</span>
                    {#if code.content}
                      <span class="badge bg-info-subtle text-info-emphasis ms-1">{KIND_LABELS[code.content.kind] ?? code.content.kind}</span>
                    {/if}
                  </div>
                  <button type="button" class="btn btn-sm btn-outline-secondary" on:click={() => copyText(code.text, i)}>
                    <i class="bi {copied === i ? 'bi-check2' : 'bi-clipboard'} me-1"></i>{copied === i ? 'Kopiert' : 'Inhalt kopieren'}
                  </button>
                </div>
                {#if code.content?.fields}
                  <table class="table table-sm mb-2">
                    <tbody>
                      {#each sortedFieldKeys(code.content.fields) as key}
                        <tr>
                          <th class="text-body-secondary fw-normal small" style="width: 40%">{fieldLabel(key)}</th>
                          <td
                            class="small {(isCheckKey(key) && code.content.fields[key] === 'false') || (key === 'expiry' && formatFieldValue(key, code.content.fields[key]).endsWith('abgelaufen'))
                              ? 'text-danger fw-semibold'
                              : ''}"
                          >
                            <span class="value">{formatFieldValue(key, code.content.fields[key])}</span>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                {/if}
                <details open={!code.content}>
                  <summary class="small text-body-secondary">Rohinhalt</summary>
                  <pre class="raw bg-body-secondary border rounded p-2 mt-1 mb-0">{code.text}</pre>
                </details>
              </div>
            {/each}
          {:else if !busy}
            <p class="text-body-secondary small mb-0">
              Liest QR, DataMatrix, Aztec, PDF417, EAN-13/8, UPC-A, Code 128, Code 39 und ITF – auch mehrere Codes pro Bild.
              Bekannte Inhalte wie GiroCode, Visitenkarte, WLAN oder Termin werden in ihre Felder zerlegt;
              beim GiroCode wird die IBAN-Prüfsumme kontrolliert, bei Arzneimittel-Codes (securPharm) PZN, Charge, Verfall und Seriennummer angezeigt.
            </p>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .drop-zone {
    min-height: 260px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border-style: dashed !important;
    transition: background-color 0.15s ease;
  }
  .drop-zone.drag-over {
    background-color: var(--bs-primary-bg-subtle);
    border-color: var(--bs-primary) !important;
  }
  /* The box hugs the video (no cropping), so the guide frame covers exactly
     the part that GUIDE_INSET cuts out of the camera frame. */
  .camera {
    position: relative;
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    border-radius: var(--bs-border-radius);
  }
  .camera video {
    display: block;
    max-width: 100%;
    max-height: 420px;
    background: #000;
  }
  .camera video.mirrored {
    transform: scaleX(-1);
  }
  .camera-frame {
    position: absolute;
    inset: 18%;
    border: 3px solid rgba(var(--bs-primary-rgb), 0.9);
    border-radius: 12px;
    box-shadow: 0 0 0 999px rgba(0, 0, 0, 0.25);
    pointer-events: none;
  }
  .preview {
    position: relative;
    display: inline-block;
    max-width: 100%;
  }
  .preview img {
    display: block;
    max-width: 100%;
    max-height: 420px;
    background: #fff;
  }
  .preview svg {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
  }
  .crop-layer {
    position: absolute;
    inset: 0;
    cursor: crosshair;
    touch-action: none;
    background: rgba(0, 0, 0, 0.15);
  }
  .crop-rect {
    position: absolute;
    border: 2px dashed var(--bs-primary);
    background: rgba(var(--bs-primary-rgb), 0.15);
  }
  .preview.cropping img {
    user-select: none;
    -webkit-user-drag: none;
  }
  .mark {
    fill: rgba(var(--bs-primary-rgb), 0.15);
    stroke: var(--bs-primary);
    stroke-width: 0.6%;
  }
  .mark-label {
    fill: var(--bs-primary);
    font-weight: 700;
  }
  .raw {
    white-space: pre-wrap;
    word-break: break-all;
    font-size: 0.75rem;
    max-height: 12rem;
    overflow-y: auto;
  }
  .value {
    white-space: pre-wrap;
    word-break: break-word;
  }
</style>
