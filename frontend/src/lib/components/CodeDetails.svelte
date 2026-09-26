<script lang="ts">
  // The fields of one read code (parsed payload) and its raw content; used
  // by the desktop "Prüfen" tab and the mobile scan result and history.
  import type { DecodedCode } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import { fieldLabel, formatFieldValue, isCheckKey, isExpired, sortedFieldKeys } from '../i18n/fields';
  import type { Lang } from '../i18n/errors';

  export let code: DecodedCode;
  export let lang: Lang = 'de';
  export let rawOpen: boolean | undefined = undefined; // default: open when nothing was parsed

  const bad = (key: string, value: string | undefined) =>
    (isCheckKey(key) && value === 'false') || (key === 'expiry' && isExpired(value));
</script>

{#if code.content?.fields}
  <table class="table table-sm mb-2">
    <tbody>
      {#each sortedFieldKeys(code.content.fields) as key}
        <tr>
          <th class="text-body-secondary fw-normal small" style="width: 40%">{fieldLabel(key, lang)}</th>
          <td class="small {bad(key, code.content.fields[key]) ? 'text-danger fw-semibold' : ''}">
            <span class="value">{formatFieldValue(key, code.content.fields[key], lang)}</span>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}
<details open={rawOpen ?? !code.content}>
  <summary class="small text-body-secondary">{lang === 'en' ? 'Raw content' : 'Rohinhalt'}</summary>
  <pre class="raw bg-body-secondary border rounded p-2 mt-1 mb-0">{code.text}</pre>
</details>

<style>
  .raw {
    white-space: pre-wrap;
    word-break: break-all;
    font-size: 0.75rem;
    max-height: 12rem;
    overflow-y: auto;
  }
  .value {
    white-space: pre-wrap;
    word-break: break-word;
  }
</style>
