// Texts of the desktop shell: navbar (Navbar.svelte) and footer (App.svelte).
import { translator } from '../lang';

const de = {
  logoAlt: 'MLC Barcode Logo',
  tabSingle: 'Einzel-Generator',
  tabBatch: 'Batch-Generator',
  tabPrint: 'Etiketten-Druck',
  tabCheck: 'Prüfen',
  tabHelp: 'Hilfe',
  toLight: 'Zu hellem Modus wechseln',
  toDark: 'Zu dunklem Modus wechseln',
  toggleTheme: 'Design umschalten',
  language: 'Sprache',
  myDetails: 'Meine Angaben',
  footerLicense: 'MIT with Attribution',
};

const en: Record<keyof typeof de, string> = {
  logoAlt: 'MLC Barcode logo',
  tabSingle: 'Single generator',
  tabBatch: 'Batch generator',
  tabPrint: 'Label printing',
  tabCheck: 'Check',
  tabHelp: 'Help',
  toLight: 'Switch to light mode',
  toDark: 'Switch to dark mode',
  toggleTheme: 'Toggle theme',
  language: 'Language',
  myDetails: 'My details',
  footerLicense: 'MIT with Attribution',
};

export const t = translator(de, en);
