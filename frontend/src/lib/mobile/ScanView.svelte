<script lang="ts">
  // Start screen of the mobile UI: live camera, the codes found with their
  // fields, a photo from the gallery as the alternative.
  import { onDestroy, tick } from 'svelte';
  import { DecodeImage, SetTorch } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
  import type { DecodedCode } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';
  import { CameraScanner, cameraErrorText, cameraSupported, coverGuide, type ScanHit } from '../scan/camera';
  import { lang, t } from '../i18n/lang';
  import { formatError } from '../i18n/errors';
  import CodeResult from './CodeResult.svelte';
  import { addScan } from './history';
  import { haptic } from './native';

  export let active = false;

  let video: HTMLVideoElement;
  let frameEl: HTMLDivElement;
  let scanner = new CameraScanner(onHit, coverGuide(() => frameEl?.getBoundingClientRect() ?? null), onAttempt);
  // Blinks once per frame the decoder looks at.
  let attempts = 0;
  function onAttempt() {
    attempts++;
  }
  let running = false;
  let starting = false;
  let error = '';
  let codes: DecodedCode[] = [];
  let torch = false;
  let galleryBusy = false;
  let notice = '';

  // The camera runs only while this tab is shown and the app is in front.
  let hidden = typeof document !== 'undefined' && document.hidden;
  $: if (active && !hidden && !running && !starting && codes.length === 0 && !error) start();
  $: if ((!active || hidden) && running) stop();

  async function start() {
    if (!cameraSupported) {
      error = $t('camNoAccess');
      return;
    }
    starting = true;
    error = '';
    await tick();
    try {
      await scanner.start(video);
      running = true;
      scanner = scanner;
    } catch (e: any) {
      error = cameraErrorText(e, $lang);
    } finally {
      starting = false;
    }
  }

  function stop() {
    if (torch) toggleTorch();
    scanner.stop();
    running = false;
  }

  function onVisibility() {
    hidden = document.hidden;
  }
  document.addEventListener('visibilitychange', onVisibility);
  onDestroy(() => {
    document.removeEventListener('visibilitychange', onVisibility);
    stop();
  });

  function show(found: DecodedCode[]) {
    codes = found;
    haptic('success');
    for (const c of found) addScan(c);
  }

  function onHit(hit: ScanHit) {
    show(hit.result.codes ?? []);
  }

  function scanAgain() {
    codes = [];
    notice = '';
    if (running) scanner.resume();
    else start();
  }

  async function toggleTorch() {
    const want = !torch;
    if ((await scanner.setTorch(want)) || (await SetTorch(want).catch(() => false))) torch = want;
  }

  async function switchCamera() {
    try {
      await scanner.switchCamera();
      scanner = scanner;
    } catch (e: any) {
      error = cameraErrorText(e, $lang);
      stop();
    }
  }

  function fromGallery(file: File | null | undefined) {
    if (!file) return;
    galleryBusy = true;
    notice = '';
    const reader = new FileReader();
    reader.onload = async () => {
      try {
        const r = await DecodeImage(reader.result as string);
        if (r.success && (r.codes?.length ?? 0) > 0) {
          scanner.pause();
          show(r.codes!);
        } else {
          notice = formatError(r, $lang);
          haptic('warning');
        }
      } finally {
        galleryBusy = false;
      }
    };
    reader.readAsDataURL(file);
  }
</script>

