<script lang="ts">
  // One read code on the mobile result sheet and in the history: what it
  // is, its fields, and what one can do with it.
  import type { DecodedCode } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import CodeDetails from '../components/CodeDetails.svelte';
  import { BARCODE_TYPES } from '../types';
  import { kindLabel } from '../i18n/fields';
  import { lang, t } from '../i18n/lang';
  import { copyText, linkOf, openLink, shareText } from './native';

  export let code: DecodedCode;
  export let index = 0; // shown when a photo had several codes

  let feedback = '';
  $: link = linkOf(code.content?.fields?.url ?? code.text);
  $: typeName = BARCODE_TYPES.find((t) => t.id === code.type)?.name ?? code.type.toUpperCase();

  async function copy() {
    await copyText(code.text);
    flash($t('copied'));
  }

  async function share() {
    if ((await shareText(code.text)) === 'copied') flash($t('copiedNoShare'));
  }

  function flash(text: string) {
    feedback = text;
    setTimeout(() => (feedback = ''), 1800);
  }
</script>

<div class="border rounded-3 p-3 mb-3 bg-body-tertiary">
  <div class="d-flex align-items-center gap-2 mb-2 flex-wrap">
    {#if index}<span class="badge bg-primary">{index}</span>{/if}
    <span class="fw-semibold">{typeName}</span>
    {#if code.content}
      <span class="badge bg-info-subtle text-info-emphasis">{kindLabel(code.content.kind, $lang)}</span>
    {/if}
  </div>
  <CodeDetails {code} lang={$lang} />
  <div class="d-flex gap-2 mt-3 flex-wrap">
    {#if link}
      <button class="btn btn-sm btn-primary" on:click={() => openLink(link)}>
        <i class="bi bi-box-arrow-up-right me-1"></i> {$t('open')}
      </button>
    {/if}
    <button class="btn btn-sm btn-outline-secondary" on:click={copy}>
      <i class="bi bi-clipboard me-1"></i> {$t('copy')}
    </button>
    <button class="btn btn-sm btn-outline-secondary" on:click={share}>
      <i class="bi bi-share me-1"></i> {$t('share')}
    </button>
    {#if feedback}<span class="small text-success align-self-center">{feedback}</span>{/if}
  </div>
</div>
