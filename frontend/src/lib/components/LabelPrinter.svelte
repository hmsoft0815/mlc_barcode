<script lang="ts">
  import { onMount } from 'svelte';
  import { GenerateBatch, PickTextFile } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import { BARCODE_TYPES, type BarcodeType, type PrintItem } from '../types';

  export let printItems: PrintItem[] = [];
  export let batchItems: PrintItem[] = [];

  // File import: first column = code, optional second column = label text.
  let importType: BarcodeType = 'qr';
  let skipHeader = false;
  let importMessage = '';
  let importOk = true;

  function takeFromBatch() {
    printItems = [...batchItems];
    importMessage = `${batchItems.length} Etiketten aus dem Batch-Generator übernommen.`;
    importOk = true;
  }

  function splitRow(line: string, sep: string): string[] {
    if (!sep) return [line];
    return line.split(sep).map((c) => c.trim().replace(/^"(.*)"$/, '$1').trim());
  }

  function detectSeparator(line: string): string {
    return [';', '\t', ','].find((s) => line.includes(s)) ?? '';
  }

  async function importFile() {
    try {
      const [path, lines] = await PickTextFile();
      if (!path || !lines?.length) return;

      const body = skipHeader ? lines.slice(1) : lines;
      const sep = path.toLowerCase().endsWith('.csv') ? detectSeparator(body[0] ?? '') : '';
      const rows = body.map((l) => splitRow(l, sep)).filter((r) => r[0]);

      const res = await GenerateBatch({
        type: importType,
        lines: rows.map((r) => r[0]),
        width: 0,
        height: 0,
        showText: false,
        fontSize: 0,
        foregroundColor: '#000000',
        backgroundColor: '#ffffff'
      });
      const items = res.items ?? [];
      printItems = items
        .filter((it) => it.success && it.svg)
        .map((it) => ({ data: rows[it.index - 1]?.[1] || it.data, svg: it.svg!, type: importType }));

      const name = path.split(/[\\/]/).pop();
      importMessage = `${printItems.length} Etiketten aus ${name} importiert`;
      if (res.errorCount) importMessage += `, ${res.errorCount} Zeilen ungültig für ${importType.toUpperCase()}`;
      importOk = !res.errorCount;
    } catch (e: any) {
      importMessage = `Import fehlgeschlagen: ${e?.message ?? e}`;
      importOk = false;
    }
  }

  let columns = 3;
  let rows = 8;
  let repeatCount = 1;
  let showDataText = true;

  const PRESETS = [
    { name: 'Avery 2x4 (8 Etiketten / 105x74mm)', cols: 2, rows: 4 },
    { name: 'Avery 3x7 (21 Etiketten / 70x42mm)', cols: 3, rows: 7 },
    { name: 'Avery 3x8 (24 Etiketten / 70x37mm)', cols: 3, rows: 8 },
    { name: 'Avery 4x10 (40 Etiketten / 52x29mm)', cols: 4, rows: 10 },
    { name: 'Einzel-Etikett / Sticker (1x1)', cols: 1, rows: 1 }
  ];

  function applyPreset(p: { cols: number; rows: number }) {
    columns = p.cols;
    rows = p.rows;
  }

  async function loadSampleLabels() {
    try {
      const res = await GenerateBatch({
        type: 'qr',
        lines: [
          'MLC-ART-001',
          'MLC-ART-002',
          'MLC-ART-003',
          'https://mlcgo.eu',
          'PROD-SN-998811',
          'BOX-ID-4402'
        ],
        width: 0,
        height: 0,
        showText: false,
        fontSize: 0,
        foregroundColor: '#000000',
        backgroundColor: '#ffffff'
      });
      if (res.items && res.items.length > 0) {
        printItems = res.items
          .filter((it) => it.success && it.svg)
          .map((it) => ({
            data: it.data,
            svg: it.svg!,
            type: 'qr'
          }));
      }
    } catch (e) {
      console.warn('Could not generate sample labels:', e);
    }
  }

  // Expanded items based on repeatCount
  $: expandedItems = printItems.flatMap((item) =>
    Array(repeatCount).fill(item)
  );

  function triggerPrint() {
    window.print();
  }

  onMount(() => {
    if (!printItems || printItems.length === 0) {
      loadSampleLabels();
    }
  });
</script>

