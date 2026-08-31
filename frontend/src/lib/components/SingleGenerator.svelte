<script lang="ts">
  import { onMount } from 'svelte';
  import { BARCODE_TYPES, type BarcodeType } from '../types';
  import {
    GenerateBarcode,
    FormatWifi,
    FormatVCard,
    FormatEvent,
    SaveSingleFile,
    CopyToClipboard
  } from '../../../bindings/github.com/mlcmcp/mlc_barcode/barcodeapp';
  import type { BarcodeResult } from '../../../bindings/github.com/mlcmcp/mlc_barcode/models';

  export let onSendToPrint: ((item: { data: string; svg: string; type: string }) => void) | undefined = undefined;

  let selectedType: BarcodeType = 'qr';
  let qrMode: 'text' | 'wifi' | 'vcard' | 'event' = 'text';

  // Free text / data
  let rawData: string = 'https://mlcgo.eu';

  // Structured QR inputs
  let wifiSSID = '';
  let wifiPass = '';
  let wifiEnc = 'WPA';
  let wifiHidden = false;

  let vcardFirst = '';
  let vcardLast = '';
  let vcardEmail = '';
  let vcardPhone = '';

  let eventSummary = '';
  let eventStart = '';
  let eventEnd = '';
  let eventTZ = 'Europe/Berlin';

  // Styling options
  let fgColor = '#000000';
  let bgColor = '#ffffff';
  let isTransparent = false;
  let showText = false;
  let customWidth = 0;
  let customHeight = 0;

  // Result state
  let result: BarcodeResult | null = null;
  let isGenerating = false;
  let feedbackMessage = '';
  let feedbackType: 'success' | 'danger' = 'success';
  let feedbackTimer: ReturnType<typeof setTimeout> | null = null;

  function showFeedback(msg: string, type: 'success' | 'danger' = 'success') {
    feedbackMessage = msg;
    feedbackType = type;
    if (feedbackTimer) clearTimeout(feedbackTimer);
    feedbackTimer = setTimeout(() => {
      feedbackMessage = '';
    }, 3000);
  }

  async function updateStructuredQR() {
    if (selectedType !== 'qr') return;

    if (qrMode === 'wifi') {
      if (!wifiSSID) return;
      rawData = await FormatWifi({
        ssid: wifiSSID,
        password: wifiPass,
        encryption: wifiEnc,
        hidden: wifiHidden
      });
    } else if (qrMode === 'vcard') {
      if (!vcardFirst && !vcardLast) return;
      rawData = await FormatVCard({
        firstName: vcardFirst,
        lastName: vcardLast,
        email: vcardEmail,
        phone: vcardPhone
      });
    } else if (qrMode === 'event') {
      if (!eventSummary) return;
      rawData = await FormatEvent({
        summary: eventSummary,
        startTime: eventStart.replace(/[-:]/g, ''),
        endTime: eventEnd.replace(/[-:]/g, ''),
        timeZone: eventTZ
      });
    }
    triggerGenerate();
  }

  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  function triggerGenerate() {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      generate();
    }, 120);
  }

  async function generate() {
    if (!rawData.trim()) {
      result = null;
      return;
    }

    isGenerating = true;
    try {
      const res = await GenerateBarcode({
        type: selectedType,
        data: rawData.trim(),
        width: customWidth > 0 ? Number(customWidth) : 0,
        height: customHeight > 0 ? Number(customHeight) : 0,
        showText,
        foregroundColor: fgColor,
        backgroundColor: isTransparent ? 'transparent' : bgColor
      });
      result = res;
    } catch (e: any) {
      result = {
        type: selectedType,
        data: rawData,
        success: false,
        error: e?.message || 'Unbekannter Fehler bei der Generierung'
      };
    } finally {
      isGenerating = false;
    }
  }

  function handleTypeChange(newType: BarcodeType) {
    selectedType = newType;
    const opt = BARCODE_TYPES.find((t) => t.id === newType);
    if (opt && rawData === 'https://mlcgo.eu' && newType !== 'qr') {
      rawData = opt.sample;
    }
    triggerGenerate();
  }

  async function copySVG() {
    if (!result?.svg) return;
    const ok = await CopyToClipboard(result.svg);
    if (ok) {
      showFeedback('SVG erfolgreich in Zwischenablage kopiert!');
    } else {
      await navigator.clipboard.writeText(result.svg);
      showFeedback('SVG in Zwischenablage kopiert!');
    }
  }

  async function copyPNG() {
    if (!result?.pngData) return;
    try {
      const res = await fetch(result.pngData);
      const blob = await res.blob();
      await navigator.clipboard.write([
        new ClipboardItem({ 'image/png': blob })
      ]);
      showFeedback('PNG-Grafik in Zwischenablage kopiert!');
    } catch (err) {
      showFeedback('Kopieren als Bild fehlgeschlagen', 'danger');
    }
  }

  async function saveSVG() {
    if (!result?.svg) return;
    try {
      const savedPath = await SaveSingleFile({
        defaultName: `${selectedType}_barcode.svg`,
        format: 'svg',
        content: result.svg
      });
      if (savedPath) {
        showFeedback(`Gespeichert unter: ${savedPath}`);
      }
    } catch (e: any) {
      showFeedback(`Fehler beim Speichern: ${e?.message}`, 'danger');
    }
  }

  async function savePNG() {
    if (!result?.pngData) return;
    try {
      const savedPath = await SaveSingleFile({
        defaultName: `${selectedType}_barcode.png`,
        format: 'png',
        content: result.pngData
      });
      if (savedPath) {
        showFeedback(`Gespeichert unter: ${savedPath}`);
      }
    } catch (e: any) {
      showFeedback(`Fehler beim Speichern: ${e?.message}`, 'danger');
    }
  }

  function sendToLabelPrint() {
    if (result?.svg && onSendToPrint) {
      onSendToPrint({
        data: rawData,
        svg: result.svg,
        type: selectedType
      });
    }
  }

  onMount(() => {
    generate();
  });
