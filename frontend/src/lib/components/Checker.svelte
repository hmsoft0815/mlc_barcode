<script lang="ts">
  import { onMount, onDestroy, tick } from 'svelte';
  import { DecodeImage, CopyToClipboard } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import type { DecodeImageResult } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import { BARCODE_TYPES } from '../types';
  import { formatError } from '../i18n/errors';
  import { kindLabel } from '../i18n/fields';
  import { lang } from '../i18n/lang';
  import { t } from '../i18n/text/checker';
  import CodeDetails from './CodeDetails.svelte';
  import { CameraScanner, cameraErrorText, cameraSupported, grab, insetGuide, FRAME_MAX_SIDE, type ScanHit } from '../scan/camera';

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
  const GUIDE_INSET = 0.18; // matches .camera-frame
  let video: HTMLVideoElement;
  let cameraOn = false;
  let cameraError = '';
  let scanner = new CameraScanner(onCameraHit, insetGuide(GUIDE_INSET));

  $: if (!active && cameraOn) stopCamera();

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

  async function startCamera() {
    cameraError = '';
    cropMode = false;
    imageUrl = '';
    result = null;
    cameraOn = true;
    await tick(); // the video element exists only after this render
    try {
      await scanner.start(video);
      scanner = scanner; // cameras and facing are known now
    } catch (e: any) {
      cameraError = cameraErrorText(e, $lang);
      cameraOn = false;
    }
  }

  async function switchCamera() {
    try {
      await scanner.switchCamera();
      scanner = scanner;
    } catch (e: any) {
      cameraError = cameraErrorText(e, $lang);
      stopCamera();
    }
  }

  function stopCamera() {
    scanner.stop();
    cameraOn = false;
  }

  function onCameraHit(hit: ScanHit) {
    stopCamera();
    imageUrl = hit.dataUrl;
    fileName = hit.guide ? $t('cameraGuide') : $t('camera');
    result = hit.result;
  }

  onDestroy(stopCamera);

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
            <i class="bi bi-qr-code-scan me-1 text-primary"></i> {$t('title')}
          </h6>
        </div>
        <div class="card-body">
          <div
            class="drop-zone rounded border p-3 text-center {dragOver ? 'drag-over' : ''}"
            role="region"
            aria-label={$t('dropLabel')}
            on:dragover|preventDefault={() => (dragOver = true)}
            on:dragleave={() => (dragOver = false)}
            on:drop={onDrop}
          >
            {#if cameraOn}
              <div class="camera">
                <!-- svelte-ignore a11y-media-has-caption -->
                <video bind:this={video} class:mirrored={scanner.facingUser} autoplay playsinline muted disablepictureinpicture></video>
                <div class="camera-frame" aria-hidden="true"></div>
              </div>
              <div class="small text-body-secondary mt-2">
                <span class="spinner-grow spinner-grow-sm text-primary me-1"></span> {$t('holdInImage')}
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
                {cropMode ? $t('dragFrame') : fileName}
              </div>
            {:else}
              <i class="bi bi-image fs-1 d-block mb-2 opacity-50"></i>
              <p class="mb-1 text-body">{$t('dropHint', { key: $t('pasteKey') })}</p>
              <p class="small text-body-secondary mb-0">{$t('formats')}</p>
            {/if}
          </div>
          <div class="d-flex flex-wrap gap-2 mt-3">
            <label class="btn btn-outline-primary btn-sm mb-0">
              <i class="bi bi-folder2-open me-1"></i> {$t('pickImage')}
              <input
                type="file"
                accept="image/png,image/jpeg,image/gif,image/webp"
                class="d-none"
                on:change={(e) => readFile(e.currentTarget.files?.[0])}
              />
            </label>
            {#if cameraSupported}
              {#if cameraOn}
                {#if scanner.cameras.length > 1}
                  <button type="button" class="btn btn-outline-primary btn-sm" on:click={switchCamera}>
                    <i class="bi bi-arrow-repeat me-1"></i> {$t('switchCamera')}
                  </button>
                {/if}
                <button type="button" class="btn btn-outline-secondary btn-sm" on:click={stopCamera}>
                  <i class="bi bi-camera-video-off me-1"></i> {$t('stopCamera')}
                </button>
              {:else}
                <button type="button" class="btn btn-outline-primary btn-sm" on:click={startCamera}>
                  <i class="bi bi-camera me-1"></i> {$t('scanWithCamera')}
                </button>
              {/if}
            {/if}
            {#if imageUrl && !cameraOn}
              <button
                type="button"
                class="btn btn-sm {cropMode ? 'btn-primary' : 'btn-outline-primary'}"
                on:click={() => ((cropMode = !cropMode), (sel = null))}
              >
                <i class="bi bi-crop me-1"></i> {cropMode ? $t('cancelSelection') : $t('selectRegion')}
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
            <i class="bi bi-list-check me-1 text-primary"></i> {$t('result')}
          </h6>
          {#if busy}
            <span class="badge bg-primary-subtle text-primary small">
              <span class="spinner-border spinner-border-sm me-1"></span> {$t('reading')}
            </span>
          {:else if result?.success}
            <span class="badge bg-success-subtle text-success small">{result.codes?.length === 1 ? $t('foundOne') : $t('foundMany', { n: result.codes?.length ?? 0 })}</span>
          {/if}
        </div>
        <div class="card-body">
          {#if result && !result.success}
            <div class="alert alert-warning mb-0">
              <i class="bi bi-exclamation-triangle me-1"></i> {formatError(result, $lang)}
            </div>
          {:else if result?.success}
            {#each result.codes ?? [] as code, i}
              <div class="border rounded p-3 mb-3">
                <div class="d-flex justify-content-between align-items-center mb-2">
                  <div>
                    <span class="badge bg-primary me-1">{i + 1}</span>
                    <span class="fw-semibold text-body">{typeName(code.type)}</span>
                    {#if code.content}
                      <span class="badge bg-info-subtle text-info-emphasis ms-1">{kindLabel(code.content.kind, $lang)}</span>
                    {/if}
                  </div>
                  <button type="button" class="btn btn-sm btn-outline-secondary" on:click={() => copyText(code.text, i)}>
                    <i class="bi {copied === i ? 'bi-check2' : 'bi-clipboard'} me-1"></i>{copied === i ? $t('copied') : $t('copyContent')}
                  </button>
                </div>
                <CodeDetails {code} lang={$lang} />
              </div>
            {/each}
          {:else if !busy}
            <p class="text-body-secondary small mb-0">
              {$t('intro')}
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
</style>
