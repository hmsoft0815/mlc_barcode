<script lang="ts">
  import { BARCODE_TYPES, type BarcodeType } from '../types';
  import {
    GenerateBatch,
    PickTextFile,
    PickExportFolder,
    ExportBatchToFolder,
    CopyToClipboard
  } from '../../../bindings/github.com/mlcmcp/mlc_barcode/barcodeapp';
  import type { BatchBarcodeResponse, BatchItemResult } from '../../../bindings/github.com/mlcmcp/mlc_barcode/models';

  export let onSendBatchToPrint: ((items: Array<{ data: string; svg: string; type: string }>) => void) | undefined = undefined;

  let selectedType: BarcodeType = 'qr';
  let rawText = 'MLC-PROD-001\nMLC-PROD-002\nMLC-PROD-003\nhttps://mlcgo.eu/demo1\nhttps://mlcgo.eu/demo2';
  let importedFilePath = '';

  // Options
  let fgColor = '#000000';
  let bgColor = '#ffffff';
  let isTransparent = false;
  let showText = false;

  // Batch states
  let isGenerating = false;
  let batchResponse: BatchBarcodeResponse | null = null;
  let isExporting = false;
  let exportMessage = '';
  let exportSuccess = false;

  // Export settings
  let exportFormat: 'svg' | 'png' = 'png';
  let namingScheme: 'index' | 'data_slug' | 'data_raw' = 'data_slug';
  let exportPrefix = 'barcode_';

  async function loadFile() {
    try {
      const [path, lines] = await PickTextFile();
      if (path && lines && lines.length > 0) {
        importedFilePath = path;
        rawText = lines.join('\n');
        runBatchGenerate();
      }
    } catch (e: any) {
      exportMessage = `Fehler beim Laden der Datei: ${e?.message}`;
      exportSuccess = false;
    }
  }

  async function runBatchGenerate() {
    const lines = rawText
      .split('\n')
      .map((l) => l.trim())
      .filter((l) => l.length > 0);

    if (lines.length === 0) {
      batchResponse = null;
      return;
    }

    isGenerating = true;
    exportMessage = '';
    try {
      const res = await GenerateBatch({
        type: selectedType,
        lines,
        width: 0,
        height: 0,
        showText,
        foregroundColor: fgColor,
        backgroundColor: isTransparent ? 'transparent' : bgColor
      });
      batchResponse = res;
    } catch (e: any) {
      exportMessage = `Fehler bei der Stapelgenerierung: ${e?.message}`;
      exportSuccess = false;
    } finally {
      isGenerating = false;
    }
  }

  async function startFolderExport() {
    if (!batchResponse?.items || batchResponse.items.length === 0) return;

    try {
      const folderPath = await PickExportFolder();
      if (!folderPath) return; // User cancelled

      isExporting = true;
      exportMessage = '';

      const res = await ExportBatchToFolder({
        folderPath,
        format: exportFormat,
        namingScheme,
        prefix: exportPrefix,
        items: batchResponse.items
      });

      if (res.error) {
        exportMessage = `Export fehlgeschlagen: ${res.error}`;
        exportSuccess = false;
      } else {
        exportMessage = `✓ ${res.exportedCount} Barcodes erfolgreich im Ordner gespeichert: ${folderPath}`;
        exportSuccess = true;
      }
    } catch (e: any) {
      exportMessage = `Export-Fehler: ${e?.message}`;
      exportSuccess = false;
    } finally {
      isExporting = false;
    }
  }

  function sendToPrintSheet() {
    if (!batchResponse?.items || !onSendBatchToPrint) return;
    const validItems = batchResponse.items
      .filter((it) => it.success && it.svg)
      .map((it) => ({
        data: it.data,
        svg: it.svg!,
        type: selectedType
      }));
    onSendBatchToPrint(validItems);
  }

  async function copyItemSVG(svg?: string) {
    if (!svg) return;
    await CopyToClipboard(svg);
  }
</script>

