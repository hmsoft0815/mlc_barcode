<script lang="ts">
  // "My details": name, contact, optional place and time zone. Changes are
  // stored as you type (profile.ts); used by the desktop dialog and the
  // mobile "About" tab.
  import { profile, organizerName } from '../profile';
  import { deviceTimeZone } from '../datetime';
  import { t } from '../i18n/text/profile';

  export let compact = false; // mobile: full-width fields, no columns

  const deviceTZ = deviceTimeZone();
  let locating = false;
  let positionError = '';

  function useMyPosition() {
    if (!navigator.geolocation) {
      positionError = $t('positionFailed', { reason: 'geolocation' });
      return;
    }
    locating = true;
    positionError = '';
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        locating = false;
        $profile.latitude = pos.coords.latitude.toFixed(6);
        $profile.longitude = pos.coords.longitude.toFixed(6);
      },
      (err) => {
        locating = false;
        positionError = $t('positionFailed', { reason: err.message || String(err.code) });
      },
      { enableHighAccuracy: true, timeout: 15000 }
    );
  }

  function clearAll() {
    profile.set({ lastName: '', firstName: '', email: '', phone: '', location: '', latitude: '', longitude: '', timeZone: '' });
  }

  $: organizer = organizerName($profile);
</script>

<p class="small text-body-secondary">{$t('intro')}</p>

<div class="row g-2 mb-1">
  <div class={compact ? 'col-12' : 'col-md-6'}>
    <label class="form-label small mb-1" for="pf-last">{$t('lastName')}</label>
    <input id="pf-last" class="form-control form-control-sm" autocomplete="family-name" bind:value={$profile.lastName} />
  </div>
  <div class={compact ? 'col-12' : 'col-md-6'}>
    <label class="form-label small mb-1" for="pf-first">{$t('firstName')}</label>
    <input id="pf-first" class="form-control form-control-sm" autocomplete="given-name" bind:value={$profile.firstName} />
  </div>
  <div class={compact ? 'col-12' : 'col-md-6'}>
    <label class="form-label small mb-1" for="pf-email">{$t('email')}</label>
    <input id="pf-email" type="email" class="form-control form-control-sm" autocomplete="email" autocapitalize="off" bind:value={$profile.email} />
  </div>
  <div class={compact ? 'col-12' : 'col-md-6'}>
    <label class="form-label small mb-1" for="pf-phone">{$t('phone')}</label>
    <input id="pf-phone" type="tel" class="form-control form-control-sm" autocomplete="tel" bind:value={$profile.phone} />
  </div>
</div>
{#if organizer}
  <div class="small text-body-secondary mb-3"><i class="bi bi-person-check me-1"></i>{$t('organizerPreview', { name: organizer })}</div>
{/if}

<h6 class="small fw-semibold mt-3 mb-2">{$t('place')}</h6>
<div class="row g-2">
  <div class="col-12">
    <label class="form-label small mb-1" for="pf-loc">{$t('location')}</label>
    <input id="pf-loc" class="form-control form-control-sm" placeholder={$t('locationPh')} bind:value={$profile.location} />
  </div>
  <div class="col-6">
    <label class="form-label small mb-1" for="pf-lat">{$t('latitude')}</label>
    <input id="pf-lat" class="form-control form-control-sm" inputmode="decimal" placeholder="52.520000" bind:value={$profile.latitude} />
  </div>
  <div class="col-6">
    <label class="form-label small mb-1" for="pf-lon">{$t('longitude')}</label>
    <input id="pf-lon" class="form-control form-control-sm" inputmode="decimal" placeholder="13.405000" bind:value={$profile.longitude} />
  </div>
  <div class="col-12">
    <button type="button" class="btn btn-sm btn-outline-primary" on:click={useMyPosition} disabled={locating}>
      <i class="bi bi-crosshair me-1"></i>{locating ? $t('locating') : $t('useMyPosition')}
    </button>
    {#if positionError}<div class="small text-warning mt-1">{positionError}</div>{/if}
  </div>
  <div class="col-12">
    <label class="form-label small mb-1" for="pf-tz">{$t('timeZone')}</label>
    <input id="pf-tz" class="form-control form-control-sm" placeholder={deviceTZ} autocapitalize="off" bind:value={$profile.timeZone} />
    {#if deviceTZ}<div class="form-text small">{$t('timeZoneDevice', { tz: deviceTZ })}</div>{/if}
  </div>
</div>

<button type="button" class="btn btn-sm btn-link text-danger px-0 mt-3" on:click={clearAll}>
  <i class="bi bi-trash me-1"></i>{$t('clear')}
</button>
