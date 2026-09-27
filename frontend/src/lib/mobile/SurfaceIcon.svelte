<script lang="ts">
  // Symbols of the scan view's shape switch: what the code sits on.
  export let kind: 'auto' | 'flat' | 'bottle' | 'tube';
</script>

<svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" aria-hidden="true">
  {#if kind === 'auto'}
    <text x="12" y="17" text-anchor="middle" font-size="14" font-weight="700" font-family="Inter, system-ui, sans-serif" fill="currentColor" stroke="none">A</text>
  {:else if kind === 'flat'}
    <!-- a flat pack with a code on its face -->
    <rect x="3.5" y="5" width="17" height="14" rx="1.5" />
    <rect x="9" y="9" width="6" height="6" fill="currentColor" stroke="none" />
  {:else if kind === 'bottle'}
    <!-- a bottle: neck, shoulders, wide body; the code about half its width -->
    <path d="M10 2.5h4v3.5l3.5 3v11.5a1.5 1.5 0 0 1-1.5 1.5H8a1.5 1.5 0 0 1-1.5-1.5V9l3.5-3z" />
    <path d="M9.5 12.5q2.5 .8 5 0v4.5q-2.5 .8-5 0z" fill="currentColor" stroke="none" />
  {:else}
    <!-- a thin tube: the code wraps far around, its sides curve away -->
    <path d="M9 2.5h6l-.5 2h-5z" />
    <rect x="8.5" y="4.5" width="7" height="17" rx="2" />
    <path d="M8.5 11q3.5 1.4 7 0v5q-3.5 1.4-7 0z" fill="currentColor" stroke="none" />
  {/if}
</svg>
