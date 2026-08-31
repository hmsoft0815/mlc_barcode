<script lang="ts">
  export let activeTab: 'single' | 'batch' | 'print' | 'about' = 'single';
  export let appVersion: string = '1.3.0';
  export let theme: 'light' | 'dark' = 'light';
  export let onToggleTheme: () => void = () => {};
</script>

<header class="navbar navbar-expand bg-body border-bottom shadow-sm px-3 py-2 sticky-top no-print">
  <div class="container-fluid d-flex align-items-center justify-content-between">
    <!-- Brand -->
    <div class="d-flex align-items-center gap-2">
      <img src="/appicon.png" width="30" height="30" class="rounded shadow-sm" alt="MLC Barcode Logo" />
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
          <i class="bi bi-upc-scan me-1"></i> Einzel-Generator
        </button>

        <button
          class="nav-link {activeTab === 'batch' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'batch')}
          type="button"
        >
          <i class="bi bi-collection me-1"></i> Batch-Generator
        </button>

        <button
          class="nav-link {activeTab === 'print' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'print')}
          type="button"
        >
          <i class="bi bi-printer me-1"></i> Etiketten-Druck
        </button>

        <button
          class="nav-link {activeTab === 'about' ? 'active' : 'text-body-secondary'}"
          on:click={() => (activeTab = 'about')}
          type="button"
        >
          <i class="bi bi-info-circle me-1"></i> Hilfe
        </button>
      </nav>

      <!-- Theme Switcher -->
      <button
        class="btn btn-sm btn-outline-secondary px-2 ms-1"
        on:click={onToggleTheme}
        title={theme === 'dark' ? 'Zu hellem Modus wechseln' : 'Zu dunklem Modus wechseln'}
        type="button"
        aria-label="Design umschalten"
      >
        <i class="bi bi-{theme === 'dark' ? 'sun-fill text-warning' : 'moon-stars-fill'}"></i>
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
</style>
