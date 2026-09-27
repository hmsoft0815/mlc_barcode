<script lang="ts">
  // "Über" in the mobile UI: what the app is, version, licence, links.
  import licenseText from '../../../../LICENSE?raw';
  import thirdPartyText from '../../../../THIRD_PARTY_NOTICES.txt?raw';
  import { openLink } from './native';
  import { lang, setLang, t } from '../i18n/lang';
  import ProfileForm from '../components/ProfileForm.svelte';

  let editingProfile = false;

  export let appVersion: string;
  export let desktop = false; // offer the way back to the desktop UI
  export let onDesktopUI: () => void = () => {};

  const PRODUCT_URL = 'https://mlcgo.eu/products/mlc-barcode/';
  const SOURCE_URL = 'https://github.com/hmsoft0815/mlc_barcode';
</script>

{#if editingProfile}
<div class="page">
  <button class="btn btn-sm btn-link px-0 mb-2" on:click={() => (editingProfile = false)}>
    <i class="bi bi-chevron-left"></i> {$t('back')}
  </button>
  <h5 class="mb-3">{$t('myDetails')}</h5>
  <ProfileForm compact />
</div>
{:else}
<div class="page">
  <div class="text-center mb-4">
    <img src="/appicon.png" alt="" width="84" height="84" class="rounded-4 mb-2" />
    <h4 class="mb-0">MLC Barcode</h4>
    <span class="badge bg-primary-subtle text-primary mt-1">v{appVersion}</span>
    <p class="text-body-secondary small mt-3 mb-0">{$t('aboutText')}</p>
  </div>

  <div class="list-group mb-4">
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => (editingProfile = true)}>
      <i class="bi bi-person-gear text-primary"></i> {$t('myDetails')} <i class="bi bi-chevron-right ms-auto small"></i>
    </button>
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => openLink(PRODUCT_URL)}>
      <i class="bi bi-globe text-primary"></i> {$t('productPage')} <i class="bi bi-box-arrow-up-right ms-auto small"></i>
    </button>
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => openLink(SOURCE_URL)}>
      <i class="bi bi-github text-primary"></i> {$t('sourceCode')} <i class="bi bi-box-arrow-up-right ms-auto small"></i>
    </button>
    {#if desktop}
      <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={onDesktopUI}>
        <i class="bi bi-display text-primary"></i> {$t('desktopUI')}
      </button>
    {/if}
  </div>

  <h6>{$t('language')}</h6>
  <div class="btn-group mb-4" role="group">
    <button class="btn btn-sm {$lang === 'de' ? 'btn-primary' : 'btn-outline-secondary'}" on:click={() => setLang('de')}>Deutsch</button>
    <button class="btn btn-sm {$lang === 'en' ? 'btn-primary' : 'btn-outline-secondary'}" on:click={() => setLang('en')}>English</button>
  </div>

  <h6>{$t('license')}</h6>
  <p class="small text-body-secondary">{$t('licenseSummary')}</p>
  <details class="mb-2">
    <summary class="small text-primary">{$t('licenseText')}</summary>
    <pre class="license">{licenseText}</pre>
  </details>
  <details>
    <summary class="small text-primary">{$t('thirdParty')}</summary>
    <pre class="license">{thirdPartyText}</pre>
  </details>
</div>
{/if}

<style>
  .page {
    padding: 1.25rem 1rem;
    max-width: 720px;
    margin: 0 auto;
  }
  .license {
    white-space: pre-wrap;
    font-size: 0.7rem;
    max-height: 16rem;
    overflow-y: auto;
    background: var(--bs-secondary-bg);
    border-radius: 0.5rem;
    padding: 0.5rem;
    margin-top: 0.5rem;
  }
</style>