<div class="container-fluid py-3">
  {#if exportMessage}
    <div class="alert alert-{exportSuccess ? 'success' : 'danger'} alert-dismissible fade show py-2 px-3 mb-3 shadow-sm" role="alert">
      <i class="bi bi-{exportSuccess ? 'check-circle' : 'exclamation-triangle'} me-2"></i>
      {exportMessage}
    </div>
  {/if}

  <div class="row g-3">
    <!-- Left Column: Input & Controls -->
    <div class="col-lg-5">
      <div class="card shadow-sm border mb-3">
        <div class="card-header bg-body border-bottom py-2 d-flex justify-content-between align-items-center">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-file-earmark-text me-1 text-primary"></i> Datenquelle (1 Zeile = 1 Barcode)
          </h6>
          <button class="btn btn-outline-primary btn-sm" on:click={loadFile}>
            <i class="bi bi-folder2-open me-1"></i> Textdatei laden...
          </button>
        </div>
        <div class="card-body">
          {#if importedFilePath}
            <div class="small text-body-secondary mb-2 text-truncate">
              <i class="bi bi-file-check text-success me-1"></i> {importedFilePath}
            </div>
          {/if}

          <div class="mb-3">
            <textarea
              class="form-control font-monospace small"
              rows={8}
              placeholder="Zeile 1\nZeile 2\nZeile 3..."
              bind:value={rawText}
              on:input={runBatchGenerate}
            ></textarea>
            <div class="form-text small">
              Jede Zeile wird als separater Barcode generiert.
            </div>
          </div>

          <!-- Configuration -->
          <div class="row g-2 mb-3">
            <div class="col-6">
              <label for="batchTypeSelect" class="form-label small text-body-secondary mb-1">Barcode-Typ</label>
              <select
                id="batchTypeSelect"
                class="form-select form-select-sm"
                bind:value={selectedType}
                on:change={runBatchGenerate}
              >
                {#each BARCODE_TYPES as t}
                  <option value={t.id}>{t.name}</option>
                {/each}
              </select>
            </div>
            <div class="col-3">
              <label for="batchFgColor" class="form-label small text-body-secondary mb-1">Vordergrund</label>
              <input
                id="batchFgColor"
                type="color"
                class="form-control form-control-sm form-control-color w-100"
                bind:value={fgColor}
                on:input={runBatchGenerate}
                aria-label="Vordergrundfarbe"
              />
            </div>
            <div class="col-3">
              <label for="batchBgColor" class="form-label small text-body-secondary mb-1">Hintergrund</label>
              <input
                id="batchBgColor"
                type="color"
                class="form-control form-control-sm form-control-color w-100"
                disabled={isTransparent}
                bind:value={bgColor}
                on:input={runBatchGenerate}
                aria-label="Hintergrundfarbe"
              />
            </div>
          </div>

          <div class="row g-2 mb-3">
            <div class="col-6">
              <div class="form-check form-switch small">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="batchTransparent"
                  bind:checked={isTransparent}
                  on:change={runBatchGenerate}
                />
                <label class="form-check-label" for="batchTransparent">Transparent</label>
              </div>
            </div>
            <div class="col-6">
              <div class="form-check form-switch small">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="batchShowText"
                  bind:checked={showText}
                  on:change={runBatchGenerate}
                />
                <label class="form-check-label" for="batchShowText">Text anzeigen</label>
              </div>
            </div>
          </div>

          <button
            class="btn btn-primary btn-sm w-100"
            disabled={isGenerating}
            on:click={runBatchGenerate}
          >
            {#if isGenerating}
              <span class="spinner-border spinner-border-sm me-1"></span> Generiere...
            {:else}
              <i class="bi bi-arrow-repeat me-1"></i> Barcodes neu generieren
            {/if}
          </button>
        </div>
      </div>

      <!-- Export Settings Box -->
      <div class="card shadow-sm border">
        <div class="card-header bg-body border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-box-arrow-up-right me-1 text-primary"></i> Batch-Export Einstellungen
          </h6>
        </div>
        <div class="card-body">
          <div class="row g-2 mb-3">
            <div class="col-6">
              <label for="exportFormatSelect" class="form-label small text-body-secondary mb-1">Format</label>
              <select id="exportFormatSelect" class="form-select form-select-sm" bind:value={exportFormat}>
                <option value="png">PNG (Rastergrafik)</option>
                <option value="svg">SVG (Vektorgrafik)</option>
              </select>
            </div>
            <div class="col-6">
              <label for="namingSchemeSelect" class="form-label small text-body-secondary mb-1">Dateibenennung</label>
              <select id="namingSchemeSelect" class="form-select form-select-sm" bind:value={namingScheme}>
                <option value="data_slug">Nummer + Inhalt</option>
                <option value="index">Nur Nummer (001, 002...)</option>
                <option value="data_raw">Nur Inhalt</option>
              </select>
            </div>
          </div>

          <div class="mb-3">
            <label for="exportPrefixInput" class="form-label small text-body-secondary mb-1">Dateinamen-Präfix</label>
            <input
              id="exportPrefixInput"
              type="text"
              class="form-control form-control-sm font-monospace"
              placeholder="z.B. code_"
              bind:value={exportPrefix}
            />
          </div>

          <div class="d-flex gap-2">
            <button
              class="btn btn-success btn-sm flex-grow-1"
              disabled={isExporting || !batchResponse?.validCount}
              on:click={startFolderExport}
            >
              {#if isExporting}
                <span class="spinner-border spinner-border-sm me-1"></span> Exportiere...
              {:else}
                <i class="bi bi-folder-symlink me-1"></i> In Ordner exportieren...
              {/if}
            </button>

            {#if onSendBatchToPrint}
              <button
                class="btn btn-outline-secondary btn-sm"
                disabled={!batchResponse?.validCount}
                on:click={sendToPrintSheet}
              >
                <i class="bi bi-printer me-1"></i> Druckbogen
              </button>
            {/if}
          </div>
        </div>
      </div>
    </div>

    <!-- Right Column: Results Table -->
    <div class="col-lg-7">
      <div class="card shadow-sm border h-100 d-flex flex-column">
        <div class="card-header bg-body border-bottom py-2 d-flex justify-content-between align-items-center">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-list-check me-1 text-primary"></i> Generierte Barcodes
          </h6>
          {#if batchResponse}
            <div class="d-flex gap-2">
              <span class="badge bg-success-subtle text-success border border-success-subtle">
                {batchResponse.validCount} Gültig
              </span>
              {#if batchResponse.errorCount > 0}
                <span class="badge bg-danger-subtle text-danger border border-danger-subtle">
                  {batchResponse.errorCount} Fehler
                </span>
              {/if}
            </div>
          {/if}
        </div>

        <div class="card-body p-0 flex-grow-1 overflow-auto" style="max-height: 600px;">
          {#if batchResponse?.items && batchResponse.items.length > 0}
            <div class="table-responsive">
              <table class="table table-hover table-sm align-middle mb-0">
                <thead class="table-light sticky-top">
                  <tr>
                    <th style="width: 45px;">#</th>
                    <th style="width: 100px;">Vorschau</th>
                    <th>Inhalt</th>
                    <th style="width: 90px;">Status</th>
                    <th style="width: 60px;">Aktion</th>
                  </tr>
                </thead>
                <tbody>
                  {#each batchResponse.items as item (item.index)}
                    <tr>
                      <td class="text-body-secondary small">{item.index}</td>
                      <td>
                        {#if item.success && item.svg}
                          <div class="mini-svg-preview checkerboard-bg">
                            <!-- eslint-disable-next-line svelte/no-at-html-tags -->
                            {@html item.svg}
                          </div>
                        {:else}
                          <span class="text-danger small"><i class="bi bi-x-circle"></i></span>
                        {/if}
                      </td>
                      <td>
                        <span class="font-monospace small text-break">{item.data}</span>
                        {#if !item.success && item.error}
                          <div class="text-danger small">{item.error}</div>
                        {/if}
                      </td>
                      <td>
                        {#if item.success}
                          <span class="badge bg-success-subtle text-success">OK</span>
                        {:else}
                          <span class="badge bg-danger-subtle text-danger">Fehler</span>
                        {/if}
                      </td>
                      <td>
                        {#if item.success}
                          <button
                            class="btn btn-outline-secondary btn-sm p-1"
                            title="SVG kopieren"
                            on:click={() => copyItemSVG(item.svg)}
                          >
                            <i class="bi bi-clipboard"></i>
                          </button>
                        {/if}
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {:else}
            <div class="text-body-secondary text-center py-5">
              <i class="bi bi-collection fs-1 opacity-50 mb-2 d-block"></i>
              <span>Keine Barcode-Daten geladen</span>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .mini-svg-preview {
    width: 80px;
    height: 38px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--bs-border-color);
    border-radius: 4px;
    padding: 2px;
  }

  .mini-svg-preview :global(svg) {
    max-width: 100%;
    max-height: 100%;
  }
</style>
