import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import '@fontsource-variable/inter';
import './theme.css';
import { mount } from 'svelte';
import App from './App.svelte';
import MobileApp from './lib/mobile/MobileApp.svelte';
import { openLink } from './lib/mobile/native';
import { Platform, SetLanguage } from '../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';
import { lang } from './lib/i18n/lang';

// The Go side words its native dialogs and a few messages itself.
lang.subscribe((l) => {
  document.documentElement.lang = l;
  SetLanguage(l).catch(() => {});
});

// External links belong in the system browser. Inside Wails a target=_blank
// link would open a second app window instead.
document.addEventListener(
  'click',
  (e) => {
    const link = (e.target as Element | null)?.closest?.('a[href]') as HTMLAnchorElement | null;
    if (!link || !/^(https?:|mailto:)/i.test(link.href)) return;
    e.preventDefault();
    openLink(link.href).catch((err) => console.warn('Could not open link:', err));
  },
  true
);

// Which UI: the mobile one on phones and tablets (iOS, Android), the
// desktop one elsewhere. On the desktop the mobile UI can be chosen by hand
// (stored per device) or with ?ui=mobile, e.g. to try it in a narrow window.
const UI_KEY = 'mlc_ui';

function storedUI(): string | null {
  try {
    return localStorage.getItem(UI_KEY);
  } catch {
    return null;
  }
}

function setUI(ui: 'mobile' | 'desktop') {
  try {
    localStorage.setItem(UI_KEY, ui);
  } catch {
    // not remembered; the reload below still switches
  }
  location.search = '';
}
(window as any).mlcSetUI = setUI;

async function start() {
  const os = await Platform().catch(() => '');
  const mobileOS = os === 'ios' || os === 'android';
  const asked = new URLSearchParams(location.search).get('ui') ?? storedUI();
  const target = document.getElementById('app')!;
  if (mobileOS || asked === 'mobile') {
    mount(MobileApp, { target, props: { desktop: !mobileOS, onDesktopUI: () => setUI('desktop') } });
  } else {
    mount(App, { target });
  }
}

start();
