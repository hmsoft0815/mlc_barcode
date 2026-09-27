<script lang="ts">
  // The last scans, newest first — e.g. the home pharmacy pack by pack.
  import { history, removeScan, clearHistory, type HistoryEntry } from './history';
  import { BARCODE_TYPES } from '../types';
  import { formatFieldValue, isExpired, kindLabel } from '../i18n/fields';
  import { lang, t } from '../i18n/lang';
  import type { Lang } from '../i18n/errors';
  import CodeResult from './CodeResult.svelte';

  let open = '';
  let confirmClear = false;

  const typeName = (id: string) => BARCODE_TYPES.find((t) => t.id === id)?.name ?? id.toUpperCase();

  // One line that says what the scan was.
  function summary(e: HistoryEntry, l: Lang): string {
    const f = e.code.content?.fields ?? {};
    if (e.code.content?.kind === 'pharma') {
      return [f.pzn && `PZN ${f.pzn}`, f.expiry && formatFieldValue('expiry', f.expiry, l)].filter(Boolean).join(' · ') || e.code.text;
    }
    return f.url ?? f.ssid ?? f.name ?? f.iban ?? e.code.text;
  }

  const expired = (e: HistoryEntry) => e.code.content?.kind === 'pharma' && isExpired(e.code.content.fields?.expiry);

  // When it was read: always date and time ("27.09.26, 14:05").
  function when(at: number, l: Lang): string {
    return new Date(at).toLocaleString(l === 'en' ? 'en-GB' : 'de-DE', {
      day: '2-digit',
      month: '2-digit',
      year: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
</script>

<div class="page">
  <div class="d-flex justify-content-between align-items-center mb-3">
    <h5 class="mb-0">{$t('historyTitle')}</h5>
    {#if $history.length > 0}
      {#if confirmClear}
        <div class="d-flex gap-2">
          <button class="btn btn-sm btn-danger" on:click={() => ((confirmClear = false), clearHistory())}>{$t('clearAll')}</button>
          <button class="btn btn-sm btn-outline-secondary" on:click={() => (confirmClear = false)}>{$t('cancel')}</button>
        </div>
      {:else}
        <button class="btn btn-sm btn-outline-secondary" on:click={() => (confirmClear = true)}>
          <i class="bi bi-trash me-1"></i> {$t('clear')}
        </button>
      {/if}
    {/if}
  </div>

  {#if $history.length === 0}
    <div class="text-center text-body-secondary py-5">
      <i class="bi bi-clock-history fs-1 d-block mb-2 opacity-50"></i>
      {$t('historyEmpty')}
    </div>
  {:else}
    <div class="list-group">
      {#each $history as e (e.id)}
        <button
          class="list-group-item list-group-item-action d-flex align-items-center gap-3 py-2"
          on:click={() => (open = open === e.id ? '' : e.id)}
        >
          <i class="bi {e.code.content?.kind === 'pharma' ? 'bi-capsule' : 'bi-upc-scan'} fs-4 {expired(e) ? 'text-danger' : 'text-primary'}"></i>
          <div class="flex-grow-1 text-start overflow-hidden">
            <div class="small text-body-secondary">
              {e.code.content ? kindLabel(e.code.content.kind, $lang) : typeName(e.code.type)} · {when(e.at, $lang)}
            </div>
            <div class="text-truncate {expired(e) ? 'text-danger fw-semibold' : ''}">{summary(e, $lang)}</div>
          </div>
          <i class="bi {open === e.id ? 'bi-chevron-up' : 'bi-chevron-down'} text-body-secondary"></i>
        </button>
        {#if open === e.id}
          <div class="list-group-item">
            <CodeResult code={e.code} />
            <button class="btn btn-sm btn-link text-danger p-0" on:click={() => removeScan(e.id)}>
              <i class="bi bi-trash me-1"></i> {$t('removeEntry')}
            </button>
          </div>
        {/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .page {
    padding: 1rem;
    max-width: 720px;
    margin: 0 auto;
  }
</style>
