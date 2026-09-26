// UI texts of the mobile UI and the shared camera, German and English. The
// English catalog is typed against the German one: a missing translation
// fails the type check (svelte-check), not the user.

const de = {
  // tabs
  tabScan: 'Scannen',
  tabHistory: 'Verlauf',
  tabCreate: 'Erzeugen',
  tabAbout: 'Über',

  // camera
  camDenied: 'Der Zugriff auf die Kamera wurde nicht erlaubt.',
  camNone: 'Keine Kamera gefunden.',
  camBusy: 'Die Kamera wird gerade von einem anderen Programm verwendet.',
  camUnavailable: 'Kamera nicht verfügbar ({name}).',
  camNoAccess: 'Diese Ansicht hat keinen Zugriff auf die Kamera. Ein Foto aus der Galerie geht trotzdem.',

  // scan
  retry: 'Erneut versuchen',
  holdInFrame: 'Code in den Rahmen halten',
  torch: 'Licht',
  switchCamera: 'Kamera wechseln',
  pickPhoto: 'Foto wählen',
  reading: 'Lese …',
  foundOne: 'Code gefunden',
  foundMany: '{n} Codes gefunden',
  scanAgain: 'Weiter scannen',

  // result
  open: 'Öffnen',
  copy: 'Kopieren',
  share: 'Teilen',
  copied: 'Kopiert',
  copiedNoShare: 'Kopiert (Teilen nicht verfügbar)',
  rawContent: 'Rohinhalt',

  // history
  historyTitle: 'Verlauf',
  clear: 'Leeren',
  clearAll: 'Alle löschen',
  cancel: 'Abbrechen',
  historyEmpty: 'Noch nichts gescannt. Gescannte Codes erscheinen hier und bleiben auf diesem Gerät.',
  removeEntry: 'Aus dem Verlauf entfernen',

  // create
  createTitle: 'Erzeugen',
  kindLink: 'Link / Text',
  kindWifi: 'WLAN',
  kindContact: 'Kontakt',
  kindGiro: 'GiroCode',
  linkPlaceholder: 'https://… oder beliebiger Text',
  ssid: 'Netzwerkname (SSID)',
  password: 'Passwort',
  encWpa: 'WPA/WPA2/WPA3',
  encWep: 'WEP',
  encOpen: 'Offen (ohne Passwort)',
  hiddenNet: 'versteckt',
  firstName: 'Vorname',
  lastName: 'Nachname',
  phone: 'Telefon',
  email: 'E-Mail',
  recipient: 'Empfänger',
  iban: 'IBAN',
  amountOptional: 'Betrag (€, optional)',
  reference: 'Verwendungszweck',
  createQr: 'QR-Code erzeugen',
  enlarge: 'Groß anzeigen',
  tapToEnlarge: 'Antippen zum Vergrößern, z. B. zum Abscannen',
  tapToClose: 'Antippen zum Schließen',
  copyContent: 'Inhalt kopieren',
  contentCopied: 'Inhalt kopiert',
  contentCopiedNoShare: 'Inhalt kopiert (Teilen nicht verfügbar)',

  // about
  aboutText:
    'Barcodes und QR-Codes scannen und erzeugen — Arzneimittel mit Verfallsdatum, GiroCodes, WLAN, Kontakte, Tickets und Pakete. Alles wird auf diesem Gerät gelesen, nichts wird hochgeladen.',
  productPage: 'Produktseite auf mlcgo.eu',
  sourceCode: 'Quellcode auf GitHub',
  desktopUI: 'Desktop-Oberfläche verwenden',
  language: 'Sprache',
  license: 'Lizenz',
  licenseSummary:
    '© 2026 Michael Lechner. Open Source unter der MIT-Lizenz mit Namensnennung: frei nutzbar, auch kommerziell, solange MLC Barcode als Quelle genannt wird.',
  licenseText: 'Lizenztext',
  thirdParty: 'Verwendete Open-Source-Komponenten',
};

export type UIKey = keyof typeof de;

const en: Record<UIKey, string> = {
  tabScan: 'Scan',
  tabHistory: 'History',
  tabCreate: 'Create',
  tabAbout: 'About',

  camDenied: 'Access to the camera was not allowed.',
  camNone: 'No camera found.',
  camBusy: 'The camera is in use by another app.',
  camUnavailable: 'Camera not available ({name}).',
  camNoAccess: 'This view has no access to the camera. A photo from the gallery still works.',

  retry: 'Try again',
  holdInFrame: 'Hold the code inside the frame',
  torch: 'Light',
  switchCamera: 'Switch camera',
  pickPhoto: 'Choose photo',
  reading: 'Reading …',
  foundOne: 'Code found',
  foundMany: '{n} codes found',
  scanAgain: 'Scan again',

  open: 'Open',
  copy: 'Copy',
  share: 'Share',
  copied: 'Copied',
  copiedNoShare: 'Copied (sharing not available)',
  rawContent: 'Raw content',

  historyTitle: 'History',
  clear: 'Clear',
  clearAll: 'Delete all',
  cancel: 'Cancel',
  historyEmpty: 'Nothing scanned yet. Scanned codes appear here and stay on this device.',
  removeEntry: 'Remove from history',

  createTitle: 'Create',
  kindLink: 'Link / text',
  kindWifi: 'Wi-Fi',
  kindContact: 'Contact',
  kindGiro: 'GiroCode',
  linkPlaceholder: 'https://… or any text',
  ssid: 'Network name (SSID)',
  password: 'Password',
  encWpa: 'WPA/WPA2/WPA3',
  encWep: 'WEP',
  encOpen: 'Open (no password)',
  hiddenNet: 'hidden',
  firstName: 'First name',
  lastName: 'Last name',
  phone: 'Phone',
  email: 'Email',
  recipient: 'Recipient',
  iban: 'IBAN',
  amountOptional: 'Amount (€, optional)',
  reference: 'Reference',
  createQr: 'Create QR code',
  enlarge: 'Show large',
  tapToEnlarge: 'Tap to enlarge, e.g. for someone to scan',
  tapToClose: 'Tap to close',
  copyContent: 'Copy content',
  contentCopied: 'Content copied',
  contentCopiedNoShare: 'Content copied (sharing not available)',

  aboutText:
    'Scan and create barcodes and QR codes — medicines with their expiry date, GiroCodes, Wi-Fi, contacts, tickets and parcels. Everything is read on this device; nothing is uploaded.',
  productPage: 'Product page on mlcgo.eu',
  sourceCode: 'Source code on GitHub',
  desktopUI: 'Use the desktop interface',
  language: 'Language',
  license: 'License',
  licenseSummary:
    '© 2026 Michael Lechner. Open source under the MIT license with attribution: free to use, also commercially, as long as MLC Barcode is credited.',
  licenseText: 'License text',
  thirdParty: 'Open-source components used',
};

export const UI_TEXT = { de, en } as const;