<div class="scan">
  <!-- svelte-ignore a11y-media-has-caption -->
  <video bind:this={video} class:mirrored={scanner.facingUser} autoplay playsinline muted disablepictureinpicture></video>

  {#if running && codes.length === 0}
    <div class="frame" bind:this={frameEl} aria-hidden="true">
      {#key attempts}<span class="pulse"></span>{/key}
    </div>
    <div class="hint">{$t('holdInFrame')}</div>
  {/if}

  {#if error}
    <div class="message">
      <i class="bi bi-camera-video-off fs-1 d-block mb-2 opacity-75"></i>
      <p class="mb-3">{error}</p>
      {#if cameraSupported}
        <button class="btn btn-outline-light btn-sm" on:click={() => ((error = ''), start())}>{$t('retry')}</button>
      {/if}
    </div>
  {/if}

  <div class="tools top">
    {#if running}
      <button class="round" class:on={torch} on:click={toggleTorch} aria-label={$t('torch')}>
        <i class="bi {torch ? 'bi-lightbulb-fill' : 'bi-lightbulb'}"></i>
      </button>
      {#if scanner.cameras.length > 1}
        <button class="round" on:click={switchCamera} aria-label={$t('switchCamera')}>
          <i class="bi bi-arrow-repeat"></i>
        </button>
      {/if}
    {/if}
  </div>

  <div class="tools bottom">
    {#if notice}
      <div class="notice">{notice}</div>
    {/if}
    <label class="btn btn-light btn-sm rounded-pill px-3 shadow">
      {#if galleryBusy}
        <span class="spinner-border spinner-border-sm me-1"></span> {$t('reading')}
      {:else}
        <i class="bi bi-images me-1"></i> {$t('pickPhoto')}
      {/if}
      <input type="file" accept="image/*" class="d-none" on:change={(e) => fromGallery(e.currentTarget.files?.[0])} />
    </label>
  </div>

  {#if codes.length > 0}
    <div class="sheet">
      <div class="sheet-handle"></div>
      <div class="sheet-body">
        <div class="small text-body-secondary mb-2">{codes.length === 1 ? $t('foundOne') : $t('foundMany', { n: codes.length })}</div>
        {#each codes as code, i}
          <CodeResult {code} index={codes.length > 1 ? i + 1 : 0} />
        {/each}
      </div>
      <div class="sheet-actions">
        <button class="btn btn-primary w-100" on:click={scanAgain}>
          <i class="bi bi-upc-scan me-1"></i> {$t('scanAgain')}
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  .scan {
    position: absolute;
    inset: 0;
    background: #000;
    overflow: hidden;
  }
  video {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    pointer-events: none; /* a tap must not open the system video player */
  }
  video.mirrored {
    transform: scaleX(-1);
  }
  /* The guide frame; its box is what coverGuide sends every other frame. */
  .frame {
    position: absolute;
    left: 50%;
    top: 45%;
    width: min(80vw, 520px, 90vh);
    aspect-ratio: 4 / 3;
    transform: translate(-50%, -50%);
    border: 3px solid rgba(var(--bs-primary-rgb), 0.95);
    border-radius: 18px;
    box-shadow: 0 0 0 200vmax rgba(0, 0, 0, 0.35);
    pointer-events: none;
  }
  .pulse {
    position: absolute;
    top: 10px;
    right: 10px;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--bs-primary);
    animation: blink 0.28s ease-out forwards;
  }
  @keyframes blink {
    from {
      opacity: 1;
      transform: scale(1.3);
    }
    to {
      opacity: 0.15;
      transform: scale(1);
    }
  }
  .hint {
    position: absolute;
    left: 0;
    right: 0;
    top: calc(45% + min(30vw, 195px, 34vh) + 1.25rem);
    text-align: center;
    color: #fff;
    font-size: 0.9rem;
    text-shadow: 0 1px 3px #000;
  }
  .message {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;
    text-align: center;
    color: #eee;
  }
  .tools {
    position: absolute;
    left: 0;
    right: 0;
    display: flex;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
  }
  .tools.top {
    top: 0;
    justify-content: flex-end;
  }
  .tools.bottom {
    bottom: 0;
    flex-direction: column;
    align-items: center;
    padding-bottom: 1.25rem;
  }
  .round {
    width: 44px;
    height: 44px;
    border-radius: 50%;
    border: none;
    background: rgba(0, 0, 0, 0.55);
    color: #fff;
    font-size: 1.2rem;
  }
  .round.on {
    background: var(--bs-primary);
    color: #000;
  }
  .notice {
    background: rgba(0, 0, 0, 0.75);
    color: #fff;
    border-radius: 0.75rem;
    padding: 0.5rem 0.9rem;
    font-size: 0.85rem;
    max-width: 90%;
    text-align: center;
  }
  .sheet {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    max-height: 82%;
    display: flex;
    flex-direction: column;
    background: var(--bs-body-bg);
    border-radius: 1.25rem 1.25rem 0 0;
    box-shadow: 0 -6px 24px rgba(0, 0, 0, 0.45);
    max-width: 720px;
    margin: 0 auto;
  }
  .sheet-handle {
    width: 40px;
    height: 5px;
    border-radius: 3px;
    background: var(--bs-secondary-bg);
    margin: 0.5rem auto 0.25rem;
  }
  .sheet-body {
    overflow-y: auto;
    padding: 0.5rem 1rem;
  }
  .sheet-actions {
    padding: 0.75rem 1rem 1rem;
    border-top: 1px solid var(--bs-border-color);
  }
</style>
