<script lang="ts">
  // The mobile UI (phones and tablets): scanning first, then history,
  // creating and "Über". Shares engine, bindings, translations and theme
  // with the desktop UI; see main.ts for which one is shown.
  import { onMount } from 'svelte';
  import { GetVersion } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import ScanView from './ScanView.svelte';
  import HistoryView from './HistoryView.svelte';
  import CreateView from './CreateView.svelte';
  import AboutView from './AboutView.svelte';
  import { haptic } from './native';
  import { lang, t } from '../i18n/lang';
  import type { UIKey } from '../i18n/ui';

  export let desktop = false; // running on a desktop OS (switched here by hand)
  export let onDesktopUI: () => void = () => {};

  type Tab = 'scan' | 'history' | 'create' | 'about';
  const TABS: { id: Tab; label: UIKey; icon: string }[] = [
    { id: 'scan', label: 'tabScan', icon: 'bi-upc-scan' },
    { id: 'history', label: 'tabHistory', icon: 'bi-clock-history' },
    { id: 'create', label: 'tabCreate', icon: 'bi-qr-code' },
    { id: 'about', label: 'tabAbout', icon: 'bi-info-circle' },
  ];
  // #history, #create, #about open that tab directly.
  const fromHash = location.hash.slice(1) as Tab;
  let tab: Tab = TABS.some((x) => x.id === fromHash) ? fromHash : 'scan';
  let appVersion = __APP_VERSION__;

  function select(t: Tab) {
    if (t !== tab) haptic('selection');
    tab = t;
  }

  onMount(async () => {
    document.documentElement.setAttribute('data-bs-theme', 'dark');
    document.documentElement.lang = $lang;
    appVersion = (await GetVersion().catch(() => '')) || appVersion;
  });
</script>

<div class="mobile">
  <main class:full={tab === 'scan'}>
    <!-- The scan view stays mounted so the camera state survives a tab
         switch; it stops the camera itself while hidden. -->
    <div class="view" hidden={tab !== 'scan'}><ScanView active={tab === 'scan'} /></div>
    {#if tab === 'history'}<HistoryView />{/if}
    {#if tab === 'create'}<CreateView />{/if}
    {#if tab === 'about'}<AboutView {appVersion} {desktop} {onDesktopUI} />{/if}
  </main>

  <nav class="tabbar">
    {#each TABS as tb}
      <button class:active={tab === tb.id} on:click={() => select(tb.id)}>
        <i class="bi {tb.icon}"></i>
        <span>{$t(tb.label)}</span>
      </button>
    {/each}
  </nav>
</div>

<style>
  .mobile {
    position: fixed;
    inset: 0;
    display: flex;
    flex-direction: column;
    background: var(--bs-body-bg);
    padding-top: env(safe-area-inset-top);
    padding-left: env(safe-area-inset-left);
    padding-right: env(safe-area-inset-right);
  }
  main {
    position: relative;
    flex: 1;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
  }
  main.full {
    overflow: hidden;
  }
  .view {
    position: absolute;
    inset: 0;
  }
  .view[hidden] {
    display: none;
  }
  .tabbar {
    display: flex;
    border-top: 1px solid var(--bs-border-color);
    background: var(--bs-tertiary-bg);
    padding-bottom: env(safe-area-inset-bottom);
  }
  .tabbar button {
    flex: 1;
    border: none;
    background: none;
    color: var(--bs-secondary-color);
    padding: 0.5rem 0 0.4rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.1rem;
    font-size: 0.7rem;
  }
  .tabbar button i {
    font-size: 1.35rem;
  }
  .tabbar button.active {
    color: var(--bs-primary);
  }
</style>
