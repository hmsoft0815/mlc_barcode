// Texts of the PZN lookup (PznLookup.svelte).
import { translator } from '../lang';

const de = {
  button: 'Beim BfArM nachschlagen',
  title: 'PZN beim BfArM nachschlagen',
  intro: 'Die Arzneimittel-Datenbank des BfArM hat keinen direkten Such-Link. So geht es:',
  step1: 'Unten „Suche öffnen“ tippen, dort „Akzeptieren und weiter zur Recherche“.',
  step2: 'Rechts neben „Suche nach“ im Feld „in“ den Eintrag „Pharmazentralnummer“ wählen.',
  step3: 'Ins Feld „Suche nach“ die PZN einfügen (lange tippen → Einsetzen) und „Suche starten“.',
  pznLabel: 'Ihre PZN (schon kopiert)',
  copyAgain: 'Nochmal kopieren',
  copied: 'Kopiert',
  open: 'Suche öffnen',
  close: 'Schließen',
  note: 'Nur zugelassene Arzneimittel sind dort verzeichnet — Kosmetik, Nahrungsergänzung und Medizinprodukte nicht.',
};

const en: Record<keyof typeof de, string> = {
  button: 'Look up at BfArM',
  title: 'Look up the PZN at BfArM',
  intro: 'The BfArM medicine database has no direct search link. This is how it works:',
  step1: 'Tap "Open search" below, then "Akzeptieren und weiter zur Recherche" (accept and continue).',
  step2: 'Next to "Suche nach", in the field "in", choose "Pharmazentralnummer".',
  step3: 'Paste the PZN into "Suche nach" (long-press → Paste) and tap "Suche starten".',
  pznLabel: 'Your PZN (already copied)',
  copyAgain: 'Copy again',
  copied: 'Copied',
  open: 'Open search',
  close: 'Close',
  note: 'Only licensed medicines are listed there — not cosmetics, food supplements or medical devices.',
};

export const t = translator(de, en);
