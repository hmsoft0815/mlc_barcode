<script lang="ts">
  // "Über" in the mobile UI: what the app is, version, licence, links.
  import licenseText from '../../../../LICENSE?raw';
  import thirdPartyText from '../../../../THIRD_PARTY_NOTICES.txt?raw';
  import { openLink } from './native';
  import { lang, setLang, t } from '../i18n/lang';
  import ProfileForm from '../components/ProfileForm.svelte';
  import PageHeader from './PageHeader.svelte';
  import { paragraphs } from './reflow';

  // A page opened from the list: own details or the licence texts.
  let sub: 'profile' | 'licenses' | null = null;

  export let appVersion: string;
  export let desktop = false; // offer the way back to the desktop UI
  export let onDesktopUI: () => void = () => {};

  const PRODUCT_URL = 'https://mlcgo.eu/products/mlc-barcode/';
  const SOURCE_URL = 'https://github.com/hmsoft0815/mlc_barcode';
  const SPONSOR_URL = 'https://www.hmsoftware.com';
</script>

<!-- LICENSE and the notices, re-wrapped for a narrow screen; the full text,
     only the line breaks differ -->
{#snippet license(text: string)}
  <div class="license">
    {#each paragraphs(text) as p}
      {#if p.kind === 'rule'}
        <hr />
      {:else if p.kind === 'gap'}
        <div class="gap"></div>
      {:else}
        <p style:margin-left="{p.level * 1.25 + (p.under ? 1.9 : 0)}em">
          {#if p.marker}<span class="marker">{p.marker}</span>{/if}<span>{p.text}</span>
        </p>
      {/if}
    {/each}
  </div>
{/snippet}

{#if sub === 'profile'}
<PageHeader title={$t('myDetails')}>
  <button slot="lead" class="btn btn-sm btn-link px-0 me-1" aria-label={$t('back')} on:click={() => (sub = null)}>
    <i class="bi bi-chevron-left fs-5"></i>
  </button>
</PageHeader>
<div class="page">
  <ProfileForm compact />
</div>
{:else if sub === 'licenses'}
<PageHeader title={$t('openSourceLicenses')}>
  <button slot="lead" class="btn btn-sm btn-link px-0 me-1" aria-label={$t('back')} on:click={() => (sub = null)}>
    <i class="bi bi-chevron-left fs-5"></i>
  </button>
</PageHeader>
<div class="page">
  <p class="small text-body-secondary">{$t('licenseSummary')}</p>
  <h6 class="mt-4">{$t('licenseText')}</h6>
  {@render license(licenseText)}
  <h6 class="mt-4">{$t('thirdParty')}</h6>
  {@render license(thirdPartyText)}
</div>
{:else}
<PageHeader title={$t('tabAbout')} icon="bi-info-circle" />
<div class="page">
  <div class="text-center mb-4">
    <img src="/appicon.png" alt="" width="84" height="84" class="rounded-4 mb-2" />
    <h4 class="mb-0">MLC Barcode</h4>
    <span class="badge bg-primary-subtle text-primary mt-1">v{appVersion}</span>
    <p class="text-body-secondary small mt-3 mb-0">{$t('aboutText')}</p>
    <p class="sponsor small mt-3 mb-0">
      {$t('sponsoredBy')}
      <button class="btn btn-link btn-sm p-0 align-baseline" on:click={() => openLink(SPONSOR_URL)}>H&amp;M Software-Netzwerktechnik-GmbH</button>
      <span class="text-body-secondary">(www.hmsoftware.com)</span>
    </p>
  </div>

  <div class="list-group mb-4">
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => (sub = 'profile')}>
      <i class="bi bi-person-gear text-primary"></i> {$t('myDetails')} <i class="bi bi-chevron-right ms-auto small"></i>
    </button>
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => openLink(PRODUCT_URL)}>
      <i class="bi bi-globe text-primary"></i> {$t('productPage')} <i class="bi bi-box-arrow-up-right ms-auto small"></i>
    </button>
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => openLink(SOURCE_URL)}>
      <i class="bi bi-github text-primary"></i> {$t('sourceCode')} <i class="bi bi-box-arrow-up-right ms-auto small"></i>
    </button>
    <button class="list-group-item list-group-item-action d-flex align-items-center gap-2" on:click={() => (sub = 'licenses')}>
      <i class="bi bi-file-earmark-text text-primary"></i> {$t('openSourceLicenses')} <i class="bi bi-chevron-right ms-auto small"></i>
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
</div>
{/if}

<style>
  .page {
    padding: 1.25rem 1rem;
    max-width: 720px;
    margin: 0 auto;
  }
  .sponsor {
    color: var(--bs-secondary-color);
  }
  .sponsor .btn-link {
    font-size: inherit;
    text-decoration: none;
  }
  .license {
    font-size: 0.75rem;
    line-height: 1.45;
    color: var(--bs-secondary-color);
    background: var(--bs-secondary-bg);
    border-radius: 0.5rem;
    padding: 0.75rem;
    margin-top: 0.5rem;
    overflow-wrap: anywhere; /* long URLs */
  }
  .license p {
    display: flex;
    margin: 0;
  }
  .license .marker {
    flex: none;
    min-width: 1.9em;
  }
  .license .gap {
    height: 0.6em;
  }
  .license hr {
    margin: 0.5rem 0;
  }
</style>
