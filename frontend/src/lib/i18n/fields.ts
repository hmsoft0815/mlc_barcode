// German labels for decoded payloads (qrformats.Parse: kind + fields).
// Unknown keys fall back to the key itself.

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
  url: 'Link'
};

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
  label: 'Bezeichnung'
};

// Order in which fields are shown; the rest follows alphabetically.
const ORDER = ['name', 'full_name', 'first_name', 'last_name', 'ssid', 'summary', 'to', 'phone', 'iban', 'iban_valid', 'bic', 'amount', 'currency', 'reference', 'start', 'end', 'all_day'];

export function fieldLabel(key: string): string {
  return FIELD_LABELS[key] ?? key;
}

export function sortedFieldKeys(fields: Record<string, string | undefined>): string[] {
  const rank = (k: string) => (ORDER.includes(k) ? ORDER.indexOf(k) : ORDER.length);
  return Object.keys(fields).sort((a, b) => rank(a) - rank(b) || a.localeCompare(b));
}

// iCalendar date (YYYYMMDD) or date-time (YYYYMMDDTHHMMSS[Z]) → 05.09.2026 18:00
function formatICalDate(v: string): string {
  const m = /^(\d{4})(\d{2})(\d{2})(?:T(\d{2})(\d{2})\d{2}(Z?))?$/.exec(v);
  if (!m) return v;
  const date = `${m[3]}.${m[2]}.${m[1]}`;
  return m[4] ? `${date} ${m[4]}:${m[5]}${m[6] ? ' UTC' : ''}` : date;
}

export function formatFieldValue(key: string, value: string | undefined): string {
  if (value === undefined) return '';
  if (value === 'true') return 'ja';
  if (value === 'false') return 'nein';
  if (key === 'start' || key === 'end') return formatICalDate(value);
  if (key === 'amount') return value.replace('.', ',');
  return value;
}
