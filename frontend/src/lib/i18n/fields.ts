// Labels for decoded payloads (qrformats.Parse: kind + fields), German and
// English. Unknown keys fall back to the key itself.

import type { Lang } from './errors';

export const KIND_LABELS: Record<string, string> = {
  epc: 'GiroCode (SEPA-Überweisung)',
  wifi: 'WLAN-Zugang',
  vcard: 'Visitenkarte',
  event: 'Termin',
  geo: 'Ort',
  tel: 'Telefonnummer',
  sms: 'SMS',
  email: 'E-Mail',
  crypto: 'Krypto-Zahlung',
  url: 'Link',
  pharma: 'Arzneimittel (securPharm)'
};

export const KIND_LABELS_EN: Record<string, string> = {
  epc: 'GiroCode (SEPA transfer)',
  wifi: 'Wi-Fi access',
  vcard: 'Contact card',
  event: 'Event',
  geo: 'Location',
  tel: 'Phone number',
  sms: 'SMS',
  email: 'Email',
  crypto: 'Crypto payment',
  url: 'Link',
  pharma: 'Medicine (securPharm)'
};

export function kindLabel(kind: string, lang: Lang = 'de'): string {
  return (lang === 'en' ? KIND_LABELS_EN : KIND_LABELS)[kind] ?? kind;
}

const FIELD_LABELS: Record<string, string> = {
  name: 'Empfänger',
  iban: 'IBAN',
  iban_valid: 'IBAN-Prüfsumme',
  bic: 'BIC',
  amount: 'Betrag',
  currency: 'Währung',
  reference: 'Verwendungszweck',
  purpose: 'Zweck',
  note: 'Hinweis',
  ssid: 'Netzwerk (SSID)',
  password: 'Passwort',
  encryption: 'Verschlüsselung',
  hidden: 'Verstecktes Netz',
  first_name: 'Vorname',
  last_name: 'Nachname',
  full_name: 'Name',
  org: 'Organisation',
  title: 'Titel / Funktion',
  phone: 'Telefon',
  email: 'E-Mail',
  url: 'Webseite',
  street: 'Straße',
  city: 'Ort',
  region: 'Region',
  zip: 'PLZ',
  country: 'Land',
  summary: 'Titel',
  start: 'Beginn',
  end: 'Ende',
  timezone: 'Zeitzone',
  all_day: 'Ganztägig',
  location: 'Ort',
  description: 'Beschreibung',
  latitude: 'Breitengrad',
  longitude: 'Längengrad',
  query: 'Suchbegriff',
  message: 'Nachricht',
  to: 'Empfänger',
  subject: 'Betreff',
  body: 'Text',
  coin: 'Währung',
  address: 'Adresse',
  label: 'Bezeichnung',
  format: 'Datenformat',
  pzn: 'PZN',
  pzn_valid: 'PZN-Prüfziffer',
  ppn: 'PPN',
  ppn_valid: 'PPN-Prüfsumme',
  gtin: 'GTIN / NTIN',
  gtin_valid: 'GTIN-Prüfziffer',
  batch: 'Charge',
  expiry: 'Verwendbar bis',
  serial: 'Seriennummer',
  production_date: 'Herstelldatum',
  nhrn: 'Nationale Nummer (NHRN)',
  unparsed: 'Nicht erkannter Rest'
};

// Order in which fields are shown; the rest follows alphabetically.
const ORDER = ['pzn', 'pzn_valid', 'expiry', 'batch', 'serial', 'name', 'full_name', 'first_name', 'last_name', 'ssid', 'summary', 'to', 'phone', 'iban', 'iban_valid', 'bic', 'amount', 'currency', 'reference', 'start', 'end', 'all_day'];

