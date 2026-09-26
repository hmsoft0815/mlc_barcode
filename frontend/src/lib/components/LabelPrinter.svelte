<script lang="ts">
  import { onMount } from 'svelte';
  import { GenerateBatch, PickTableFile } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import { BARCODE_TYPES, type BarcodeType, type PrintItem } from '../types';
  import { t, type LabelsKey } from '../i18n/text/labels';

  export let printItems: PrintItem[] = [];
  export let batchItems: PrintItem[] = [];

  // File import: first column = code, optional second column = label text.
  let importType: BarcodeType = 'qr';
  let skipHeader = false;
  // Message parts as key + params, so they follow a language switch.
  type Msg = { key: LabelsKey; params?: Record<string, string | number> };
  let importMessage: Msg[] = [];
  let importOk = true;

  function takeFromBatch() {
    printItems = [...batchItems];
    importMessage = [{ key: 'takenFromBatch', params: { n: batchItems.length } }];
    importOk = true;
  }

  async function importFile() {
    try {
      const file = await PickTableFile();
      if (!file.path) return;
      const all = (file.rows ?? []).map((r) => r ?? []);
      const rows = skipHeader ? all.slice(1) : all;
      if (rows.length === 0) {
        importMessage = [{ key: 'fileEmpty' }];
        importOk = false;
        return;
      }

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

      const name = file.path.split(/[\\/]/).pop();
      importMessage = [{ key: 'imported', params: { n: printItems.length, name: `${name}` } }];
      if (res.errorCount)
        importMessage = [...importMessage, { key: 'importedInvalid', params: { n: res.errorCount, type: importType.toUpperCase() } }];
      importOk = !res.errorCount;
    } catch (e: any) {
      importMessage = [{ key: 'importFailed', params: { msg: `${e?.message ?? e}` } }];
      importOk = false;
    }
  }

  let columns = 3;
  let rows = 8;
  let repeatCount = 1;
  let showDataText = true;

  const PRESETS = [
    { key: 'presetAvery2x4', cols: 2, rows: 4 },
    { key: 'presetAvery3x7', cols: 3, rows: 7 },
    { key: 'presetAvery3x8', cols: 3, rows: 8 },
    { key: 'presetAvery4x10', cols: 4, rows: 10 },
    { key: 'presetSingle', cols: 1, rows: 1 }
  ] as const;

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
        <i class="bi bi-printer me-1 text-primary"></i> {$t('title')}
      </h6>
      <div class="d-flex gap-2">
        <button
          class="btn btn-outline-primary btn-sm"
          disabled={batchItems.length === 0}
          title={batchItems.length ? $t('batchTitle', { n: batchItems.length }) : $t('batchEmpty')}
          on:click={takeFromBatch}
        >
          <i class="bi bi-collection me-1"></i> {$t('takeFromBatch')}
        </button>
        <button class="btn btn-outline-secondary btn-sm" on:click={loadSampleLabels}>
          <i class="bi bi-magic me-1"></i> {$t('loadSample')}
        </button>
        <button
          class="btn btn-primary btn-sm"
          disabled={expandedItems.length === 0}
          on:click={triggerPrint}
        >
          <i class="bi bi-printer-fill me-1"></i> {$t('print')}
        </button>
      </div>
    </div>
    <div class="card-body">
      <div class="row g-3 align-items-center">
        <div class="col-md-4">
          <label for="labelPresetSelect" class="form-label small text-body-secondary mb-1">{$t('presets')}</label>
          <select
            id="labelPresetSelect"
            class="form-select form-select-sm"
            on:change={(e) => {
              const val = e.currentTarget.value;
              const found = PRESETS.find((p) => p.key === val);
              if (found) applyPreset(found);
            }}
          >
            {#each PRESETS as p}
              <option value={p.key} selected={p.cols === columns && p.rows === rows}>
                {$t(p.key)}
              </option>
            {/each}
          </select>
        </div>

        <div class="col-md-2 col-6">
          <label for="labelColumnsInput" class="form-label small text-body-secondary mb-1">{$t('columns')}</label>
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
          <label for="labelRowsInput" class="form-label small text-body-secondary mb-1">{$t('rows')}</label>
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
          <label for="labelRepeatInput" class="form-label small text-body-secondary mb-1">{$t('repeat')}</label>
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
            <label class="form-check-label" for="showDataLabel">{$t('showText')}</label>
          </div>
        </div>
      </div>

      <div class="row g-2 align-items-end border-top pt-3 mt-2">
        <div class="col-md-4">
          <label for="labelImportType" class="form-label small text-body-secondary mb-1">
            {$t('importLabel')}
          </label>
          <select id="labelImportType" class="form-select form-select-sm" bind:value={importType}>
            {#each BARCODE_TYPES as bt}
              <option value={bt.id}>{bt.name}</option>
            {/each}
          </select>
        </div>
        <div class="col-md-3 col-6">
          <div class="form-check form-switch small mb-1">
            <input class="form-check-input" type="checkbox" id="labelSkipHeader" bind:checked={skipHeader} />
            <label class="form-check-label" for="labelSkipHeader">{$t('skipHeader')}</label>
          </div>
        </div>
        <div class="col-md-2 col-6">
          <button class="btn btn-outline-primary btn-sm w-100" on:click={importFile}>
            <i class="bi bi-filetype-csv me-1"></i> {$t('importButton')}
          </button>
        </div>
        {#if importMessage.length}
          <div class="col-12 small {importOk ? 'text-success' : 'text-warning'}">{importMessage.map((m) => $t(m.key, m.params)).join('')}</div>
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
      <h5 class="text-body">{$t('emptyTitle')}</h5>
      <p class="small text-body-secondary mb-3">
        {$t('emptyHintBefore')} <strong>{$t('single')}</strong> {$t('or')} <strong>{$t('batch')}</strong> {$t('emptyHintAfter')}
      </p>
      <div>
        <button class="btn btn-outline-primary btn-sm" on:click={loadSampleLabels}>
          <i class="bi bi-magic me-1"></i> {$t('loadExamples')}
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
      padding: 10mm !important;
      max-width: 100% !important;
      width: 100% !important;
    }

    .label-cell {
      border: 1px solid #ddd !important;
      page-break-inside: avoid;
    }
  }
</style>
