// Texts of the desktop "Prüfen" tab (Checker.svelte).
import { translator } from '../lang';

const de = {
  title: 'Barcode prüfen',
  dropLabel: 'Bild ablegen',
  holdInImage: 'Code ins Bild halten …',
  dragFrame: 'Rahmen um den Code ziehen …',
  dropHint: 'Bild hierher ziehen, auswählen oder mit {key}+V einfügen',
  pasteKey: 'Strg',
  formats: 'PNG, JPEG, GIF oder WebP · Fotos, Scans, Screenshots',
  pickImage: 'Bild auswählen …',
  switchCamera: 'Kamera wechseln',
  stopCamera: 'Kamera beenden',
  scanWithCamera: 'Mit Kamera scannen',
  cancelSelection: 'Auswahl abbrechen',
  selectRegion: 'Ausschnitt wählen',
  camera: 'Kamera',
  cameraGuide: 'Kamera (Zielrahmen)',
  result: 'Ergebnis',
  reading: 'Lese …',
  foundOne: '1 Code gefunden',
  foundMany: '{n} Codes gefunden',
  copied: 'Kopiert',
  copyContent: 'Inhalt kopieren',
  intro:
    'Liest QR, DataMatrix, Aztec, PDF417, EAN-13/8, UPC-A, Code 128, Code 39 und ITF – auch mehrere Codes pro Bild. Bekannte Inhalte wie GiroCode, Visitenkarte, WLAN oder Termin werden in ihre Felder zerlegt; beim GiroCode wird die IBAN-Prüfsumme kontrolliert, bei Arzneimittel-Codes (securPharm) PZN, Charge, Verfall und Seriennummer angezeigt.',
};

const en: Record<keyof typeof de, string> = {
  title: 'Check a barcode',
  dropLabel: 'Drop image',
  holdInImage: 'Hold the code in view …',
  dragFrame: 'Drag a frame around the code …',
  dropHint: 'Drag an image here, choose one or paste it with {key}+V',
  pasteKey: 'Ctrl',
  formats: 'PNG, JPEG, GIF or WebP · photos, scans, screenshots',
  pickImage: 'Choose image …',
  switchCamera: 'Switch camera',
  stopCamera: 'Stop camera',
  scanWithCamera: 'Scan with camera',
  cancelSelection: 'Cancel selection',
  selectRegion: 'Select region',
  camera: 'Camera',
  cameraGuide: 'Camera (guide frame)',
  result: 'Result',
  reading: 'Reading …',
  foundOne: '1 code found',
  foundMany: '{n} codes found',
  copied: 'Copied',
  copyContent: 'Copy content',
  intro:
    'Reads QR, DataMatrix, Aztec, PDF417, EAN-13/8, UPC-A, Code 128, Code 39 and ITF – several codes per image, too. Known contents such as GiroCode, contact card, Wi-Fi or event are split into their fields; for GiroCodes the IBAN checksum is verified, for medicine codes (securPharm) PZN, batch, expiry and serial number are shown.',
};

export const t = translator(de, en);
