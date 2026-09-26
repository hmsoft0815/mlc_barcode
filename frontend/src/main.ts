import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import '@fontsource-variable/inter';
import './theme.css';
import { mount } from 'svelte';
import { Browser } from '@wailsio/runtime';
import App from './App.svelte';

// External links belong in the system browser. Inside Wails a target=_blank
// link would open a second app window instead.
document.addEventListener(
  'click',
  (e) => {
    const link = (e.target as Element | null)?.closest?.('a[href]') as HTMLAnchorElement | null;
    if (!link || !/^(https?:|mailto:)/i.test(link.href)) return;
    e.preventDefault();
    Browser.OpenURL(link.href).catch((err) => console.warn('Could not open link:', err));
  },
  true
);

mount(App, { target: document.getElementById('app')! });