const FIELD_LABELS_EN: Record<string, string> = {
  name: 'Recipient',
  iban: 'IBAN',
  iban_valid: 'IBAN checksum',
  bic: 'BIC',
  amount: 'Amount',
  currency: 'Currency',
  reference: 'Reference',
  purpose: 'Purpose',
  note: 'Note',
  ssid: 'Network (SSID)',
  password: 'Password',
  encryption: 'Encryption',
  hidden: 'Hidden network',
  first_name: 'First name',
  last_name: 'Last name',
  full_name: 'Name',
  org: 'Organization',
  title: 'Title / role',
  phone: 'Phone',
  email: 'Email',
  url: 'Website',
  street: 'Street',
  city: 'City',
  region: 'Region',
  zip: 'Postcode',
  country: 'Country',
  summary: 'Title',
  start: 'Start',
  end: 'End',
  timezone: 'Time zone',
  all_day: 'All day',
  location: 'Location',
  description: 'Description',
  latitude: 'Latitude',
  longitude: 'Longitude',
  query: 'Search term',
  message: 'Message',
  to: 'To',
  subject: 'Subject',
  body: 'Text',
  coin: 'Currency',
  address: 'Address',
  label: 'Label',
  format: 'Data format',
  pzn: 'PZN',
  pzn_valid: 'PZN check digit',
  ppn: 'PPN',
  ppn_valid: 'PPN checksum',
  gtin: 'GTIN / NTIN',
  gtin_valid: 'GTIN check digit',
  batch: 'Batch',
  expiry: 'Use by',
  serial: 'Serial number',
  production_date: 'Production date',
  nhrn: 'National number (NHRN)',
  unparsed: 'Unrecognized rest'
};

export function fieldLabel(key: string, lang: Lang = 'de'): string {
  return (lang === 'en' ? FIELD_LABELS_EN : FIELD_LABELS)[key] ?? key;
}

export function sortedFieldKeys(fields: Record<string, string | undefined>): string[] {
  const rank = (k: string) => (ORDER.includes(k) ? ORDER.indexOf(k) : ORDER.length);
  return Object.keys(fields).sort((a, b) => rank(a) - rank(b) || a.localeCompare(b));
}

// iCalendar date (YYYYMMDD) or date-time (YYYYMMDDTHHMMSS[Z]) → 05.09.2026 18:00
// (en: 2026-09-05 18:00)
function formatICalDate(v: string, lang: Lang): string {
  const m = /^(\d{4})(\d{2})(\d{2})(?:T(\d{2})(\d{2})\d{2}(Z?))?$/.exec(v);
  if (!m) return v;
  const date = lang === 'en' ? `${m[1]}-${m[2]}-${m[3]}` : `${m[3]}.${m[2]}.${m[1]}`;
  return m[4] ? `${date} ${m[4]}:${m[5]}${m[6] ? ' UTC' : ''}` : date;
}

// GS1/IFA date YYMMDD; DD 00 means the end of the month.
function parsePackDate(v: string) {
  const m = /^(\d{2})(\d{2})(\d{2})$/.exec(v);
  if (!m) return null;
  const year = 2000 + Number(m[1]);
  const month = Number(m[2]);
  const day = Number(m[3]);
  const last = day === 0 ? new Date(year, month, 0) : new Date(year, month - 1, day);
  last.setHours(23, 59, 59);
  return { year, mm: m[2], dd: m[3], endOfMonth: day === 0, last };
}

// isExpired reports whether a pack date (the expiry field) has passed.
export function isExpired(value: string | undefined): boolean {
  const d = value ? parsePackDate(value) : null;
  return d !== null && d.last < new Date();
}

function formatPackDate(v: string, markExpired: boolean, lang: Lang): string {
  const d = parsePackDate(v);
  if (!d) return v;
  const text = d.endOfMonth
    ? `${d.mm}/${d.year}`
    : lang === 'en'
      ? `${d.year}-${d.mm}-${d.dd}`
      : `${d.dd}.${d.mm}.${d.year}`;
  if (!markExpired || !isExpired(v)) return text;
  return `${text} – ${lang === 'en' ? 'expired' : 'abgelaufen'}`;
}

export function isCheckKey(key: string): boolean {
  return key.endsWith('_valid');
}

export function formatFieldValue(key: string, value: string | undefined, lang: Lang = 'de'): string {
  const en = lang === 'en';
  if (value === undefined) return '';
  if (isCheckKey(key)) {
    if (value === 'true') return en ? '✓ valid' : '✓ gültig';
    return en ? '✗ check digit wrong – please check!' : '✗ Prüfziffer falsch – bitte prüfen!';
  }
  if (key === 'expiry') return formatPackDate(value, true, lang);
  if (key === 'production_date') return formatPackDate(value, false, lang);
  if (key === 'format') return value === 'ifa' ? 'IFA (PPN)' : value === 'gs1' ? 'GS1' : value;
  if (value === 'true') return en ? 'yes' : 'ja';
  if (value === 'false') return en ? 'no' : 'nein';
  if (key === 'start' || key === 'end') return formatICalDate(value, lang);
  if (key === 'amount') return en ? value : value.replace('.', ',');
  return value;
}
