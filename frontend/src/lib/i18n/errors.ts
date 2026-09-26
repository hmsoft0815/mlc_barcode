// Translations of the engine's input errors (barcodes.InputError).
// Go sends a stable code plus parameters; CLI and MCP print the English
// text, the GUI words it here. Unknown codes and languages without a
// catalog fall back to the English text, so a message is never lost.

import { BARCODE_TYPES } from '../types';

export type Lang = 'de' | 'en';
export const UI_LANG: Lang = 'de';

type Params = Record<string, string | undefined>;
type Catalog = Record<string, (p: Params) => string>;

export interface ErrorSource {
  error?: string;
  errorCode?: string;
  errorParams?: Record<string, string | undefined> | null;
}

const typeName = (id?: string) => BARCODE_TYPES.find((t) => t.id === id)?.name ?? (id ?? '').toUpperCase();
const quoted = (c?: string) => (c && c.trim() ? `„${c}“` : 'ein Steuerzeichen');

const capacityHintDe: Record<string, string> = {
  datamatrix: 'kürzen oder QR nehmen, das fasst mehr.',
  pdf417: 'kürzen oder QR nehmen.',
  qr: 'kürzen, auf mehrere Codes aufteilen oder Aztec nehmen, das fasst mehr Text.'
};

const de: Catalog = {
  empty: () => 'Bitte einen Inhalt eingeben.',
  retail_non_digit: (p) => `${typeName(p.type)} erlaubt nur Ziffern.`,
  retail_length: (p) => `${typeName(p.type)} braucht ${p.max} Ziffern (oder ${p.min} ohne Prüfziffer).`,
  retail_checksum: (p) => `Prüfziffer ${p.given} ist falsch, richtig wäre ${p.expected} (${p.code}).`,
  code39_charset: (p) =>
    `Code 39 kann ${quoted(p.char)} an Position ${p.pos} nicht darstellen; erlaubt sind A–Z, 0–9, Leerzeichen und - . $ / + %. ` +
    (p.fix === 'upper'
      ? 'In Großbuchstaben umwandeln – oder Code 128 nehmen, um Kleinbuchstaben zu behalten.'
      : 'Für alle ASCII-Zeichen Code 128 nehmen, für beliebigen Text QR.'),
  ascii_only: (p) =>
    `${typeName(p.type)} kann nur ASCII-Zeichen; ${quoted(p.char)} an Position ${p.pos} gehört nicht dazu. ` +
    'Für Umlaute und anderen Text QR, DataMatrix oder Aztec nehmen.',
  itf_digits: (p) =>
    `ITF erlaubt nur Ziffern; ${quoted(p.char)} an Position ${p.pos} ist keine. Für Buchstaben Code 128 nehmen.`,
  itf_even: (p) =>
    `ITF braucht eine gerade Anzahl Ziffern, eingegeben sind ${p.count}. Eine führende 0 ergänzen: ${p.suggestion}`,
  capacity: (p) =>
    `Zu viele Daten für ${typeName(p.type)}: ${p.chars} Zeichen (${p.bytes} Bytes), höchstens ${p.max} Zeichen dieses Inhalts passen hinein – ` +
    (capacityHintDe[p.type ?? ''] ?? 'kürzen oder auf mehrere Codes aufteilen.'),
  unsupported_type: (p) => `Unbekannter Barcode-Typ „${p.type}“.`,
  encoder: (p) => `${typeName(p.type)} kann diese Eingabe nicht codieren (${p.detail}).`,
  aztec_charset: (p) =>
    `Aztec kann nur Zeichen aus ISO-8859-1 (Latin-1, mit äöüß) darstellen; ${quoted(p.char)} an Position ${p.pos} gehört nicht dazu. Dafür QR oder DataMatrix nehmen.`,
  image_format: () =>
    'Dieses Bildformat wird nicht unterstützt. Bitte PNG, JPEG, GIF oder WebP verwenden (iPhone-Fotos im HEIC-Format vorher als JPEG exportieren).',
  image_too_large: (p) =>
    p.max_mp
      ? `Das Bild ist ${p.width}×${p.height} Pixel groß, höchstens ${p.max_mp} Megapixel sind erlaubt. Bitte verkleinern.`
      : `Das Bild ist ${p.size_mb} MB groß, höchstens ${p.max_mb} MB sind erlaubt.`,
  nothing_found: () =>
    'Im Bild wurde kein Barcode gefunden. Ist der Code scharf, vollständig sichtbar und nicht zu klein?'
};

const catalogs: Partial<Record<Lang, Catalog>> = { de };

export function formatError(src: ErrorSource, lang: Lang = UI_LANG): string {
  const fn = src.errorCode ? catalogs[lang]?.[src.errorCode] : undefined;
  if (fn) return fn(src.errorParams ?? {});
  return src.error || (lang === 'de' ? 'Unbekannter Fehler' : 'Unknown error');
}
