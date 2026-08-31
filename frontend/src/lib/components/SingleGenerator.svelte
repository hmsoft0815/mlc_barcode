<script lang="ts">
  import { onMount } from 'svelte';
  import { BARCODE_TYPES, type BarcodeType } from '../types';
  import {
    GenerateBarcode,
    FormatWifi,
    FormatVCard,
    FormatEvent,
    FormatEPC,
    FormatCrypto,
    FormatGeo,
    FormatTel,
    FormatSMS,
    FormatEmail,
    SaveSingleFile,
    CopyToClipboard
  } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import type { BarcodeResult } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';

  export let onSendToPrint: ((item: { data: string; svg: string; type: string }) => void) | undefined = undefined;

  let selectedType: BarcodeType = 'qr';
  type QRMode = 'text' | 'epc' | 'wifi' | 'vcard' | 'event' | 'crypto' | 'geo' | 'tel' | 'sms' | 'email';
  let qrMode: QRMode = 'text';

  // Free text / data
  let rawData: string = 'https://mlcgo.eu';
  let customLabelText: string = '';
  let customLabelTouched: boolean = false;

  // Structured QR inputs: EPC / GiroCode
  let epcName = 'Michael Lechner';
  let epcIBAN = 'DE89370400440532013000';
  let epcBIC = '';
  let epcAmount: number | string = 19.99;
  let epcRef = 'Rechnung-1002';

  // Structured QR inputs: WIFI
  let wifiSSID = '';
  let wifiPass = '';
  let wifiEnc = 'WPA';
  let wifiHidden = false;

  // Structured QR inputs: vCard
  let vcardFirst = '';
  let vcardLast = '';
  let vcardEmail = '';
  let vcardPhone = '';

  // Structured QR inputs: Event
  let eventSummary = '';
  let eventStart = '';
  let eventEnd = '';
  let eventTZ = 'Europe/Berlin';

  // Structured QR inputs: Crypto
  let cryptoCoin = 'bitcoin';
  let cryptoAddress = '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa';
  let cryptoAmount: number | string = 0.005;
  let cryptoMessage = 'Spende';

  // Structured QR inputs: Geo
  let geoLat = 52.5200;
  let geoLon = 13.4050;
  let geoQuery = 'Berlin Fernsehturm';

  // Structured QR inputs: Tel & SMS
  let telNumber = '+49 170 1234567';
  let smsNumber = '+49 170 1234567';
  let smsMessage = 'Hallo, ich interessiere mich für MLC Barcode!';

  // Structured QR inputs: Email
  let mailTo = 'support@mlcgo.eu';
  let mailSubject = 'Anfrage Barcode Software';
  let mailBody = 'Guten Tag,\n\nich habe eine Frage zum Produkt.';

  // Styling options
  let fgColor = '#000000';
  let bgColor = '#ffffff';
  let isTransparent = false;
  let showText = false;
  let customWidth = 0;
  let customHeight = 0;

  // Preset Color Palettes
  const PRESET_BG_COLORS = [
    { label: 'Weiß', value: '#ffffff' },
    { label: 'Gelb', value: '#fff59d' },
    { label: 'Hellblau', value: '#e1f5fe' },
    { label: 'Hellgrün', value: '#e8f5e9' },
    { label: 'Hellgrau', value: '#f8f9fa' }
  ];

  const PRESET_FG_COLORS = [
    { label: 'Schwarz', value: '#000000' },
    { label: 'Dunkelblau', value: '#0d47a1' },
    { label: 'Dunkelrot', value: '#b71c1c' },
    { label: 'Dunkelgrün', value: '#1b5e20' }
  ];

  function setBgColor(color: string) {
    bgColor = color;
    isTransparent = false;
    triggerGenerate();
  }

  function setFgColor(color: string) {
    fgColor = color;
    triggerGenerate();
  }

  function toggleTransparent() {
    isTransparent = !isTransparent;
    triggerGenerate();
  }

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
    if (selectedType !== 'qr' && selectedType !== 'datamatrix') {
      selectedType = 'qr';
    }

    if (qrMode === 'epc') {
      if (!epcIBAN && !epcName) return;
      const numAmount = typeof epcAmount === 'string' ? parseFloat(epcAmount) || 0 : (epcAmount || 0);
      rawData = await FormatEPC({
        name: epcName,
        iban: epcIBAN,
        bic: epcBIC,
        amount: numAmount,
        reference: epcRef,
        purpose: ''
      });
      if (!customLabelTouched) {
        customLabelText = numAmount > 0 ? `Überweisung: ${numAmount.toFixed(2)} € an ${epcName}` : `GiroCode: ${epcName}`;
      }
    } else if (qrMode === 'crypto') {
      if (!cryptoAddress) return;
      const numAmount = typeof cryptoAmount === 'string' ? parseFloat(cryptoAmount) || 0 : (cryptoAmount || 0);
      rawData = await FormatCrypto({
        coin: cryptoCoin,
        address: cryptoAddress,
        amount: numAmount,
        label: '',
        message: cryptoMessage
      });
      if (!customLabelTouched) {
        customLabelText = `${cryptoCoin.toUpperCase()}: ${cryptoAddress.slice(0, 8)}...${cryptoAddress.slice(-6)}`;
      }
    } else if (qrMode === 'geo') {
      rawData = await FormatGeo({
        latitude: Number(geoLat) || 0,
        longitude: Number(geoLon) || 0,
        query: geoQuery
      });
      if (!customLabelTouched) {
        customLabelText = geoQuery ? geoQuery : `Maps: ${geoLat}, ${geoLon}`;
      }
    } else if (qrMode === 'tel') {
      if (!telNumber) return;
      rawData = await FormatTel({
        phoneNumber: telNumber
      });
      if (!customLabelTouched) {
        customLabelText = `Tel: ${telNumber}`;
      }
    } else if (qrMode === 'sms') {
      if (!smsNumber) return;
      rawData = await FormatSMS({
        phoneNumber: smsNumber,
        message: smsMessage
      });
      if (!customLabelTouched) {
        customLabelText = `SMS: ${smsNumber}`;
      }
    } else if (qrMode === 'email') {
      if (!mailTo) return;
      rawData = await FormatEmail({
        to: mailTo,
        subject: mailSubject,
        body: mailBody
      });
      if (!customLabelTouched) {
        customLabelText = `E-Mail: ${mailTo}`;
      }
    } else if (qrMode === 'wifi') {
      if (!wifiSSID) return;
      rawData = await FormatWifi({
        ssid: wifiSSID,
        password: wifiPass,
        encryption: wifiEnc,
        hidden: wifiHidden
      });
      if (!customLabelTouched) {
        customLabelText = wifiSSID ? `WLAN: ${wifiSSID}` : '';
      }
    } else if (qrMode === 'vcard') {
      if (!vcardFirst && !vcardLast) return;
      rawData = await FormatVCard({
        firstName: vcardFirst,
        lastName: vcardLast,
        email: vcardEmail,
        phone: vcardPhone
      });
      if (!customLabelTouched) {
        customLabelText = `${vcardFirst} ${vcardLast}`.trim();
      }
    } else if (qrMode === 'event') {
      if (!eventSummary) return;
      rawData = await FormatEvent({
        summary: eventSummary,
        startTime: eventStart.replace(/[-:]/g, ''),
        endTime: eventEnd.replace(/[-:]/g, ''),
        timeZone: eventTZ
      });
      if (!customLabelTouched) {
        customLabelText = eventSummary;
      }
    }
    triggerGenerate();
  }

  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  function triggerGenerate() {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      generate();
    }, 100);
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
        customText: showText ? customLabelText.trim() : '',
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

  function handleModeChange(mode: QRMode) {
    qrMode = mode;
    customLabelTouched = false;
    if (mode !== 'text') {
      if (selectedType !== 'qr' && selectedType !== 'datamatrix') {
        selectedType = 'qr';
      }
      updateStructuredQR();
    } else {
      triggerGenerate();
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
        data: customLabelText.trim() || rawData,
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
            <i class="bi bi-sliders me-1 text-primary"></i> Barcode-Inhalt & Format
          </h6>
        </div>
        <div class="card-body">
          <!-- Inhaltsformat / Typ-Modus Switcher -->
          <div class="mb-3">
            <span class="form-label fw-medium small text-body-secondary d-block mb-1">Inhaltstyp / Vorlage</span>
            <div class="d-flex flex-wrap gap-1 bg-body-secondary p-1 rounded">
              <button
                class="btn btn-sm {qrMode === 'text' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('text')}
              >
                <i class="bi bi-fonts"></i> Freitext/URL
              </button>
              <button
                class="btn btn-sm {qrMode === 'epc' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('epc')}
              >
                <i class="bi bi-bank"></i> GiroCode / SEPA
              </button>
              <button
                class="btn btn-sm {qrMode === 'wifi' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('wifi')}
              >
                <i class="bi bi-wifi"></i> WLAN
              </button>
              <button
                class="btn btn-sm {qrMode === 'vcard' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('vcard')}
              >
                <i class="bi bi-person-badge"></i> vCard
              </button>
              <button
                class="btn btn-sm {qrMode === 'event' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('event')}
              >
                <i class="bi bi-calendar-event"></i> Termin
              </button>
              <button
                class="btn btn-sm {qrMode === 'crypto' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('crypto')}
              >
                <i class="bi bi-currency-bitcoin"></i> Krypto
              </button>
              <button
                class="btn btn-sm {qrMode === 'geo' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('geo')}
              >
                <i class="bi bi-geo-alt"></i> Maps (Geo)
              </button>
              <button
                class="btn btn-sm {qrMode === 'tel' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('tel')}
              >
                <i class="bi bi-telephone"></i> Telefon
              </button>
              <button
                class="btn btn-sm {qrMode === 'sms' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('sms')}
              >
                <i class="bi bi-chat-dots"></i> SMS
              </button>
              <button
                class="btn btn-sm {qrMode === 'email' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('email')}
              >
                <i class="bi bi-envelope"></i> E-Mail
              </button>
            </div>
          </div>

          <!-- Symbology Selection -->
          <div class="mb-3">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label for="singleTypeSelect" class="form-label fw-medium small text-body-secondary mb-0">Symbologie</label>
              {#if qrMode !== 'text'}
                <span class="badge bg-info-subtle text-info-emphasis small">2D Matrix empfohlen</span>
              {/if}
            </div>
            <select
              id="singleTypeSelect"
              class="form-select form-select-sm"
              value={selectedType}
              on:change={(e) => handleTypeChange(e.currentTarget.value as BarcodeType)}
            >
              <optgroup label="2D Matrix (Mehrzeilig / Große Datenmengen)">
                {#each BARCODE_TYPES.filter((t) => t.category === '2D Matrix') as t}
                  <option value={t.id}>{t.name} ({t.description})</option>
                {/each}
              </optgroup>
              {#if qrMode === 'text'}
                <optgroup label="1D Linear (Einzelhandels- & Industrie-Codes)">
                  {#each BARCODE_TYPES.filter((t) => t.category === '1D Linear') as t}
                    <option value={t.id}>{t.name} - {t.description}</option>
                  {/each}
                </optgroup>
              {/if}
            </select>
          </div>

          <!-- Structured Forms -->
          {#if qrMode === 'epc'}
            <!-- EPC / GiroCode -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="d-flex align-items-center mb-2">
                <i class="bi bi-bank text-primary me-2"></i>
                <span class="fw-medium small">GiroCode / SEPA-Überweisung (EPC-QR)</span>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-7">
                  <label for="epcNameInput" class="form-label small mb-1 fw-medium">Empfängername</label>
                  <input
                    id="epcNameInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="z.B. Mustermann GmbH"
                    bind:value={epcName}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-5">
                  <label for="epcAmountInput" class="form-label small mb-1 fw-medium">Betrag (€)</label>
                  <input
                    id="epcAmountInput"
                    type="number"
                    step="0.01"
                    class="form-control form-control-sm"
                    placeholder="0.00"
                    bind:value={epcAmount}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-8">
                  <label for="epcIbanInput" class="form-label small mb-1 fw-medium">IBAN</label>
                  <input
                    id="epcIbanInput"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="DE89 3704 0044 0532 0130 00"
                    bind:value={epcIBAN}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-4">
                  <label for="epcBicInput" class="form-label small mb-1 fw-medium">BIC (optional)</label>
                  <input
                    id="epcBicInput"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="GENODEFFXXX"
                    bind:value={epcBIC}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
              <div>
                <label for="epcRefInput" class="form-label small mb-1 fw-medium">Verwendungszweck</label>
                <input
                  id="epcRefInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="z.B. Rechnungsnummer 1002"
                  bind:value={epcRef}
                  on:input={updateStructuredQR}
                />
              </div>
            </div>
          {:else if qrMode === 'crypto'}
            <!-- Crypto -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="d-flex align-items-center mb-2">
                <i class="bi bi-currency-bitcoin text-warning me-2"></i>
                <span class="fw-medium small">Krypto Wallet & Zahlungsadresse</span>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-4">
                  <label for="cryptoCoinSelect" class="form-label small mb-1 fw-medium">Kryptowährung</label>
                  <select
                    id="cryptoCoinSelect"
                    class="form-select form-select-sm"
                    bind:value={cryptoCoin}
                    on:change={updateStructuredQR}
                  >
                    <option value="bitcoin">Bitcoin (BTC)</option>
                    <option value="ethereum">Ethereum (ETH)</option>
                    <option value="solana">Solana (SOL)</option>
                    <option value="dogecoin">Dogecoin (DOGE)</option>
                    <option value="litecoin">Litecoin (LTC)</option>
                  </select>
                </div>
                <div class="col-8">
                  <label for="cryptoAddressInput" class="form-label small mb-1 fw-medium">Wallet-Adresse</label>
                  <input
                    id="cryptoAddressInput"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="Adresse einfügen"
                    bind:value={cryptoAddress}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
              <div class="row g-2">
                <div class="col-5">
                  <label for="cryptoAmountInput" class="form-label small mb-1 fw-medium">Betrag (optional)</label>
                  <input
                    id="cryptoAmountInput"
                    type="number"
                    step="0.0001"
                    class="form-control form-control-sm"
                    placeholder="0.00"
                    bind:value={cryptoAmount}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-7">
                  <label for="cryptoMessageInput" class="form-label small mb-1 fw-medium">Nachricht / Verwendungszweck</label>
                  <input
                    id="cryptoMessageInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="z.B. Spende"
                    bind:value={cryptoMessage}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
            </div>
          {:else if qrMode === 'geo'}
            <!-- Maps Geo -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="d-flex align-items-center mb-2">
                <i class="bi bi-geo-alt text-danger me-2"></i>
                <span class="fw-medium small">Geo-Koordinaten (Google Maps & Apple Maps)</span>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-6">
                  <label for="geoLatInput" class="form-label small mb-1 fw-medium">Breitengrad (Latitude)</label>
                  <input
                    id="geoLatInput"
                    type="number"
                    step="0.000001"
                    class="form-control form-control-sm"
                    placeholder="52.5200"
                    bind:value={geoLat}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-6">
                  <label for="geoLonInput" class="form-label small mb-1 fw-medium">Längengrad (Longitude)</label>
                  <input
                    id="geoLonInput"
                    type="number"
                    step="0.000001"
                    class="form-control form-control-sm"
                    placeholder="13.4050"
                    bind:value={geoLon}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
              <div>
                <label for="geoQueryInput" class="form-label small mb-1 fw-medium">Ortsname / Suchbegriff (optional)</label>
                <input
                  id="geoQueryInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="z.B. Berliner Fernsehturm"
                  bind:value={geoQuery}
                  on:input={updateStructuredQR}
                />
              </div>
            </div>
          {:else if qrMode === 'tel'}
            <!-- Telephone -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="d-flex align-items-center mb-2">
                <i class="bi bi-telephone text-success me-2"></i>
                <span class="fw-medium small">Telefonanruf (tel:)</span>
              </div>
              <div>
                <label for="telInput" class="form-label small mb-1 fw-medium">Telefonnummer</label>
                <input
                  id="telInput"
                  type="tel"
                  class="form-control form-control-sm"
                  placeholder="+49 170 1234567"
                  bind:value={telNumber}
                  on:input={updateStructuredQR}
                />
              </div>
            </div>
          {:else if qrMode === 'sms'}
            <!-- SMS -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="d-flex align-items-center mb-2">
                <i class="bi bi-chat-dots text-info me-2"></i>
                <span class="fw-medium small">SMS-Nachricht (smsto:)</span>
              </div>
              <div class="mb-2">
                <label for="smsNumberInput" class="form-label small mb-1 fw-medium">Empfängernummer</label>
                <input
                  id="smsNumberInput"
                  type="tel"
                  class="form-control form-control-sm"
                  placeholder="+49 170 1234567"
                  bind:value={smsNumber}
                  on:input={updateStructuredQR}
                />
              </div>
              <div>
                <label for="smsMsgInput" class="form-label small mb-1 fw-medium">SMS-Text</label>
                <input
                  id="smsMsgInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="Vorgefertigte Nachricht..."
                  bind:value={smsMessage}
                  on:input={updateStructuredQR}
                />
              </div>
            </div>
          {:else if qrMode === 'email'}
            <!-- Email -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="d-flex align-items-center mb-2">
                <i class="bi bi-envelope text-primary me-2"></i>
                <span class="fw-medium small">E-Mail verfassen (mailto:)</span>
              </div>
              <div class="mb-2">
                <label for="mailToInput" class="form-label small mb-1 fw-medium">Empfängeradresse</label>
                <input
                  id="mailToInput"
                  type="email"
                  class="form-control form-control-sm"
                  placeholder="info@beispiel.de"
                  bind:value={mailTo}
                  on:input={updateStructuredQR}
                />
              </div>
              <div class="mb-2">
                <label for="mailSubInput" class="form-label small mb-1 fw-medium">Betreffzeile</label>
                <input
                  id="mailSubInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="Betreff eingeben"
                  bind:value={mailSubject}
                  on:input={updateStructuredQR}
                />
              </div>
              <div>
                <label for="mailBodyInput" class="form-label small mb-1 fw-medium">Nachrichtentext</label>
                <textarea
                  id="mailBodyInput"
                  rows="2"
                  class="form-control form-control-sm"
                  placeholder="Mailtext..."
                  bind:value={mailBody}
                  on:input={updateStructuredQR}
                ></textarea>
              </div>
            </div>
          {:else if qrMode === 'wifi'}
            <!-- WLAN -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="row g-2 mb-2">
                <div class="col-8">
                  <label for="wifiSsidInput" class="form-label small mb-1 fw-medium">WLAN-Name (SSID)</label>
                  <input
                    id="wifiSsidInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder="z.B. MeinWLAN"
                    bind:value={wifiSSID}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-4">
                  <label for="wifiEncSelect" class="form-label small mb-1 fw-medium">Verschlüsselung</label>
                  <select
                    id="wifiEncSelect"
                    class="form-select form-select-sm"
                    bind:value={wifiEnc}
                    on:change={updateStructuredQR}
                  >
                    <option value="WPA">WPA / WPA2 / WPA3</option>
                    <option value="WEP">WEP</option>
                    <option value="nopass">Offen (Kein Passwort)</option>
                  </select>
                </div>
              </div>
              {#if wifiEnc !== 'nopass'}
                <div class="mb-2">
                  <label for="wifiPassInput" class="form-label small mb-1 fw-medium">WLAN-Passwort</label>
                  <input
                    id="wifiPassInput"
                    type="password"
                    class="form-control form-control-sm font-monospace"
                    placeholder="Passwort eingeben"
                    bind:value={wifiPass}
                    on:input={updateStructuredQR}
                  />
                </div>
              {/if}
            </div>
          {:else if qrMode === 'vcard'}
            <!-- vCard -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="row g-2 mb-2">
                <div class="col-6">
                  <label for="vcardFirstInput" class="form-label small mb-1 fw-medium">Vorname</label>
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
                  <label for="vcardLastInput" class="form-label small mb-1 fw-medium">Nachname</label>
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
                  <label for="vcardEmailInput" class="form-label small mb-1 fw-medium">E-Mail</label>
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
                  <label for="vcardPhoneInput" class="form-label small mb-1 fw-medium">Telefon</label>
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
            <!-- Event -->
            <div class="p-3 bg-body-secondary rounded mb-3 border">
              <div class="mb-2">
                <label for="eventSummaryInput" class="form-label small mb-1 fw-medium">Titel / Anlass</label>
                <input
                  id="eventSummaryInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder="z.B. Jahres-Hauptversammlung"
                  bind:value={eventSummary}
                  on:input={updateStructuredQR}
                />
              </div>
              <div class="row g-2">
                <div class="col-6">
                  <label for="eventStartInput" class="form-label small mb-1 fw-medium">Beginn (YYYYMMDDTHHMMSS)</label>
                  <input
                    id="eventStartInput"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="20260901T100000"
                    bind:value={eventStart}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-6">
                  <label for="eventEndInput" class="form-label small mb-1 fw-medium">Ende (YYYYMMDDTHHMMSS)</label>
                  <input
                    id="eventEndInput"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder="20260901T120000"
                    bind:value={eventEnd}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
            </div>
          {/if}

          <!-- Raw Data Input -->
          <div class="mb-3">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label for="singleDataInput" class="form-label fw-medium small text-body-secondary mb-0">
                {qrMode === 'text' ? 'Inhalt / Nutzdaten' : 'Generierte Nutzdaten (Raw)'}
              </label>
              <span class="badge bg-secondary-subtle text-secondary-emphasis small font-monospace">
                {rawData.length} Zeichen
              </span>
            </div>
            <textarea
              id="singleDataInput"
              class="form-control form-control-sm font-monospace"
              rows={qrMode === 'text' ? 3 : 2}
              placeholder="Text oder Link eingeben..."
              bind:value={rawData}
              on:input={triggerGenerate}
              readonly={qrMode !== 'text'}
            ></textarea>
          </div>

          <!-- Custom Text / Caption Input -->
          <div class="p-3 bg-body-secondary rounded mb-3 border">
            <div class="form-check form-switch mb-2">
              <input
                id="showTextCheck"
                type="checkbox"
                class="form-check-input"
                bind:checked={showText}
                on:change={triggerGenerate}
              />
              <label for="showTextCheck" class="form-check-label small fw-medium">
                Klartext-Beschriftung unter dem Barcode anzeigen
              </label>
            </div>
            {#if showText}
              <div>
                <label for="customLabelInput" class="form-label small mb-1 text-body-secondary">
                  Beschriftungstext (Freitext für Etikett oder Scan-Hinweis)
                </label>
                <input
                  id="customLabelInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={rawData.slice(0, 40) || 'z.B. Artikel-Nr. 12345'}
                  bind:value={customLabelText}
                  on:input={() => {
                    customLabelTouched = true;
                    triggerGenerate();
                  }}
                />
              </div>
            {/if}
          </div>

          <!-- Styling & Colors -->
          <div class="row g-2 mb-3">
            <!-- Background -->
            <div class="col-6">
              <span class="form-label fw-medium small text-body-secondary mb-1 d-block">Hintergrund</span>
              <div class="d-flex align-items-center gap-2 mb-2">
                <input
                  id="singleBgColor"
                  type="color"
                  class="form-control form-control-color form-control-sm"
                  bind:value={bgColor}
                  on:input={() => {
                    isTransparent = false;
                    triggerGenerate();
                  }}
                  disabled={isTransparent}
                />
                <button
                  type="button"
                  class="btn btn-sm {isTransparent ? 'btn-primary' : 'btn-outline-secondary'} flex-grow-1"
                  on:click={toggleTransparent}
                >
                  <i class="bi bi-circle-half me-1"></i> {isTransparent ? 'Transparent ✓' : 'Transparent'}
                </button>
              </div>
              <div class="d-flex flex-wrap gap-1">
                {#each PRESET_BG_COLORS as preset}
                  <button
                    type="button"
                    class="btn btn-xs btn-outline-secondary p-1"
                    style="width: 22px; height: 22px; background-color: {preset.value}; border-radius: 4px;"
                    title="Hintergrund: {preset.label}"
                    on:click={() => setBgColor(preset.value)}
                  ></button>
                {/each}
              </div>
            </div>

            <!-- Foreground -->
            <div class="col-6">
              <label for="singleFgColor" class="form-label fw-medium small text-body-secondary mb-1">Barcode-Farbe</label>
              <div class="d-flex align-items-center gap-2 mb-2">
                <input
                  id="singleFgColor"
                  type="color"
                  class="form-control form-control-color form-control-sm"
                  bind:value={fgColor}
                  on:input={triggerGenerate}
                />
                <span class="small font-monospace text-body-secondary">{fgColor}</span>
              </div>
              <div class="d-flex flex-wrap gap-1">
                {#each PRESET_FG_COLORS as preset}
                  <button
                    type="button"
                    class="btn btn-xs btn-outline-secondary p-1"
                    style="width: 22px; height: 22px; background-color: {preset.value}; border-radius: 4px;"
                    title="Barcode-Farbe: {preset.label}"
                    on:click={() => setFgColor(preset.value)}
                  ></button>
                {/each}
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
            <span class="badge bg-primary-subtle text-primary small">
              <span class="spinner-border spinner-border-sm me-1"></span> Render...
            </span>
          {:else if result?.success}
            <span class="badge bg-success-subtle text-success small">Bereit</span>
          {/if}
        </div>

        <div class="card-body d-flex flex-column align-items-center justify-content-center p-4">
          {#if result?.success && result?.svg}
            <div
              class="svg-container rounded p-3 mb-3 d-flex align-items-center justify-content-center shadow-sm"
              style="background-color: {isTransparent ? 'repeating-conic-gradient(#808080 0% 25%, transparent 0% 50%) 50% / 16px 16px' : bgColor};"
            >
              {@html result.svg}
            </div>

            <!-- Format Badge -->
            <div class="d-flex gap-2 mb-3">
              <span class="badge bg-body-secondary text-body border small">
                Typ: <strong class="text-uppercase">{result.type}</strong>
              </span>
              <span class="badge bg-body-secondary text-body border small">
                Vektor: <strong>SVG / Crisp</strong>
              </span>
            </div>
          {:else if result?.error}
            <div class="alert alert-danger w-100 text-center py-4 my-auto">
              <i class="bi bi-exclamation-octagon fs-2 d-block mb-2 text-danger"></i>
              <strong class="d-block mb-1">Ungültige Eingabedaten für {selectedType.toUpperCase()}</strong>
              <small class="text-body-secondary">{result.error}</small>
            </div>
          {:else}
            <div class="text-center text-body-secondary py-5 my-auto">
              <i class="bi bi-upc-scan fs-1 d-block mb-2 opacity-50"></i>
              <span>Geben Sie Daten ein, um den Barcode in Echtzeit zu generieren.</span>
            </div>
          {/if}
        </div>

        <!-- Export Buttons -->
        {#if result?.success}
          <div class="card-footer bg-body border-top p-3">
            <div class="row g-2">
              <div class="col-6">
                <button type="button" class="btn btn-outline-primary btn-sm w-100" on:click={copySVG}>
                  <i class="bi bi-clipboard me-1"></i> SVG kopieren
                </button>
              </div>
              <div class="col-6">
                <button type="button" class="btn btn-outline-secondary btn-sm w-100" on:click={copyPNG}>
                  <i class="bi bi-image me-1"></i> PNG kopieren
                </button>
              </div>
              <div class="col-6">
                <button type="button" class="btn btn-primary btn-sm w-100" on:click={saveSVG}>
                  <i class="bi bi-download me-1"></i> SVG speichern...
                </button>
              </div>
              <div class="col-6">
                <button type="button" class="btn btn-primary btn-sm w-100" on:click={savePNG}>
                  <i class="bi bi-download me-1"></i> PNG speichern...
                </button>
              </div>
              {#if onSendToPrint}
                <div class="col-12 mt-2">
                  <button type="button" class="btn btn-success btn-sm w-100" on:click={sendToLabelPrint}>
                    <i class="bi bi-printer me-1"></i> Zum Etikettendruck (DIN A4) hinzufügen
                  </button>
                </div>
              {/if}
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .svg-container {
    max-width: 380px;
    width: 100%;
    height: 260px;
    border: 1px solid rgba(128, 128, 128, 0.2);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .svg-container :global(svg) {
    width: 100% !important;
    height: 100% !important;
    max-width: 100%;
    max-height: 100%;
    display: block;
    object-fit: contain;
  }
</style>