<div class="container-fluid py-3">
  <!-- Controls (hidden when printing) -->
  <div class="card shadow-sm border mb-3 no-print">
    <div class="card-header bg-body border-bottom py-2 d-flex justify-content-between align-items-center">
      <h6 class="mb-0 fw-semibold text-body">
        <i class="bi bi-printer me-1 text-primary"></i> Etikettenbogen-Layout & Druckeinstellungen
      </h6>
      <div class="d-flex gap-2">
        <button
          class="btn btn-outline-primary btn-sm"
          disabled={batchItems.length === 0}
          title={batchItems.length ? `${batchItems.length} Barcodes aus dem Batch-Generator` : 'Im Batch-Generator ist noch nichts erzeugt'}
          on:click={takeFromBatch}
        >
          <i class="bi bi-collection me-1"></i> Aus Batch-Generator übernehmen
        </button>
        <button class="btn btn-outline-secondary btn-sm" on:click={loadSampleLabels}>
          <i class="bi bi-magic me-1"></i> Muster laden
        </button>
        <button
          class="btn btn-primary btn-sm"
          disabled={expandedItems.length === 0}
          on:click={triggerPrint}
        >
          <i class="bi bi-printer-fill me-1"></i> Drucken (Print Dialog)
        </button>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-3 align-items-center">
        <div class="col-md-4">
          <label for="labelPresetSelect" class="form-label small text-body-secondary mb-1">Standard-Vorlagen (DIN A4)</label>
          <select
            id="labelPresetSelect"
            class="form-select form-select-sm"
            on:change={(e) => {
              const val = e.currentTarget.value;
              const found = PRESETS.find((p) => p.name === val);
              if (found) applyPreset(found);
            }}
          >
            {#each PRESETS as p}
              <option value={p.name} selected={p.cols === columns && p.rows === rows}>
                {p.name}
              </option>
            {/each}
          </select>
        </div>

        <div class="col-md-2 col-6">
          <label for="labelColumnsInput" class="form-label small text-body-secondary mb-1">Spalten (Columns)</label>
          <input
            id="labelColumnsInput"
            type="number"
            class="form-control form-control-sm"
            min="1"
            max="6"
            bind:value={columns}
          />
        </div>

        <div class="col-md-2 col-6">
          <label for="labelRowsInput" class="form-label small text-body-secondary mb-1">Zeilen (Rows)</label>
          <input
            id="labelRowsInput"
            type="number"
            class="form-control form-control-sm"
            min="1"
            max="15"
            bind:value={rows}
          />
        </div>

        <div class="col-md-2 col-6">
          <label for="labelRepeatInput" class="form-label small text-body-secondary mb-1">Wiederholungen</label>
          <input
            id="labelRepeatInput"
            type="number"
            class="form-control form-control-sm"
            min="1"
            max="100"
            bind:value={repeatCount}
          />
        </div>

        <div class="col-md-2 col-6 pt-md-3">
          <div class="form-check form-switch small">
            <input
              class="form-check-input"
              type="checkbox"
              id="showDataLabel"
              bind:checked={showDataText}
            />
            <label class="form-check-label" for="showDataLabel">Text anzeigen</label>
          </div>
        </div>
      </div>

      <div class="row g-2 align-items-end border-top pt-3 mt-2">
        <div class="col-md-4">
          <label for="labelImportType" class="form-label small text-body-secondary mb-1">
            Datei importieren (TXT/CSV: 1. Spalte Code, 2. Spalte Etikett-Text)
          </label>
          <select id="labelImportType" class="form-select form-select-sm" bind:value={importType}>
            {#each BARCODE_TYPES as t}
              <option value={t.id}>{t.name}</option>
            {/each}
          </select>
        </div>
        <div class="col-md-3 col-6">
          <div class="form-check form-switch small mb-1">
            <input class="form-check-input" type="checkbox" id="labelSkipHeader" bind:checked={skipHeader} />
            <label class="form-check-label" for="labelSkipHeader">Erste Zeile ist Kopfzeile</label>
          </div>
        </div>
        <div class="col-md-2 col-6">
          <button class="btn btn-outline-primary btn-sm w-100" on:click={importFile}>
            <i class="bi bi-filetype-csv me-1"></i> Importieren …
          </button>
        </div>
        {#if importMessage}
          <div class="col-12 small {importOk ? 'text-success' : 'text-warning'}">{importMessage}</div>
        {/if}
      </div>
    </div>
  </div>

  <!-- Print Sheet Container -->
  {#if expandedItems.length > 0}
    <div class="print-page-wrapper">
      <div
        class="print-grid"
        style="--cols: {columns}; --rows: {rows};"
      >
        {#each expandedItems as item}
          <div class="label-cell">
            <div class="label-svg-wrapper">
              <!-- eslint-disable-next-line svelte/no-at-html-tags -->
              {@html item.svg}
            </div>
            {#if showDataText}
              <div class="label-text">{item.data}</div>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  {:else}
    <div class="card shadow-sm border py-5 text-center text-body-secondary no-print">
      <i class="bi bi-printer fs-1 opacity-50 mb-2"></i>
      <h5 class="text-body">Keine Etiketten in der Druck-Warteschlange</h5>
      <p class="small text-body-secondary mb-3">
        Erstelle Barcodes im <strong>Einzel-Generator</strong> oder <strong>Batch-Generator</strong> und klicke auf "Als Etikett drucken".
      </p>
      <div>
        <button class="btn btn-outline-primary btn-sm" on:click={loadSampleLabels}>
          <i class="bi bi-magic me-1"></i> Beispiel-Etiketten laden
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .print-page-wrapper {
    background: #ffffff;
    color: #000000;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    margin: 0 auto;
    padding: 10mm;
    max-width: 210mm;
    min-height: 297mm;
    box-sizing: border-box;
  }

  .print-grid {
    display: grid;
    grid-template-columns: repeat(var(--cols, 3), 1fr);
    grid-auto-rows: minmax(28mm, auto);
    gap: 3mm;
    width: 100%;
  }

  .label-cell {
    border: 1px dashed #dee2e6;
    padding: 2mm;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
    overflow: hidden;
    background: #ffffff;
    color: #000000;
  }

  .label-svg-wrapper {
    width: 100%;
    height: 20mm;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .label-svg-wrapper :global(svg) {
    width: 100%;
    height: 100%;
    max-width: 100%;
    max-height: 100%;
  }

  .label-text {
    font-size: 8pt;
    font-family: monospace;
    margin-top: 1mm;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
  }

  @media print {
    :global(body) {
      background: #ffffff !important;
      color: #000000 !important;
      margin: 0 !important;
      padding: 0 !important;
    }

    :global(.no-print) {
      display: none !important;
    }

    .print-page-wrapper {
      box-shadow: none !important;
      margin: 0 !important;
      padding: 5mm !important;
      max-width: 100% !important;
      width: 100% !important;
    }

    .label-cell {
      border: 1px solid #ddd !important;
      page-break-inside: avoid;
    }
  }
</style>
