<script lang="ts">
  import { onMount } from 'svelte';
  import { BARCODE_TYPES, type BarcodeType } from '../types';
  import { formatError } from '../i18n/errors';
  import { lang } from '../i18n/lang';
  import { t, type SingleKey } from '../i18n/text/single';
  import { typeText, typeDescKey } from '../i18n/text/types';
  import { toLocalInput, nextFullHour, shiftLocal, icalDateTime, icalDate, withDate, deviceTimeZone } from '../datetime';
  import { profile, organizerName, profileCoords } from '../profile';
  import { get } from 'svelte/store';
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
    ValidateBarcode,
    CopyToClipboard
  } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import type { RetailValidation } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
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
  let epcRef = $t('sampleEpcRef');

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
  const HOUR_MS = 3_600_000;
  const DAY_MS = 24 * HOUR_MS;

  // Held as <input type="datetime-local"> values (YYYY-MM-DDTHH:MM, local time).
  let eventAllDay = false;
  let eventStart = toLocalInput(nextFullHour());
  let eventEnd = shiftLocal(eventStart, HOUR_MS);
  let eventTZ = get(profile).timeZone || deviceTimeZone() || 'Europe/Berlin';
  let eventLocation = '';
  let eventLat = 0;
  let eventLon = 0;
  let eventOrganizer = '';
  let eventOrganizerEmail = '';

  // Structured QR inputs: Crypto
  let cryptoCoin = 'bitcoin';
  let cryptoAddress = '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa';
  let cryptoAmount: number | string = 0.005;
  let cryptoMessage = $t('sampleCryptoMessage');

  // Structured QR inputs: Geo
  let geoLat = 52.5200;
  let geoLon = 13.4050;
  let geoQuery = $t('sampleGeoQuery');

  // Structured QR inputs: Tel & SMS
  let telNumber = '+49 170 1234567';
  let smsNumber = '+49 170 1234567';
  let smsMessage = $t('sampleSms');

  // Structured QR inputs: Email
  let mailTo = 'support@mlcgo.eu';
  let mailSubject = $t('sampleMailSubject');
  let mailBody = $t('sampleMailBody');

  // Styling options
  let fgColor = '#000000';
  let bgColor = '#ffffff';
  let isTransparent = false;
  let quietZone = true; // blank margin scanners need
  let showText = false;
  let fontSize = 0; // caption px, 0 = automatic
  let customWidth = 0;
  let customHeight = 0;

  // Preset Color Palettes
  const PRESET_BG_COLORS = [
    { label: 'colorWhite', value: '#ffffff' },
    { label: 'colorYellow', value: '#fff59d' },
    { label: 'colorLightBlue', value: '#e1f5fe' },
    { label: 'colorLightGreen', value: '#e8f5e9' },
    { label: 'colorLightGray', value: '#f8f9fa' }
  ] satisfies { label: SingleKey; value: string }[];

  const PRESET_FG_COLORS = [
    { label: 'colorBlack', value: '#000000' },
    { label: 'colorDarkBlue', value: '#0d47a1' },
    { label: 'colorDarkRed', value: '#b71c1c' },
    { label: 'colorDarkGreen', value: '#1b5e20' }
  ] satisfies { label: SingleKey; value: string }[];

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
        customLabelText =
          numAmount > 0
            ? $t('captionTransfer', { amount: numAmount.toFixed(2), name: epcName })
            : $t('captionGiro', { name: epcName });
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
        customLabelText = $t('captionEmail', { to: mailTo });
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
        customLabelText = wifiSSID ? $t('captionWifi', { ssid: wifiSSID }) : '';
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
        startTime: eventAllDay ? icalDate(eventStart) : icalDateTime(eventStart),
        // iCal all-day DTEND is exclusive: the day after the last day.
        endTime: eventAllDay ? icalDate(shiftLocal(eventEnd, DAY_MS)) : icalDateTime(eventEnd),
        timeZone: eventTZ,
        location: eventLocation,
        latitude: eventLat,
        longitude: eventLon,
        organizer: eventOrganizer,
        organizerEmail: eventOrganizerEmail
      });
      if (!customLabelTouched) {
        customLabelText = eventSummary;
      }
    }
    triggerGenerate();
  }


  // Moving the start moves the end along and keeps the duration.
  function setEventStart(value: string) {
    if (!value) return;
    const duration = new Date(eventEnd).getTime() - new Date(eventStart).getTime();
    eventStart = value;
    eventEnd = shiftLocal(value, duration > 0 ? duration : eventAllDay ? 0 : HOUR_MS);
    updateStructuredQR();
  }

  function setEventEnd(value: string) {
    if (!value) return;
    eventEnd = value < eventStart ? eventStart : value;
    updateStructuredQR();
  }


  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  function triggerGenerate() {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      generate();
    }, 100);
  }

  // Live check of EAN/UPC input (length, check digit); null for other types.
  let retailCheck: RetailValidation | null = null;

  async function generate() {
    if (!rawData.trim()) {
      result = null;
      retailCheck = null;
      return;
    }

    try {
      const check = await ValidateBarcode(selectedType, rawData);
      retailCheck = check.applies ? check : null;
    } catch {
      retailCheck = null;
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
        fontSize,
        foregroundColor: fgColor,
        backgroundColor: isTransparent ? 'transparent' : bgColor,
        noQuietZone: !quietZone
      });
      result = res;
    } catch (e: any) {
      result = {
        type: selectedType,
        data: rawData,
        success: false,
        error: e?.message || $t('generateError')
      };
    } finally {
      isGenerating = false;
    }
  }

  // Empty fields of the chosen mode are filled from "My details"; what the
  // user typed stays.
  let prefilled = false;
  function prefillFromProfile(mode: QRMode) {
    const p = get(profile);
    const coords = profileCoords(p);
    prefilled = false;
    if (mode === 'event') {
      if (!eventOrganizer && organizerName(p)) (eventOrganizer = organizerName(p)), (prefilled = true);
      if (!eventOrganizerEmail && p.email) (eventOrganizerEmail = p.email), (prefilled = true);
      if (!eventLocation && p.location) (eventLocation = p.location), (prefilled = true);
      if (coords && !eventLat && !eventLon) (eventLat = coords.latitude), (eventLon = coords.longitude);
      if (p.timeZone) eventTZ = p.timeZone;
    } else if (mode === 'geo' && coords) {
      geoLat = coords.latitude;
      geoLon = coords.longitude;
      if (p.location) geoQuery = p.location;
      prefilled = true;
    } else if (mode === 'vcard') {
      if (!vcardFirst && !vcardLast && (p.firstName || p.lastName)) {
        vcardFirst = p.firstName;
        vcardLast = p.lastName;
        prefilled = true;
      }
      if (!vcardEmail && p.email) vcardEmail = p.email;
      if (!vcardPhone && p.phone) vcardPhone = p.phone;
    }
  }

  function handleModeChange(mode: QRMode) {
    qrMode = mode;
    customLabelTouched = false;
    prefillFromProfile(mode);
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
    const opt = BARCODE_TYPES.find((b) => b.id === newType);
    if (opt && rawData === 'https://mlcgo.eu' && newType !== 'qr') {
      rawData = opt.sample;
    }
    triggerGenerate();
  }

  async function copySVG() {
    if (!result?.svg) return;
    const ok = await CopyToClipboard(result.svg);
    if (ok) {
      showFeedback($t('svgCopiedNative'));
    } else {
      await navigator.clipboard.writeText(result.svg);
      showFeedback($t('svgCopied'));
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
      showFeedback($t('pngCopied'));
    } catch (err) {
      showFeedback($t('pngCopyFailed'), 'danger');
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
        showFeedback($t('savedTo', { path: savedPath }));
      }
    } catch (e: any) {
      showFeedback($t('saveError', { msg: e?.message ?? '' }), 'danger');
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
        showFeedback($t('savedTo', { path: savedPath }));
      }
    } catch (e: any) {
      showFeedback($t('saveError', { msg: e?.message ?? '' }), 'danger');
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

  // Type descriptions in the current language.
  $: typeDesc = (id: string) => $typeText(typeDescKey(id));

  // A language switch rewords the default caption of a structured code.
  let captionLang = $lang;
  $: if ($lang !== captionLang) {
    captionLang = $lang;
    if (qrMode !== 'text' && !customLabelTouched) updateStructuredQR();
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
            <i class="bi bi-sliders me-1 text-primary"></i> {$t('contentHeading')}
          </h6>
        </div>
        <div class="card-body">
          <!-- Inhaltsformat / Typ-Modus Switcher -->
          <div class="mb-3">
            <span class="form-label fw-medium small text-body-secondary d-block mb-1">{$t('contentType')}</span>
            <div class="d-flex flex-wrap gap-1 bg-body-secondary p-1 rounded">
              <button
                class="btn btn-sm {qrMode === 'text' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('text')}
              >
                <i class="bi bi-fonts"></i> {$t('modeText')}
              </button>
              <button
                class="btn btn-sm {qrMode === 'epc' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('epc')}
              >
                <i class="bi bi-bank"></i> {$t('modeEpc')}
              </button>
              <button
                class="btn btn-sm {qrMode === 'wifi' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('wifi')}
              >
                <i class="bi bi-wifi"></i> {$t('modeWifi')}
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
                <i class="bi bi-calendar-event"></i> {$t('modeEvent')}
              </button>
              <button
                class="btn btn-sm {qrMode === 'crypto' ? 'btn-primary' : 'btn-light border-0'} py-1 px-2"
                type="button"
                on:click={() => handleModeChange('crypto')}
              >
                <i class="bi bi-currency-bitcoin"></i> {$t('modeCrypto')}
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
                <i class="bi bi-telephone"></i> {$t('modeTel')}
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
                <i class="bi bi-envelope"></i> {$t('email')}
              </button>
            </div>
          </div>

          <!-- Symbology Selection -->
          <div class="mb-3">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label for="singleTypeSelect" class="form-label fw-medium small text-body-secondary mb-0">{$t('symbology')}</label>
              {#if qrMode !== 'text'}
                <span class="badge bg-info-subtle text-info-emphasis small">{$t('matrixRecommended')}</span>
              {/if}
            </div>
            <select
              id="singleTypeSelect"
              class="form-select form-select-sm"
              value={selectedType}
              on:change={(e) => handleTypeChange(e.currentTarget.value as BarcodeType)}
            >
              <optgroup label={$t('group2dMatrix')}>
                <!-- Structured payloads (vCard, Wi-Fi, …) are read reliably by phones only as QR or DataMatrix. -->
                {#each BARCODE_TYPES.filter((b) => b.category === '2D Matrix' && (qrMode === 'text' || b.id === 'qr' || b.id === 'datamatrix')) as bt}
                  <option value={bt.id}>{bt.name} ({typeDesc(bt.id)})</option>
                {/each}
              </optgroup>
              {#if qrMode === 'text'}
                <optgroup label={$t('group2dStacked')}>
                  {#each BARCODE_TYPES.filter((b) => b.category === '2D Stacked') as bt}
                    <option value={bt.id}>{bt.name} - {typeDesc(bt.id)}</option>
                  {/each}
                </optgroup>
                <optgroup label={$t('group1dLinear')}>
                  {#each BARCODE_TYPES.filter((b) => b.category === '1D Linear') as bt}
                    <option value={bt.id}>{bt.name} - {typeDesc(bt.id)}</option>
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
                <span class="fw-medium small">{$t('epcTitle')}</span>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-7">
                  <label for="epcNameInput" class="form-label small mb-1 fw-medium">{$t('epcName')}</label>
                  <input
                    id="epcNameInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder={$t('epcNamePh')}
                    bind:value={epcName}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-5">
                  <label for="epcAmountInput" class="form-label small mb-1 fw-medium">{$t('epcAmount')}</label>
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
                <label for="epcRefInput" class="form-label small mb-1 fw-medium">{$t('epcRef')}</label>
                <input
                  id="epcRefInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={$t('epcRefPh')}
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
                <span class="fw-medium small">{$t('cryptoTitle')}</span>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-4">
                  <label for="cryptoCoinSelect" class="form-label small mb-1 fw-medium">{$t('cryptoCoin')}</label>
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
                  <label for="cryptoAddressInput" class="form-label small mb-1 fw-medium">{$t('cryptoAddress')}</label>
                  <input
                    id="cryptoAddressInput"
                    type="text"
                    class="form-control form-control-sm font-monospace"
                    placeholder={$t('cryptoAddressPh')}
                    bind:value={cryptoAddress}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
              <div class="row g-2">
                <div class="col-5">
                  <label for="cryptoAmountInput" class="form-label small mb-1 fw-medium">{$t('cryptoAmount')}</label>
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
                  <label for="cryptoMessageInput" class="form-label small mb-1 fw-medium">{$t('cryptoMessage')}</label>
                  <input
                    id="cryptoMessageInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder={$t('cryptoMessagePh')}
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
                <span class="fw-medium small">{$t('geoTitle')}</span>
              </div>
              <div class="row g-2 mb-2">
                <div class="col-6">
                  <label for="geoLatInput" class="form-label small mb-1 fw-medium">{$t('geoLat')}</label>
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
                  <label for="geoLonInput" class="form-label small mb-1 fw-medium">{$t('geoLon')}</label>
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
                <label for="geoQueryInput" class="form-label small mb-1 fw-medium">{$t('geoQuery')}</label>
                <input
                  id="geoQueryInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={$t('geoQueryPh')}
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
                <span class="fw-medium small">{$t('telTitle')}</span>
              </div>
              <div>
                <label for="telInput" class="form-label small mb-1 fw-medium">{$t('telNumber')}</label>
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
                <span class="fw-medium small">{$t('smsTitle')}</span>
              </div>
              <div class="mb-2">
                <label for="smsNumberInput" class="form-label small mb-1 fw-medium">{$t('smsNumber')}</label>
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
                <label for="smsMsgInput" class="form-label small mb-1 fw-medium">{$t('smsText')}</label>
                <input
                  id="smsMsgInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={$t('smsTextPh')}
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
                <span class="fw-medium small">{$t('mailTitle')}</span>
              </div>
              <div class="mb-2">
                <label for="mailToInput" class="form-label small mb-1 fw-medium">{$t('mailTo')}</label>
                <input
                  id="mailToInput"
                  type="email"
                  class="form-control form-control-sm"
                  placeholder={$t('mailToPh')}
                  bind:value={mailTo}
                  on:input={updateStructuredQR}
                />
              </div>
              <div class="mb-2">
                <label for="mailSubInput" class="form-label small mb-1 fw-medium">{$t('mailSubject')}</label>
                <input
                  id="mailSubInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={$t('mailSubjectPh')}
                  bind:value={mailSubject}
                  on:input={updateStructuredQR}
                />
              </div>
              <div>
                <label for="mailBodyInput" class="form-label small mb-1 fw-medium">{$t('mailBody')}</label>
                <textarea
                  id="mailBodyInput"
                  rows="2"
                  class="form-control form-control-sm"
                  placeholder={$t('mailBodyPh')}
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
                  <label for="wifiSsidInput" class="form-label small mb-1 fw-medium">{$t('wifiSsid')}</label>
                  <input
                    id="wifiSsidInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder={$t('wifiSsidPh')}
                    bind:value={wifiSSID}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-4">
                  <label for="wifiEncSelect" class="form-label small mb-1 fw-medium">{$t('wifiEnc')}</label>
                  <select
                    id="wifiEncSelect"
                    class="form-select form-select-sm"
                    bind:value={wifiEnc}
                    on:change={updateStructuredQR}
                  >
                    <option value="WPA">WPA / WPA2 / WPA3</option>
                    <option value="WEP">WEP</option>
                    <option value="nopass">{$t('wifiOpen')}</option>
                  </select>
                </div>
              </div>
              {#if wifiEnc !== 'nopass'}
                <div class="mb-2">
                  <label for="wifiPassInput" class="form-label small mb-1 fw-medium">{$t('wifiPass')}</label>
                  <input
                    id="wifiPassInput"
                    type="password"
                    class="form-control form-control-sm font-monospace"
                    placeholder={$t('wifiPassPh')}
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
                  <label for="vcardFirstInput" class="form-label small mb-1 fw-medium">{$t('vcardFirst')}</label>
                  <input
                    id="vcardFirstInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder={$t('vcardFirstPh')}
                    bind:value={vcardFirst}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-6">
                  <label for="vcardLastInput" class="form-label small mb-1 fw-medium">{$t('vcardLast')}</label>
                  <input
                    id="vcardLastInput"
                    type="text"
                    class="form-control form-control-sm"
                    placeholder={$t('vcardLastPh')}
                    bind:value={vcardLast}
                    on:input={updateStructuredQR}
                  />
                </div>
              </div>
              <div class="row g-2">
                <div class="col-6">
                  <label for="vcardEmailInput" class="form-label small mb-1 fw-medium">{$t('email')}</label>
                  <input
                    id="vcardEmailInput"
                    type="email"
                    class="form-control form-control-sm"
                    placeholder={$t('vcardEmailPh')}
                    bind:value={vcardEmail}
                    on:input={updateStructuredQR}
                  />
                </div>
                <div class="col-6">
                  <label for="vcardPhoneInput" class="form-label small mb-1 fw-medium">{$t('phone')}</label>
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
                <label for="eventSummaryInput" class="form-label small mb-1 fw-medium">{$t('eventTitle')}</label>
                <input
                  id="eventSummaryInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={$t('eventTitlePh')}
                  bind:value={eventSummary}
                  on:input={updateStructuredQR}
                />
              </div>
              <div class="form-check form-switch mb-2">
                <input
                  id="eventAllDay"
                  type="checkbox"
                  class="form-check-input"
                  bind:checked={eventAllDay}
                  on:change={updateStructuredQR}
                />
                <label for="eventAllDay" class="form-check-label small">{$t('allDay')}</label>
              </div>
              <div class="row g-2">
                <div class="col-6">
                  <label for="eventStartInput" class="form-label small mb-1 fw-medium">{$t('start')}</label>
                  {#if eventAllDay}
                    <input
                      id="eventStartInput"
                      type="date"
                      class="form-control form-control-sm"
                      value={eventStart.slice(0, 10)}
                      on:change={(e) => setEventStart(withDate(eventStart, e.currentTarget.value))}
                    />
                  {:else}
                    <input
                      id="eventStartInput"
                      type="datetime-local"
                      class="form-control form-control-sm"
                      value={eventStart}
                      on:change={(e) => setEventStart(e.currentTarget.value)}
                    />
                  {/if}
                </div>
                <div class="col-6">
                  <label for="eventEndInput" class="form-label small mb-1 fw-medium">{$t('end')}</label>
                  {#if eventAllDay}
                    <input
                      id="eventEndInput"
                      type="date"
                      class="form-control form-control-sm"
                      min={eventStart.slice(0, 10)}
                      value={eventEnd.slice(0, 10)}
                      on:change={(e) => setEventEnd(withDate(eventEnd, e.currentTarget.value))}
                    />
                  {:else}
                    <input
                      id="eventEndInput"
                      type="datetime-local"
                      class="form-control form-control-sm"
                      min={eventStart}
                      value={eventEnd}
                      on:change={(e) => setEventEnd(e.currentTarget.value)}
                    />
                  {/if}
                </div>
              </div>
              <div class="row g-2 mt-1">
                <div class="col-12">
                  <label for="eventLocationInput" class="form-label small mb-1 fw-medium">{$t('eventLocation')}</label>
                  <input id="eventLocationInput" type="text" class="form-control form-control-sm" placeholder={$t('eventLocationPh')} bind:value={eventLocation} on:input={updateStructuredQR} />
                </div>
                <div class="col-6">
                  <label for="eventOrganizerInput" class="form-label small mb-1 fw-medium">{$t('eventOrganizer')}</label>
                  <input id="eventOrganizerInput" type="text" class="form-control form-control-sm" placeholder={$t('eventOrganizerPh')} bind:value={eventOrganizer} on:input={updateStructuredQR} />
                </div>
                <div class="col-6">
                  <label for="eventOrganizerEmailInput" class="form-label small mb-1 fw-medium">{$t('eventOrganizerEmail')}</label>
                  <input id="eventOrganizerEmailInput" type="email" class="form-control form-control-sm" bind:value={eventOrganizerEmail} on:input={updateStructuredQR} />
                </div>
                {#if !eventAllDay}
                  <div class="col-12">
                    <label for="eventTZInput" class="form-label small mb-1 fw-medium">{$t('eventTimeZone')}</label>
                    <input id="eventTZInput" type="text" class="form-control form-control-sm" bind:value={eventTZ} on:input={updateStructuredQR} />
                  </div>
                {/if}
              </div>
              {#if prefilled}
                <div class="form-text small mt-2"><i class="bi bi-person-check me-1"></i>{$t('fromProfile')}</div>
              {/if}
            </div>
          {/if}

          <!-- Raw Data Input -->
          <div class="mb-3">
            <div class="d-flex justify-content-between align-items-center mb-1">
              <label for="singleDataInput" class="form-label fw-medium small text-body-secondary mb-0">
                {qrMode === 'text' ? $t('dataLabel') : $t('dataLabelRaw')}
              </label>
              <span class="badge bg-secondary-subtle text-secondary-emphasis small font-monospace">
                {$t('chars', { n: rawData.length })}
              </span>
            </div>
            <textarea
              id="singleDataInput"
              class="form-control form-control-sm font-monospace"
              rows={qrMode === 'text' ? 3 : 2}
              placeholder={$t('dataPh')}
              class:is-invalid={retailCheck && !retailCheck.valid}
              class:is-valid={retailCheck?.valid}
              bind:value={rawData}
              on:input={triggerGenerate}
              readonly={qrMode !== 'text'}
            ></textarea>
            {#if retailCheck}
              {@const typeName = BARCODE_TYPES.find((b) => b.id === selectedType)?.name ?? selectedType}
              {#if retailCheck.valid}
                <div class="valid-feedback d-block">
                  {#if retailCheck.checkDigitAdded}
                    {$t('checkDigitAddedPre')} <strong>{retailCheck.code.slice(-1)}</strong> {$t('checkDigitAddedPost')}
                    <span class="font-monospace">{retailCheck.code}</span>
                  {:else}
                    {$t('checkDigitOk')}
                  {/if}
                </div>
              {:else}
                <div class="invalid-feedback d-block">
                  {#if retailCheck.reason === 'non_digit'}
                    {$t('retailNonDigit', { type: typeName })}
                  {:else if retailCheck.reason === 'length'}
                    {$t('retailLength', {
                      type: typeName,
                      max: retailCheck.maxLength,
                      min: retailCheck.minLength,
                      n: rawData.trim().length
                    })}
                  {:else}
                    {$t('retailChecksum', { given: retailCheck.given })} <strong>{retailCheck.expected}</strong>.
                    <button
                      type="button"
                      class="btn btn-link btn-sm p-0 align-baseline"
                      on:click={() => {
                        rawData = retailCheck?.code ?? rawData;
                        triggerGenerate();
                      }}>{$t('fixIt')} <span class="font-monospace">{retailCheck.code}</span></button
                    >
                  {/if}
                </div>
              {/if}
            {/if}
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
                {$t('showCaption')}
              </label>
            </div>
            {#if showText}
              <div>
                <label for="customLabelInput" class="form-label small mb-1 text-body-secondary">
                  {$t('captionText')}
                </label>
                <input
                  id="customLabelInput"
                  type="text"
                  class="form-control form-control-sm"
                  placeholder={rawData.slice(0, 40) || $t('captionPh')}
                  bind:value={customLabelText}
                  on:input={() => {
                    customLabelTouched = true;
                    triggerGenerate();
                  }}
                />
                <div class="d-flex justify-content-between align-items-center mt-2 mb-1">
                  <label for="fontSizeRange" class="form-label small mb-0 text-body-secondary">{$t('fontSize')}</label>
                  <div class="d-flex align-items-center gap-2 small">
                    <span class="font-monospace">{fontSize > 0 ? `${fontSize} px` : 'Auto'}</span>
                    {#if fontSize > 0}
                      <button
                        type="button"
                        class="btn btn-link btn-sm p-0 small"
                        on:click={() => {
                          fontSize = 0;
                          triggerGenerate();
                        }}>Auto</button
                      >
                    {/if}
                  </div>
                </div>
                <input
                  id="fontSizeRange"
                  type="range"
                  class="form-range"
                  min="8"
                  max="96"
                  step="1"
                  value={fontSize || 24}
                  on:input={(e) => {
                    fontSize = Number(e.currentTarget.value);
                    triggerGenerate();
                  }}
                />
              </div>
            {/if}
          </div>

          <div class="form-check form-switch mb-3">
            <input
              id="quietZoneCheck"
              type="checkbox"
              class="form-check-input"
              bind:checked={quietZone}
              on:change={triggerGenerate}
            />
            <label for="quietZoneCheck" class="form-check-label small">
              {$t('quietZone')}
              <span class="text-body-secondary">{$t('quietZoneHint')}</span>
            </label>
          </div>

          <!-- Styling & Colors -->
          <div class="row g-2 mb-3">
            <!-- Background -->
            <div class="col-6">
              <span class="form-label fw-medium small text-body-secondary mb-1 d-block">{$t('background')}</span>
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
                    title={$t('backgroundTitle', { color: $t(preset.label) })}
                    on:click={() => setBgColor(preset.value)}
                  ></button>
                {/each}
              </div>
            </div>

            <!-- Foreground -->
            <div class="col-6">
              <label for="singleFgColor" class="form-label fw-medium small text-body-secondary mb-1">{$t('barcodeColor')}</label>
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
                    title={$t('barcodeColorTitle', { color: $t(preset.label) })}
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
            <i class="bi bi-eye me-1 text-primary"></i> {$t('livePreview')}
          </h6>
          {#if isGenerating}
            <span class="badge bg-primary-subtle text-primary small">
              <span class="spinner-border spinner-border-sm me-1"></span> {$t('rendering')}
            </span>
          {:else if result?.success}
            <span class="badge bg-success-subtle text-success small">{$t('ready')}</span>
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
                {$t('typeBadge')} <strong class="text-uppercase">{result.type}</strong>
              </span>
              <span class="badge bg-body-secondary text-body border small">
                {$t('vectorBadge')} <strong>SVG / Crisp</strong>
              </span>
              {#if result.readBack === 'ok'}
                <span class="badge bg-success-subtle text-success-emphasis border border-success-subtle small" title={$t('readOkTitle')}>
                  <i class="bi bi-check2-circle me-1"></i>{$t('readOk')}
                </span>
              {:else if result.readBack === 'unreadable'}
                <span class="badge bg-danger-subtle text-danger-emphasis border border-danger-subtle small" title={$t('readUnreadableTitle')}>
                  <i class="bi bi-exclamation-triangle me-1"></i>{$t('readUnreadable')}
                </span>
              {:else if result.readBack === 'mismatch'}
                <span class="badge bg-warning-subtle text-warning-emphasis border border-warning-subtle small" title={$t('readMismatchTitle')}>
                  <i class="bi bi-exclamation-triangle me-1"></i>{$t('readMismatch')}
                </span>
              {:else if result.readBack === 'unsupported'}
                <span class="badge bg-body-secondary text-body-secondary border small" title={$t('readUnsupportedTitle')}>
                  {$t('readUnsupported')}
                </span>
              {/if}
            </div>
          {:else if result?.error}
            <div class="alert alert-danger w-100 text-center py-4 my-auto">
              <i class="bi bi-exclamation-octagon fs-2 d-block mb-2 text-danger"></i>
              <strong class="d-block mb-1">{$t('invalidInput', { type: selectedType.toUpperCase() })}</strong>
              <small class="text-body-secondary">{formatError(result, $lang)}</small>
            </div>
          {:else}
            <div class="text-center text-body-secondary py-5 my-auto">
              <i class="bi bi-upc-scan fs-1 d-block mb-2 opacity-50"></i>
              <span>{$t('emptyHint')}</span>
            </div>
          {/if}
        </div>

        <!-- Export Buttons -->
        {#if result?.success}
          <div class="card-footer bg-body border-top p-3">
            <div class="row g-2">
              <div class="col-6">
                <button type="button" class="btn btn-outline-primary btn-sm w-100" on:click={copySVG}>
                  <i class="bi bi-clipboard me-1"></i> {$t('copySvg')}
                </button>
              </div>
              <div class="col-6">
                <button type="button" class="btn btn-outline-primary btn-sm w-100" on:click={copyPNG}>
                  <i class="bi bi-image me-1"></i> {$t('copyPng')}
                </button>
              </div>
              <div class="col-6">
                <button type="button" class="btn btn-primary btn-sm w-100" on:click={saveSVG}>
                  <i class="bi bi-download me-1"></i> {$t('saveSvg')}
                </button>
              </div>
              <div class="col-6">
                <button type="button" class="btn btn-primary btn-sm w-100" on:click={savePNG}>
                  <i class="bi bi-download me-1"></i> {$t('savePng')}
                </button>
              </div>
              {#if onSendToPrint}
                <div class="col-12 mt-2">
                  <button type="button" class="btn btn-success btn-sm w-100" on:click={sendToLabelPrint}>
                    <i class="bi bi-printer me-1"></i> {$t('addToLabels')}
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
