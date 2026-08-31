<script lang="ts">
  import { onMount } from 'svelte';
  import Navbar from './lib/components/Navbar.svelte';
  import SingleGenerator from './lib/components/SingleGenerator.svelte';
  import BatchGenerator from './lib/components/BatchGenerator.svelte';
  import LabelPrinter from './lib/components/LabelPrinter.svelte';
  import AboutView from './lib/components/AboutView.svelte';
  import { GetVersion } from '../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';

  let activeTab: 'single' | 'batch' | 'print' | 'about' = 'single';
  let appVersion = '1.3.0';
  let theme: 'light' | 'dark' = 'light';

  let printItems: Array<{ data: string; svg: string; type: string }> = [];

  function applyTheme(newTheme: 'light' | 'dark') {
    theme = newTheme;
    document.documentElement.setAttribute('data-bs-theme', newTheme);
    try {
      localStorage.setItem('mlc_theme', newTheme);
    } catch (e) {
      // ignore in restricted environments
    }
  }

  function handleToggleTheme() {
    applyTheme(theme === 'dark' ? 'light' : 'dark');
  }

  function handleSendSingleToPrint(item: { data: string; svg: string; type: string }) {
    printItems = [item];
    activeTab = 'print';
  }

  function handleSendBatchToPrint(items: Array<{ data: string; svg: string; type: string }>) {
    printItems = items;
    activeTab = 'print';
  }

  onMount(async () => {
    // Theme preference
    const saved = localStorage.getItem('mlc_theme') as 'light' | 'dark' | null;
    if (saved === 'dark' || saved === 'light') {
      applyTheme(saved);
    } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      applyTheme('dark');
    } else {
      applyTheme('light');
    }

    try {
      const v = await GetVersion();
      if (v) appVersion = v;
    } catch (e) {
      console.warn('Could not fetch app version:', e);
    }
  });
</script>

<div class="app-root d-flex flex-column min-vh-100 bg-body-tertiary">
  <Navbar bind:activeTab {appVersion} {theme} onToggleTheme={handleToggleTheme} />

  <main class="flex-grow-1">
    {#if activeTab === 'single'}
      <SingleGenerator onSendToPrint={handleSendSingleToPrint} />
    {:else if activeTab === 'batch'}
      <BatchGenerator onSendBatchToPrint={handleSendBatchToPrint} />
    {:else if activeTab === 'print'}
      <LabelPrinter {printItems} />
    {:else if activeTab === 'about'}
      <AboutView {appVersion} />
    {/if}
  </main>
</div>

<style>
  .app-root {
    font-family: system-ui, -apple-system, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  }
</style>
