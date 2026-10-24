<script lang="ts">
  // The title bar of a mobile view: stays at the top while the page
  // scrolls, same surface as the tab bar. Slot "lead" (e.g. a back button)
  // sits left of the title, the default slot holds actions on the right.
  export let title: string;
  export let icon = ''; // Bootstrap Icons class, e.g. "bi-clock-history"
</script>

<header>
  <slot name="lead" />
  <h1>{#if icon}<i class="bi {icon}"></i>{/if}{title}</h1>
  <div class="actions"><slot /></div>
</header>

<style>
  header {
    position: sticky;
    top: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-height: 3.25rem;
    padding: 0.5rem 1rem;
    background: var(--bs-tertiary-bg);
    border-bottom: 1px solid var(--bs-border-color);
  }
  h1 {
    flex: 1;
    margin: 0;
    font-size: 1.15rem;
    font-weight: 600;
    letter-spacing: 0.01em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  h1 i {
    color: var(--bs-primary);
    margin-right: 0.5rem;
  }
  .actions {
    display: flex;
    gap: 0.5rem;
  }
</style>
