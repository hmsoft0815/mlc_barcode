// Short descriptions of the symbologies (BARCODE_TYPES), shown in the
// generator's type list and the help table.
import { translator } from '../lang';

const de = {
  desc_qr: '2D Matrixcode für Text, URLs, WLAN, vCard, Events',
  desc_datamatrix: 'Kompakter 2D-Code für Industrie und Bauteile',
  desc_aztec: 'Kompakter 2D-Code für Bahn- und Flugtickets – für Fachscanner, Handy-Kameras lesen ihn meist nicht',
  desc_pdf417: 'Gestapelter 2D-Code für Bordkarten und Versandetiketten – für Fachscanner, Handy-Kameras lesen ihn meist nicht',
  desc_code128: 'Universeller 1D-Barcode für alle ASCII-Zeichen',
  desc_code39: 'Alphanumerischer Barcode (Großbuchstaben, Ziffern)',
  desc_ean13: '13-stelliger Einzelhandels-Barcode (12 Ziffern + Prüfziffer)',
  desc_ean8: '8-stelliger kompakter Handels-Barcode',
  desc_upca: '12-stelliger US-Einzelhandels-Barcode',
  desc_itf: 'Kompakter numerischer Barcode (nur Ziffern, gerade Anzahl)',
};

const en: Record<keyof typeof de, string> = {
  desc_qr: '2D matrix code for text, URLs, Wi-Fi, vCard, events',
  desc_datamatrix: 'Compact 2D code for industry and components',
  desc_aztec: 'Compact 2D code for train and plane tickets – for dedicated scanners, phone cameras mostly cannot read it',
  desc_pdf417: 'Stacked 2D code for boarding passes and shipping labels – for dedicated scanners, phone cameras mostly cannot read it',
  desc_code128: 'Universal 1D barcode for all ASCII characters',
  desc_code39: 'Alphanumeric barcode (capital letters, digits)',
  desc_ean13: '13-digit retail barcode (12 digits + check digit)',
  desc_ean8: '8-digit compact retail barcode',
  desc_upca: '12-digit US retail barcode',
  desc_itf: 'Compact numeric barcode (digits only, even count)',
};

export type TypeDescKey = keyof typeof de;
export const typeText = translator(de, en);

export const typeDescKey = (id: string) => `desc_${id}` as TypeDescKey;
