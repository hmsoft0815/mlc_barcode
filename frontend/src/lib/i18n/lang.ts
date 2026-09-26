// The UI language: German or English. It follows the device unless the
// user picked one (stored per device). The desktop UI is German for now;
// the mobile UI and the shared parts follow this store.
import { derived, writable } from 'svelte/store';
import type { Lang } from './errors';
import { UI_TEXT, type UIKey } from './ui';

const KEY = 'mlc_lang';

function initial(): Lang {
  const asked = typeof location !== 'undefined' ? new URLSearchParams(location.search).get('lang') : null;
  if (asked === 'de' || asked === 'en') return asked; // ?lang=de, e.g. for screenshots
  try {
    const saved = localStorage.getItem(KEY);
    if (saved === 'de' || saved === 'en') return saved;
  } catch {
    // storage blocked: follow the device
  }
  return typeof navigator !== 'undefined' && navigator.language?.toLowerCase().startsWith('de') ? 'de' : 'en';
}

export const lang = writable<Lang>(initial());

export function setLang(l: Lang) {
  lang.set(l);
  try {
    localStorage.setItem(KEY, l);
  } catch {
    // not remembered, still switched
  }
  document.documentElement.lang = l;
}

type Params = Record<string, string | number>;

// translator makes a $t store for one dictionary. The English catalog is
// typed against the German one, so a missing translation fails the type
// check. Each UI part (mobile, each desktop tab) has its own dictionary in
// i18n/text/, so they can grow independently.
export function translator<K extends string>(de: Record<K, string>, en: Record<K, string>) {
  const catalogs = { de, en };
  return derived(lang, (l) => (key: K, params: Params = {}) =>
    catalogs[l][key].replace(/\{(\w+)\}/g, (_, p) => String(params[p] ?? ''))
  );
}

// $t('key', { n: 3 }) — {name} placeholders are filled from params.
export const t = translator<UIKey>(UI_TEXT.de, UI_TEXT.en);
