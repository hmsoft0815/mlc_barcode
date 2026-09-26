<script lang="ts">
  export let activeTab: 'single' | 'batch' | 'print' | 'check' | 'about' = 'single';
  export let appVersion: string = __APP_VERSION__;
  export let theme: 'light' | 'dark' = 'light';
  export let onToggleTheme: () => void = () => {};
  import { lang, setLang } from '../i18n/lang';
  import { t } from '../i18n/text/shell';
</script>

<header class="navbar navbar-expand bg-body border-bottom shadow-sm px-3 py-2 sticky-top no-print">
  <div class="container-fluid d-flex align-items-center justify-content-between">
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
        <button
          class="nav-link {activeTab === 'single' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'single')}
          type="button"
        >
          <i class="bi bi-upc-scan me-1"></i> {$t('tabSingle')}
        </button>

        <button
          class="nav-link {activeTab === 'batch' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'batch')}
          type="button"
        >
          <i class="bi bi-collection me-1"></i> {$t('tabBatch')}
        </button>

        <button
          class="nav-link {activeTab === 'print' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'print')}
          type="button"
        >
          <i class="bi bi-printer me-1"></i> {$t('tabPrint')}
        </button>

        <button
          class="nav-link {activeTab === 'check' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'check')}
          type="button"
        >
          <i class="bi bi-qr-code-scan me-1"></i> {$t('tabCheck')}
        </button>

        <button
          class="nav-link {activeTab === 'about' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'about')}
          type="button"
        >
          <i class="bi bi-info-circle me-1"></i> {$t('tabHelp')}
        </button>
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

      <!-- Theme Switcher -->
      <button
        class="btn btn-sm btn-outline-secondary px-2 ms-1"
        on:click={onToggleTheme}
        title={theme === 'dark' ? $t('toLight') : $t('toDark')}
        type="button"
        aria-label={$t('toggleTheme')}
      >
        <i class="bi bi-{theme === 'dark' ? 'sun-fill text-primary' : 'moon-stars-fill'}"></i>
      </button>
    </div>
  </div>
</header>

<style>
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
