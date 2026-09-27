<script lang="ts">
  // For a medicine code with a PZN: copy the PZN and open the public
  // medicine database of the BfArM (AMIce). It has no link that searches
  // directly — a session and accepting its terms come first — so the user
  // pastes the PZN there.
  import type { DecodedCode } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import { copyText, openLink } from '../mobile/native';
  import { t } from '../i18n/text/lookup';

  export let code: DecodedCode;

  const AMICE_URL = 'https://auth.bfarm.de/sso-home/amguifree';
  $: pzn = code.content?.kind === 'pharma' ? code.content.fields?.pzn ?? '' : '';
  let note = '';

  async function lookup() {
    await copyText(pzn).catch(() => {});
    note = $t('copied', { pzn });
    await openLink(AMICE_URL);
  }
</script>

{#if pzn}
  <button type="button" class="btn btn-sm btn-outline-secondary" on:click={lookup}>
    <i class="bi bi-search me-1"></i>{$t('button')}
  </button>
  {#if note}<div class="small text-body-secondary mt-1 w-100">{note}</div>{/if}
{/if}
