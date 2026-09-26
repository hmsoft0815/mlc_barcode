<script lang="ts">
  import { onMount } from 'svelte';
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

  const typeName = (id: string) => BARCODE_TYPES.find((t) => t.id === id)?.name ?? id.toUpperCase();

  async function check(dataUrl: string, name: string) {
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
            {#if imageUrl}
              <div class="preview">
                <img src={imageUrl} alt={fileName} />
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
              </div>
              <div class="small text-body-secondary mt-2 text-truncate">{fileName}</div>
            {:else}
              <i class="bi bi-image fs-1 d-block mb-2 opacity-50"></i>
              <p class="mb-1 text-body">Bild hierher ziehen, auswählen oder mit <kbd>Strg</kbd>+<kbd>V</kbd> einfügen</p>
              <p class="small text-body-secondary mb-0">PNG, JPEG, GIF oder WebP · Fotos, Scans, Screenshots</p>
            {/if}
          </div>
          <label class="btn btn-outline-primary btn-sm mt-3">
            <i class="bi bi-folder2-open me-1"></i> Bild auswählen …
            <input
              type="file"
              accept="image/png,image/jpeg,image/gif,image/webp"
              class="d-none"
              on:change={(e) => readFile(e.currentTarget.files?.[0])}
            />
          </label>
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
