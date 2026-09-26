<script lang="ts">
  // Creating on the phone: the few codes one shows or sends to someone —
  // a link, the Wi-Fi, one's contact, a GiroCode. Shared, not saved.
  import {
    FormatEPC,
    FormatVCard,
    FormatWifi,
    GenerateBarcode,
  } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import { formatError } from '../i18n/errors';
  import { copyText, haptic, shareImage } from './native';
  import { lang, t } from '../i18n/lang';
  import type { UIKey } from '../i18n/ui';

  type Kind = 'link' | 'wifi' | 'contact' | 'giro';
  const KINDS: { id: Kind; label: UIKey; icon: string }[] = [
    { id: 'link', label: 'kindLink', icon: 'bi-link-45deg' },
    { id: 'wifi', label: 'kindWifi', icon: 'bi-wifi' },
    { id: 'contact', label: 'kindContact', icon: 'bi-person-vcard' },
    { id: 'giro', label: 'kindGiro', icon: 'bi-bank' },
  ];

  let kind: Kind = 'link';
  let text = '';
  let wifi = { ssid: '', password: '', encryption: 'WPA', hidden: false };
  let contact = { firstName: '', lastName: '', phone: '', email: '' };
  let giro = { name: '', iban: '', amount: '', reference: '' };

  let png = '';
  let data = '';
  let error = '';
  let busy = false;
  let feedback = '';
  let fullscreen = false;

  $: kind, reset();
  function reset() {
    png = '';
    error = '';
  }

  async function payload(): Promise<string> {
    switch (kind) {
      case 'link':
        return text.trim();
      case 'wifi':
        return FormatWifi(wifi);
      case 'contact':
        return FormatVCard(contact);
      case 'giro':
        return FormatEPC({
          name: giro.name,
          iban: giro.iban,
          bic: '',
          amount: Number(giro.amount.replace(',', '.')) || 0,
          reference: giro.reference,
          purpose: '',
        });
    }
  }

  async function create() {
    busy = true;
    error = '';
    png = '';
    try {
      data = await payload();
      const r = await GenerateBarcode({
        type: 'qr',
        data,
        customText: '',
        width: 720,
        height: 720,
        showText: false,
        fontSize: 0,
        noQuietZone: false,
        foregroundColor: '#000000',
        backgroundColor: '#ffffff',
      });
      if (r.success && r.pngData) {
        png = r.pngData;
        haptic('success');
      } else {
        error = formatError(r, $lang);
        haptic('error');
      }
    } catch (e: any) {
      error = e?.message ?? String(e);
    } finally {
      busy = false;
    }
  }

  async function share() {
    const name = { link: 'link', wifi: `wlan-${wifi.ssid}`, contact: 'kontakt', giro: 'girocode' }[kind];
    if ((await shareImage(png, name, data)) === 'copied') flash($t('contentCopiedNoShare'));
  }

  async function copy() {
    await copyText(data);
    flash($t('contentCopied'));
  }

  function flash(t: string) {
    feedback = t;
    setTimeout(() => (feedback = ''), 2000);
  }
</script>

<div class="page">
  <h5 class="mb-3">{$t('createTitle')}</h5>

  <div class="kinds mb-3">
    {#each KINDS as k}
      <button class="btn btn-sm {kind === k.id ? 'btn-primary' : 'btn-outline-secondary'}" on:click={() => (kind = k.id)}>
        <i class="bi {k.icon} me-1"></i>{$t(k.label)}
      </button>
    {/each}
  </div>

  <form on:submit|preventDefault={create} class="mb-3">
    {#if kind === 'link'}
      <textarea class="form-control" rows="3" placeholder={$t('linkPlaceholder')} bind:value={text}></textarea>
    {:else if kind === 'wifi'}
      <input class="form-control mb-2" placeholder={$t('ssid')} bind:value={wifi.ssid} autocapitalize="off" />
      <input class="form-control mb-2" placeholder={$t('password')} bind:value={wifi.password} autocapitalize="off" autocomplete="off" />
      <div class="d-flex gap-2 align-items-center">
        <select class="form-select" bind:value={wifi.encryption}>
          <option value="WPA">{$t('encWpa')}</option>
          <option value="WEP">{$t('encWep')}</option>
          <option value="nopass">{$t('encOpen')}</option>
        </select>
        <div class="form-check text-nowrap">
          <input class="form-check-input" type="checkbox" id="wifi-hidden" bind:checked={wifi.hidden} />
          <label class="form-check-label small" for="wifi-hidden">{$t('hiddenNet')}</label>
        </div>
      </div>
    {:else if kind === 'contact'}
      <div class="d-flex gap-2 mb-2">
        <input class="form-control" placeholder={$t('firstName')} bind:value={contact.firstName} />
        <input class="form-control" placeholder={$t('lastName')} bind:value={contact.lastName} />
      </div>
      <input class="form-control mb-2" type="tel" placeholder={$t('phone')} bind:value={contact.phone} />
      <input class="form-control" type="email" placeholder={$t('email')} bind:value={contact.email} autocapitalize="off" />
    {:else}
      <input class="form-control mb-2" placeholder={$t('recipient')} bind:value={giro.name} />
      <input class="form-control mb-2" placeholder={$t('iban')} bind:value={giro.iban} autocapitalize="characters" />
      <div class="d-flex gap-2">
        <input class="form-control" inputmode="decimal" placeholder={$t('amountOptional')} bind:value={giro.amount} />
        <input class="form-control" placeholder={$t('reference')} bind:value={giro.reference} />
      </div>
    {/if}
    <button type="submit" class="btn btn-primary w-100 mt-3" disabled={busy}>
      {#if busy}<span class="spinner-border spinner-border-sm me-1"></span>{:else}<i class="bi bi-qr-code me-1"></i>{/if}
      {$t('createQr')}
    </button>
  </form>

  {#if error}
    <div class="alert alert-warning small"><i class="bi bi-exclamation-triangle me-1"></i> {error}</div>
  {/if}

  {#if png}
    <div class="text-center">
      <button class="qr" on:click={() => (fullscreen = true)} aria-label={$t('enlarge')}>
        <img src={png} alt="QR-Code" />
      </button>
      <div class="small text-body-secondary mt-1">{$t('tapToEnlarge')}</div>
      <div class="d-flex gap-2 justify-content-center mt-3 flex-wrap">
        <button class="btn btn-primary" on:click={share}><i class="bi bi-share me-1"></i> {$t('share')}</button>
        <button class="btn btn-outline-secondary" on:click={copy}><i class="bi bi-clipboard me-1"></i> {$t('copyContent')}</button>
      </div>
      {#if feedback}<div class="small text-success mt-2">{feedback}</div>{/if}
    </div>
  {/if}
</div>

{#if fullscreen}
  <button class="fullscreen" on:click={() => (fullscreen = false)} aria-label={$t('tapToClose')}>
    <img src={png} alt="QR-Code" />
    <span class="small text-dark mt-3">{$t('tapToClose')}</span>
  </button>
{/if}

<style>
  .page {
    padding: 1rem;
    max-width: 720px;
    margin: 0 auto;
  }
  .kinds {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  .qr {
    border: none;
    padding: 0.75rem;
    background: #fff;
    border-radius: 1rem;
  }
  .qr img {
    width: min(70vw, 320px);
    display: block;
  }
  .fullscreen {
    position: fixed;
    inset: 0;
    z-index: 2000;
    border: none;
    background: #fff;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
  .fullscreen img {
    width: min(92vw, 92vh);
  }
</style>
