<script lang="ts">
  // For a medicine code with a PZN: look it up in the public medicine
  // database of the BfArM (AMIce). It has no link that searches directly —
  // a session and accepting its terms come first — so a short guide shows
  // the three steps, the PZN is copied, and "open search" goes there.
  import type { DecodedCode } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import { copyText, openLink } from '../mobile/native';
  import { t } from '../i18n/text/lookup';

  export let code: DecodedCode;

  const AMICE_URL = 'https://auth.bfarm.de/sso-home/amguifree';
  $: pzn = code.content?.kind === 'pharma' ? code.content.fields?.pzn ?? '' : '';
  let open = false;
  let copied = false;

  async function copy() {
    await copyText(pzn).catch(() => {});
    copied = true;
    setTimeout(() => (copied = false), 1500);
  }

  async function start() {
    open = true;
    await copy();
  }
</script>

{#if pzn}
  <button type="button" class="btn btn-sm btn-outline-secondary" on:click={start}>
    <i class="bi bi-search me-1"></i>{$t('button')}
  </button>
{/if}

{#if open}
  <div class="lookup-backdrop" role="presentation" on:click|self={() => (open = false)}>
    <div class="card shadow-lg lookup-card" role="dialog" aria-modal="true" aria-label={$t('title')}>
      <div class="card-header d-flex justify-content-between align-items-center">
        <h6 class="mb-0"><i class="bi bi-search me-1 text-primary"></i>{$t('title')}</h6>
        <button type="button" class="btn-close" aria-label={$t('close')} on:click={() => (open = false)}></button>
      </div>
      <div class="card-body">
        <div class="small text-body-secondary mb-1">{$t('pznLabel')}</div>
        <div class="d-flex align-items-center gap-2 mb-3">
          <span class="pzn">{pzn}</span>
          <button type="button" class="btn btn-sm btn-outline-secondary" on:click={copy}>
            <i class="bi {copied ? 'bi-check2' : 'bi-clipboard'} me-1"></i>{copied ? $t('copied') : $t('copyAgain')}
          </button>
        </div>
        <p class="small mb-2">{$t('intro')}</p>
        <ol class="small ps-3 mb-3">
          <li class="mb-1">{$t('step1')}</li>
          <li class="mb-1">{$t('step2')}</li>
          <li>{$t('step3')}</li>
        </ol>
        <p class="small text-body-secondary mb-0">{$t('note')}</p>
      </div>
      <div class="card-footer d-flex justify-content-end gap-2">
        <button type="button" class="btn btn-sm btn-outline-secondary" on:click={() => (open = false)}>{$t('close')}</button>
        <button type="button" class="btn btn-sm btn-primary" on:click={() => openLink(AMICE_URL)}>
          <i class="bi bi-box-arrow-up-right me-1"></i>{$t('open')}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .lookup-backdrop {
    position: fixed;
    inset: 0;
    z-index: 2000;
    background: rgba(0, 0, 0, 0.55);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }
  .lookup-card {
    width: min(520px, 100%);
    max-height: 90vh;
    overflow-y: auto;
  }
  .pzn {
    font-family: var(--bs-font-monospace);
    font-size: 1.6rem;
    font-weight: 700;
    letter-spacing: 0.08em;
  }
</style>
