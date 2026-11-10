<script lang="ts">
  export let activeTab: 'single' | 'batch' | 'print' | 'check' | 'about' = 'single';
  export let appVersion: string = __APP_VERSION__;
  export let theme: 'light' | 'dark' = 'light';
  export let onToggleTheme: () => void = () => {};
  export let onOpenProfile: () => void = () => {};
  import { lang, setLang } from '../i18n/lang';
  import { t } from '../i18n/text/shell';
  import { System } from '@wailsio/runtime';
  import { onMount, tick } from 'svelte';

  // The macOS window draws its close/minimise/zoom buttons over the top left
  // of the page (TitleBarHiddenInset); leave them room there.
  const macInset = System.IsMac();

  type Tab = typeof activeTab;
  const tabs: { id: Tab; icon: string; label: 'tabSingle' | 'tabBatch' | 'tabPrint' | 'tabCheck' | 'tabHelp' }[] = [
    { id: 'single', icon: 'upc-scan', label: 'tabSingle' },
    { id: 'batch', icon: 'collection', label: 'tabBatch' },
    { id: 'print', icon: 'printer', label: 'tabPrint' },
    { id: 'check', icon: 'qr-code-scan', label: 'tabCheck' },
    { id: 'about', icon: 'info-circle', label: 'tabHelp' },
  ];

  // When the labels do not fit in one row, the tabs show only their icons
  // and name themselves in a tooltip. The width they need is measured with
  // the labels shown, again whenever the language changes.
  let bar: HTMLElement;
  let compact = false;
  let need = 0;

  function fit() {
    if (bar) compact = bar.clientWidth < need;
  }

  async function measure() {
    compact = false;
    await tick();
    if (!bar) return;
    need = [...bar.children].reduce((w, el) => w + (el as HTMLElement).scrollWidth, 0) + 16;
    fit();
  }

  onMount(() => {
    const ro = new ResizeObserver(fit);
    ro.observe(bar);
    return () => ro.disconnect();
  });

  $: $lang, measure();
</script>

<header class="navbar navbar-expand bg-body border-bottom shadow-sm px-3 py-2 sticky-top no-print" class:mac-inset={macInset}>
  <div class="container-fluid d-flex align-items-center justify-content-between" bind:this={bar}>
    <!-- Brand -->
    <div class="d-flex align-items-center gap-2">
      <img src="/appicon.png" width="30" height="30" class="rounded shadow-sm" alt={$t('logoAlt')} />
      <div>
        <span class="fw-bold fs-5 text-body">MLC Barcode</span>
        <span class="badge bg-primary-subtle text-primary border border-primary-subtle ms-1">v{appVersion}</span>
      </div>
    </div>

    <!-- Right side: Navigation Tabs & Theme Toggle -->
    <div class="d-flex align-items-center gap-2">
      <nav class="nav nav-pills gap-1">
        {#each tabs as tab (tab.id)}
          <button
            class="nav-link {activeTab === tab.id ? 'active' : 'text-body-secondary'}"
            on:click={() => (activeTab = tab.id)}
            type="button"
            title={compact ? $t(tab.label) : undefined}
            aria-label={compact ? $t(tab.label) : undefined}
          >
            <i class="bi bi-{tab.icon}" class:me-1={!compact}></i>{#if !compact}&nbsp;{$t(tab.label)}{/if}
          </button>
        {/each}
      </nav>

      <!-- Language switch (stored per device, see i18n/lang.ts) -->
      <div class="btn-group btn-group-sm ms-1" role="group" aria-label={$t('language')} title={$t('language')}>
        <button
          type="button"
          class="btn {$lang === 'de' ? 'btn-primary' : 'btn-outline-secondary'} lang-btn"
          aria-pressed={$lang === 'de'}
          title="Deutsch"
          on:click={() => setLang('de')}>DE</button
        >
        <button
          type="button"
          class="btn {$lang === 'en' ? 'btn-primary' : 'btn-outline-secondary'} lang-btn"
          aria-pressed={$lang === 'en'}
          title="English"
          on:click={() => setLang('en')}>EN</button
        >
      </div>

      <!-- My details (profile) -->
      <button
        type="button"
        class="btn btn-sm btn-outline-secondary px-2 ms-1"
        on:click={onOpenProfile}
        title={$t('myDetails')}
        aria-label={$t('myDetails')}
      >
        <i class="bi bi-person-gear"></i>
      </button>

      <!-- Theme Switcher -->
      <button
        type="button"
        class="btn btn-sm btn-outline-secondary px-2 ms-1"
        on:click={onToggleTheme}
        title={theme === 'dark' ? $t('toLight') : $t('toDark')}
        aria-label={$t('toggleTheme')}
      >
        <i class="bi bi-{theme === 'dark' ? 'sun-fill text-primary' : 'moon-stars-fill'}"></i>
      </button>
    </div>
  </div>
</header>

<style>
  .mac-inset {
    padding-left: 84px !important;
  }
  /* one row: when it gets narrow the tabs drop their labels (see fit) */
  .container-fluid,
  .container-fluid > div,
  .nav-pills {
    flex-wrap: nowrap;
    white-space: nowrap;
  }
  .mac-inset .nav-pills .nav-link {
    padding: 0.35rem 0.55rem;
  }
  .nav-pills .nav-link {
    font-weight: 500;
    padding: 0.35rem 0.85rem;
    border-radius: 0.375rem;
    transition: all 0.15s ease-in-out;
  }
  .lang-btn {
    font-size: 0.75rem;
    font-weight: 600;
    padding: 0.25rem 0.5rem;
  }
</style>