</script>

<div class="container-fluid py-3">
  {#if feedbackMessage}
    <div class="alert alert-{feedbackType} alert-dismissible fade show py-2 px-3 mb-3 shadow-sm" role="alert">
      <i class="bi bi-{feedbackType === 'success' ? 'check-circle' : 'exclamation-triangle'} me-2"></i>
      {feedbackMessage}
    </div>
  {/if}

  <div class="row g-3">
    <!-- Left Column: Settings & Input -->
    <div class="col-lg-6">
      <div class="card shadow-sm border mb-3">
        <div class="card-header bg-body border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-sliders me-1 text-primary"></i> Barcode-Typ & Konfiguration
          </h6>
        </div>
        <div class="card-body">
          <!-- Type Selection -->
          <div class="mb-3">
            <label for="singleTypeSelect" class="form-label fw-medium small text-body-secondary">Symbologie / Typ</label>
            <select
              id="singleTypeSelect"
              class="form-select form-select-sm"
              value={selectedType}
              on:change={(e) => handleTypeChange(e.currentTarget.value as BarcodeType)}
            >
              <optgroup label="2D Matrix">
                {#each BARCODE_TYPES.filter((t) => t.category === '2D Matrix') as t}
                  <option value={t.id}>{t.name} ({t.description})</option>
                {/each}
              </optgroup>
              <optgroup label="1D Linear">
                {#each BARCODE_TYPES.filter((t) => t.category === '1D Linear') as t}
                  <option value={t.id}>{t.name} - {t.description}</option>
                {/each}
              </optgroup>
            </select>
          </div>

          <!-- QR Special Modes (if QR selected) -->
          {#if selectedType === 'qr'}
            <div class="mb-3">
              <span class="form-label fw-medium small text-body-secondary d-block">QR-Inhaltsformat</span>
              <ul class="nav nav-pills nav-fill bg-body-secondary p-1 rounded small">
                <li class="nav-item">
                  <button
                    class="nav-link py-1 {qrMode === 'text' ? 'active' : 'text-body'}"
                    type="button"
                    on:click={() => {
                      qrMode = 'text';
                      triggerGenerate();
                    }}
                  >
                    <i class="bi bi-link-45deg"></i> Text / URL
                  </button>
                </li>
                <li class="nav-item">
                  <button
                    class="nav-link py-1 {qrMode === 'wifi' ? 'active' : 'text-body'}"
                    type="button"
                    on:click={() => {
                      qrMode = 'wifi';
                      updateStructuredQR();
                    }}
                  >
                    <i class="bi bi-wifi"></i> WLAN
                  </button>
                </li>
                <li class="nav-item">
                  <button
                    class="nav-link py-1 {qrMode === 'vcard' ? 'active' : 'text-body'}"
                    type="button"
                    on:click={() => {
                      qrMode = 'vcard';
                      updateStructuredQR();
                    }}
                  >
                    <i class="bi bi-person-badge"></i> vCard
                  </button>
                </li>
                <li class="nav-item">
                  <button
                    class="nav-link py-1 {qrMode === 'event' ? 'active' : 'text-body'}"
                    type="button"
                    on:click={() => {
                      qrMode = 'event';
                      updateStructuredQR();
                    }}
                  >
                    <i class="bi bi-calendar-event"></i> Termin
                  </button>
                </li>
              </ul>
            </div>

            <!-- Structured Forms for QR -->
            {#if qrMode === 'wifi'}
              <div class="p-3 bg-body-secondary rounded mb-3 border">
                <div class="row g-2 mb-2">
                  <div class="col-8">
                    <label for="wifiSsidInput" class="form-label small mb-1">WLAN-Name (SSID)</label>
                    <input
                      id="wifiSsidInput"
                      type="text"
                      class="form-control form-control-sm"
                      placeholder="z.B. HomeNetwork"
                      bind:value={wifiSSID}
                      on:input={updateStructuredQR}
                    />
                  </div>
                  <div class="col-4">
                    <label for="wifiEncSelect" class="form-label small mb-1">Verschlüsselung</label>
                    <select
                      id="wifiEncSelect"
                      class="form-select form-select-sm"
                      bind:value={wifiEnc}
                      on:change={updateStructuredQR}
                    >
                      <option value="WPA">WPA/WPA2/WPA3</option>
                      <option value="WEP">WEP</option>
                      <option value="nopass">Kein Passwort</option>
                    </select>
                  </div>
                </div>
                {#if wifiEnc !== 'nopass'}
                  <div class="mb-2">
                    <label for="wifiPassInput" class="form-label small mb-1">WLAN-Passwort</label>
                    <input
                      id="wifiPassInput"
                      type="password"
                      class="form-control form-control-sm"
                      placeholder="Passwort"
                      bind:value={wifiPass}
                      on:input={updateStructuredQR}
                    />
                  </div>
                {/if}
              </div>
            {:else if qrMode === 'vcard'}
              <div class="p-3 bg-body-secondary rounded mb-3 border">
                <div class="row g-2 mb-2">
                  <div class="col-6">
                    <label for="vcardFirstInput" class="form-label small mb-1">Vorname</label>
                    <input
                      id="vcardFirstInput"
                      type="text"
                      class="form-control form-control-sm"
                      placeholder="Max"
                      bind:value={vcardFirst}
                      on:input={updateStructuredQR}
                    />
                  </div>
                  <div class="col-6">
                    <label for="vcardLastInput" class="form-label small mb-1">Nachname</label>
                    <input
                      id="vcardLastInput"
                      type="text"
                      class="form-control form-control-sm"
                      placeholder="Mustermann"
                      bind:value={vcardLast}
                      on:input={updateStructuredQR}
                    />
                  </div>
                </div>
                <div class="row g-2">
                  <div class="col-6">
                    <label for="vcardEmailInput" class="form-label small mb-1">E-Mail</label>
                    <input
                      id="vcardEmailInput"
                      type="email"
                      class="form-control form-control-sm"
                      placeholder="max@beispiel.de"
                      bind:value={vcardEmail}
                      on:input={updateStructuredQR}
                    />
                  </div>
                  <div class="col-6">
                    <label for="vcardPhoneInput" class="form-label small mb-1">Telefon</label>
                    <input
                      id="vcardPhoneInput"
                      type="tel"
                      class="form-control form-control-sm"
                      placeholder="+49 123 456789"
                      bind:value={vcardPhone}
                      on:input={updateStructuredQR}
                    />
                  </div>
                </div>
              </div>
            {:else if qrMode === 'event'}
              <div class="p-3 bg-body-secondary rounded mb-3 border">
                <div class="mb-2">
                  <label for="eventSummaryInput" class="form-label small mb-1">Event-Titel</label>
                  <input
                    id="eventSummaryInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="Team Meeting"
                    bind:value={eventSummary}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="row g-2">
                  <div class="col-6">
                    <label for="eventStartInput" class="form-label small mb-1">Startzeit (YYYYMMDDTHHMMSS)</label>
                    <input
                      id="eventStartInput"
                      type="text"
                      class="form-control form-control-sm"
                      placeholder="20260901T100000"
                      bind:value={eventStart}
                      on:input={updateStructuredQR}
                    />
                  </div>
                  <div class="col-6">
                    <label for="eventEndInput" class="form-label small mb-1">Endzeit (YYYYMMDDTHHMMSS)</label>
                    <input
                      id="eventEndInput"
                      type="text"
                      class="form-control form-control-sm"
                      placeholder="20260901T113000"
                      bind:value={eventEnd}
                      on:input={updateStructuredQR}
                    />
                  </div>
                </div>
              </div>
            {/if}
          {/if}

          <!-- Direct Data Input -->
          <div class="mb-3">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label for="rawDataTextarea" class="form-label fw-medium small text-body-secondary mb-0">Daten / Inhalt</label>
              {#if selectedType === 'ean13'}
                <span class="badge bg-secondary-subtle text-secondary small">12 oder 13 Ziffern</span>
              {:else if selectedType === 'ean8'}
                <span class="badge bg-secondary-subtle text-secondary small">7 oder 8 Ziffern</span>
              {:else if selectedType === 'upca'}
                <span class="badge bg-secondary-subtle text-secondary small">11 oder 12 Ziffern</span>
              {/if}
            </div>
            <textarea
              id="rawDataTextarea"
              class="form-control form-control-sm font-monospace"
              rows={qrMode === 'text' ? 3 : 2}
              placeholder="Zu codierender Text oder Ziffern..."
              bind:value={rawData}
              on:input={triggerGenerate}
            ></textarea>
          </div>

          <!-- Color & Appearance Options -->
          <div class="row g-2 mb-3">
            <div class="col-sm-6">
              <label for="singleFgColorText" class="form-label small text-body-secondary mb-1">Vordergrundfarbe</label>
              <div class="input-group input-group-sm">
                <input
                  type="color"
                  class="form-control form-control-color"
                  bind:value={fgColor}
                  on:input={triggerGenerate}
                  aria-label="Vordergrund-Farbwähler"
                />
                <input
                  id="singleFgColorText"
                  type="text"
                  class="form-control font-monospace"
                  bind:value={fgColor}
                  on:input={triggerGenerate}
                />
              </div>
            </div>

            <div class="col-sm-6">
              <label for="singleBgColorText" class="form-label small text-body-secondary mb-1">Hintergrundfarbe</label>
              <div class="input-group input-group-sm">
                <input
                  type="color"
                  class="form-control form-control-color"
                  disabled={isTransparent}
                  bind:value={bgColor}
                  on:input={triggerGenerate}
                  aria-label="Hintergrund-Farbwähler"
                />
                <input
                  id="singleBgColorText"
                  type="text"
                  class="form-control font-monospace"
                  disabled={isTransparent}
                  bind:value={bgColor}
                  on:input={triggerGenerate}
                />
              </div>
            </div>
          </div>

          <div class="row g-2 mb-2">
            <div class="col-6">
              <div class="form-check form-switch small">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="transparentCheck"
                  bind:checked={isTransparent}
                  on:change={triggerGenerate}
                />
                <label class="form-check-label" for="transparentCheck">Transparenter Hintergrund</label>
              </div>
            </div>
            <div class="col-6">
              <div class="form-check form-switch small">
                <input
                  class="form-check-input"
                  type="checkbox"
                  id="showTextCheck"
                  bind:checked={showText}
                  on:change={triggerGenerate}
                />
                <label class="form-check-label" for="showTextCheck">Klartext unter Barcode</label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Column: Live Preview & Export -->
    <div class="col-lg-6">
      <div class="card shadow-sm border h-100">
        <div class="card-header bg-body border-bottom py-2 d-flex justify-content-between align-items-center">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-eye me-1 text-primary"></i> Live-Vorschau
          </h6>
          {#if isGenerating}
            <span class="spinner-border spinner-border-sm text-primary" role="status"></span>
          {/if}
        </div>
        <div class="card-body d-flex flex-column align-items-center justify-content-center p-4">
          {#if result?.success && result.svg}
            <div class="preview-box checkerboard-bg p-3 rounded d-flex align-items-center justify-content-center w-100 mb-3 border">
              <div class="svg-container">
                <!-- eslint-disable-next-line svelte/no-at-html-tags -->
                {@html result.svg}
              </div>
            </div>
          {:else if result && !result.success}
            <div class="alert alert-danger w-100 text-center py-4 my-auto">
              <i class="bi bi-exclamation-octagon fs-2 text-danger mb-2 d-block"></i>
              <div class="fw-semibold mb-1">Ungültige Barcode-Daten</div>
              <div class="small text-danger-emphasis">{result.error}</div>
            </div>
          {:else}
            <div class="text-body-secondary text-center py-5">
              <i class="bi bi-qr-code fs-1 opacity-50 mb-2 d-block"></i>
              <span>Keine Daten eingegeben</span>
            </div>
          {/if}

          <!-- Action Buttons -->
          {#if result?.success}
            <div class="w-100 mt-auto pt-3 border-top">
              <div class="row g-2">
                <div class="col-6 col-md-3">
                  <button class="btn btn-outline-primary btn-sm w-100" on:click={copySVG}>
                    <i class="bi bi-clipboard me-1"></i> SVG Copy
                  </button>
                </div>
                <div class="col-6 col-md-3">
                  <button class="btn btn-outline-primary btn-sm w-100" on:click={copyPNG}>
                    <i class="bi bi-file-earmark-image me-1"></i> PNG Copy
                  </button>
                </div>
                <div class="col-6 col-md-3">
                  <button class="btn btn-primary btn-sm w-100" on:click={saveSVG}>
                    <i class="bi bi-download me-1"></i> SVG Save
                  </button>
                </div>
                <div class="col-6 col-md-3">
                  <button class="btn btn-primary btn-sm w-100" on:click={savePNG}>
                    <i class="bi bi-download me-1"></i> PNG Save
                  </button>
                </div>
              </div>

              {#if onSendToPrint}
                <div class="mt-2 text-center">
                  <button class="btn btn-outline-secondary btn-sm" on:click={sendToLabelPrint}>
                    <i class="bi bi-printer me-1"></i> Als Etikett drucken
                  </button>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .preview-box {
    min-height: 280px;
    max-height: 400px;
    overflow: auto;
  }

  .svg-container {
    max-width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .svg-container :global(svg) {
    max-width: 100%;
    max-height: 260px;
    height: auto;
    filter: drop-shadow(0 2px 4px rgba(0,0,0,0.12));
  }
</style>
