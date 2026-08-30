<script lang="ts">
  export let printItems: Array<{ data: string; svg: string; type: string }> = [];

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

  // Expanded items based on repeatCount
  $: expandedItems = printItems.flatMap((item) =>
    Array(repeatCount).fill(item)
  );

  function triggerPrint() {
    window.print();
  }
</script>

<div class="container-fluid py-3">
  <!-- Controls (hidden when printing) -->
  <div class="card shadow-sm border-0 mb-3 no-print">
    <div class="card-header bg-white border-bottom py-2 d-flex justify-content-between align-items-center">
      <h6 class="mb-0 fw-semibold text-secondary">
        <i class="bi bi-printer me-1"></i> Etikettenbogen-Layout & Druckeinstellungen
      </h6>
      <button
        class="btn btn-primary btn-sm"
        disabled={expandedItems.length === 0}
        on:click={triggerPrint}
      >
        <i class="bi bi-printer-fill me-1"></i> Drucken (Print Dialog)
      </button>
    </div>
    <div class="card-body">
      <div class="row g-3 align-items-center">
        <div class="col-md-4">
          <label class="form-label small text-muted mb-1">Standard-Vorlagen (DIN A4)</label>
          <select
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
          <label class="form-label small text-muted mb-1">Spalten (Columns)</label>
          <input
            type="number"
            class="form-control form-control-sm"
            min="1"
            max="6"
            bind:value={columns}
          />
        </div>

        <div class="col-md-2 col-6">
          <label class="form-label small text-muted mb-1">Zeilen (Rows)</label>
          <input
            type="number"
            class="form-control form-control-sm"
            min="1"
            max="15"
            bind:value={rows}
          />
        </div>

        <div class="col-md-2 col-6">
          <label class="form-label small text-muted mb-1">Wiederholungen</label>
          <input
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
    </div>
  </div>

  <!-- Print Sheet Container -->
  {#if expandedItems.length > 0}
    <div class="print-page-wrapper">
      <div
        class="print-grid"
        style="--cols: {columns}; --rows: {rows};"
      >
        {#each expandedItems as item, idx}
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
    <div class="card shadow-sm border-0 py-5 text-center text-muted no-print">
      <i class="bi bi-printer fs-1 text-secondary opacity-50 mb-2"></i>
      <h5>Keine Etiketten in der Druck-Warteschlange</h5>
      <p class="small text-muted mb-0">
        Erstelle Barcodes im <strong>Einzel-Generator</strong> oder <strong>Batch-Generator</strong> und klicke auf "Als Etikett drucken".
      </p>
    </div>
  {/if}
</div>

<style>
  .print-page-wrapper {
    background: #fff;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
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
    background: #fff;
  }

  .label-svg-wrapper {
    max-width: 90%;
    max-height: 20mm;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .label-svg-wrapper :global(svg) {
    max-width: 100%;
    max-height: 100%;
    height: auto;
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
      background: #fff !important;
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
