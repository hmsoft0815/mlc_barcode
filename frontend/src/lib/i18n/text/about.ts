// Texts of the desktop "Hilfe" tab (AboutView.svelte). Licence text and
// third-party notices are imported files and stay as they are.
import { translator } from '../lang';

const de = {
  logoAlt: 'MLC Barcode Logo',
  lead: 'Plattformübergreifendes Desktop-Werkzeug zur Generierung, Stapelverarbeitung und dem Druck von Barcodes & QR-Codes.',
  version: 'Version {v}',
  githubRepo: 'GitHub Repository',
  productPage: 'Produktseite (mlcgo.eu)',

  symbologiesTitle: 'Unterstützte Symbologien & Formate',
  colType: 'Typ',
  colCategory: 'Kategorie',
  colDescription: 'Beschreibung & Spezifikation',
  colSample: 'Beispiel',

  featuresTitle: 'Funktionen im Überblick',
  featPreviewLabel: 'Live-Vorschau:',
  featPreview: 'Direkte Vektordarstellung in SVG und hochauflösendem PNG.',
  featQrLabel: '9 QR-Sonderformate:',
  featQr: 'GiroCode (SEPA-Überweisung), Krypto-Wallets, Maps (Geo), Telefon, SMS, E-Mail, WLAN, vCard 3.0 und Kalender-Events.',
  featCheckLabel: 'Prüfen:',
  featCheckA: 'Barcodes aus Bildern lesen (Datei, Drag & Drop,',
  ctrlKey: 'Strg',
  featCheckB:
    ') – Format, Inhalt und zerlegte Felder, beim GiroCode mit IBAN-Prüfung. Jeder erzeugte Code wird automatisch gegengelesen („Lesbar geprüft“).',
  featCaptionLabel: 'Beschriftung:',
  featCaption: 'Freitext unter dem Barcode mit einstellbarer Schriftgröße, in SVG und PNG.',
  featBatchLabel: 'Batch-Generierung:',
  featBatchA: 'Import von',
  featBatchB: '(Format siehe unten) und Stapel-Export in einen Zielordner.',
  featSheetLabel: 'Etiketten-Druckbogen:',
  featSheet:
    'DIN-A4-Layouts nach gängigen Avery-Zweckform-Rastern; Etiketten aus dem Batch-Generator übernehmen oder direkt aus TXT/CSV mit eigenem Etikett-Text importieren.',

  cliTitle: 'CLI & MCP Server',
  cliIntro:
    'MLC Barcode ist auch als stand-alone Kommandozeilentool und MCP-Server für KI-Assistenten (Claude Desktop, Cursor, Gemini-CLI) verfügbar:',
  cliComment: '# Kommandozeile:',
  mcpComment1: '# MCP-Server: startet Ihr KI-Client selbst.',
  mcpComment2: '# Eintrag z. B. in claude_desktop_config.json (Windows-Setup):',
  mcpGuide: 'Anleitung: MCP-Server in Claude, Gemini & Cursor einbinden',

  importTitle: 'Datei-Import (Batch-Generator & Etiketten-Druck)',
  importTxt: 'Jede Zeile ist ein Barcode-Inhalt. Es wird nichts aufgeteilt — Kommas und Semikolons gehören zum Inhalt.',
  importCsvA: 'Spalte 1 = Inhalt (z. B. für den QR-Code), Spalte 2 =',
  importCsvOptional: 'optionaler',
  importCsvB: 'Etikett-Text. Der Batch-Generator nutzt nur Spalte 1.',
  importSepLabel: 'Trennzeichen:',
  importSepA: 'oder Tab werden erkannt (Excel speichert deutsch mit',
  importSepB:
    '). Ein Komma gilt nur als Trennzeichen, wenn jede Zeile gleich viele hat. Ohne erkennbares Trennzeichen ist jede Zeile ein Inhalt.',
  importQuoteA: 'Enthält der Inhalt selbst das Trennzeichen (z. B.',
  importQuoteB: '), das Feld in Anführungszeichen setzen.',
  importEmpty: 'Leere Zeilen werden übersprungen, eine Kopfzeile lässt sich im Etiketten-Druck abschalten.',
  importSample: 'ART-1001;Schraube M4\nART-1002\n"WIFI:T:WPA;S:Gast;P:geheim;;";Gäste-WLAN',

  mobileTitle: 'Mobile Oberfläche',
  mobileText:
    '— Scannen mit Kamera, Verlauf, Erzeugen im Handy-Stil. Auf iPhone, iPad und Android ist sie die Standardansicht; hier lässt sie sich ausprobieren. Zurück über „Über“.',
  mobileButton: 'Mobile Oberfläche',

  licenseTitle: 'Lizenz & Copyright',
  licenseA: 'Veröffentlicht als Open-Source-Software unter der',
  licenseName: 'MIT-Lizenz mit Namensnennung',
  licenseB:
    '(MIT with Attribution). Der Quellcode darf frei genutzt, verändert und weitergegeben werden — auch kommerziell. Produkte, die ihn verwenden, müssen den Autor „Michael Lechner“ sichtbar nennen (z. B. in der Dokumentation oder einem Info-Dialog). Eine kommerzielle Lizenz ohne Namensnennung ist auf Anfrage erhältlich.',
  licenseLine: 'Lizenz: MIT with Attribution (siehe unten)',
  showLicense: 'Lizenztext anzeigen',
  showThirdParty: 'Drittlizenzen anzeigen (ZXing, boombuler/barcode, gozxing)',
  libsTitle: 'Verwendete Open-Source Bibliotheken:',
  libWails: 'Desktop App Framework',
  libSvelte: 'Reactive UI Framework',
  libBoombuler: 'Barcode Engine',
  libGozxing: 'Barcodes lesen',
  libPdf417: 'nach Go portiert',
  libMcp: 'Model Context Protocol',
  libBootstrap: 'UI Styling & Icons',
  libInter: 'Schrift der Oberfläche',
  libGoImage: 'Beschriftung im PNG',
  libGoImageName: 'golang.org/x/image & Go-Schriften',
};

const en: Record<keyof typeof de, string> = {
  logoAlt: 'MLC Barcode logo',
  lead: 'Cross-platform desktop tool for generating, batch-processing and printing barcodes & QR codes.',
  version: 'Version {v}',
  githubRepo: 'GitHub repository',
  productPage: 'Product page (mlcgo.eu)',

  symbologiesTitle: 'Supported symbologies & formats',
  colType: 'Type',
  colCategory: 'Category',
  colDescription: 'Description & specification',
  colSample: 'Example',

  featuresTitle: 'Features at a glance',
  featPreviewLabel: 'Live preview:',
  featPreview: 'Direct vector rendering as SVG and high-resolution PNG.',
  featQrLabel: '9 special QR formats:',
  featQr: 'GiroCode (SEPA transfer), crypto wallets, maps (geo), phone, SMS, email, Wi-Fi, vCard 3.0 and calendar events.',
  featCheckLabel: 'Check:',
  featCheckA: 'Read barcodes from images (file, drag & drop,',
  ctrlKey: 'Ctrl',
  featCheckB:
    ') – format, content and parsed fields, with IBAN validation for GiroCodes. Every generated code is read back automatically ("Verified readable").',
  featCaptionLabel: 'Caption:',
  featCaption: 'Free text below the barcode with adjustable font size, in SVG and PNG.',
  featBatchLabel: 'Batch generation:',
  featBatchA: 'Import from',
  featBatchB: '(format see below) and batch export to a target folder.',
  featSheetLabel: 'Label sheets:',
  featSheet:
    'A4 layouts for common Avery Zweckform grids; take labels from the batch generator or import them directly from TXT/CSV with their own label text.',

  cliTitle: 'CLI & MCP server',
  cliIntro:
    'MLC Barcode is also available as a stand-alone command-line tool and as an MCP server for AI assistants (Claude Desktop, Cursor, Gemini CLI):',
  cliComment: '# Command line:',
  mcpComment1: '# MCP server: your AI client starts it itself.',
  mcpComment2: '# Entry e.g. in claude_desktop_config.json (Windows setup):',
  mcpGuide: 'Guide: add the MCP server to Claude, Gemini & Cursor',

  importTitle: 'File import (batch generator & label printing)',
  importTxt: 'Each line is one barcode content. Nothing is split — commas and semicolons are part of the content.',
  importCsvA: 'Column 1 = content (e.g. for the QR code), column 2 =',
  importCsvOptional: 'optional',
  importCsvB: 'label text. The batch generator only uses column 1.',
  importSepLabel: 'Separator:',
  importSepA: 'or tab are detected (German Excel saves with',
  importSepB:
    '). A comma only counts as separator if every line has the same number of them. Without a recognisable separator each line is one content.',
  importQuoteA: 'If the content itself contains the separator (e.g.',
  importQuoteB: '), put the field in quotes.',
  importEmpty: 'Empty lines are skipped; a header line can be turned off in label printing.',
  importSample: 'ART-1001;Screw M4\nART-1002\n"WIFI:T:WPA;S:Guest;P:secret;;";Guest Wi-Fi',

  mobileTitle: 'Mobile interface',
  mobileText:
    '— camera scanning, history, phone-style generating. It is the default view on iPhone, iPad and Android; you can try it here. Go back via "About".',
  mobileButton: 'Mobile interface',

  licenseTitle: 'License & copyright',
  licenseA: 'Published as open-source software under the',
  licenseName: 'MIT license with attribution',
  licenseB:
    '(MIT with Attribution). The source code may be freely used, modified and distributed — commercially too. Products using it must visibly credit the author "Michael Lechner" (e.g. in the documentation or an about dialog). A commercial license without attribution is available on request.',
  licenseLine: 'License: MIT with Attribution (see below)',
  showLicense: 'Show license text',
  showThirdParty: 'Show third-party licenses (ZXing, boombuler/barcode, gozxing)',
  libsTitle: 'Open-source libraries used:',
  libWails: 'Desktop app framework',
  libSvelte: 'Reactive UI framework',
  libBoombuler: 'Barcode engine',
  libGozxing: 'Barcode reading',
  libPdf417: 'ported to Go',
  libMcp: 'Model Context Protocol',
  libBootstrap: 'UI styling & icons',
  libInter: 'Interface font',
  libGoImage: 'Captions in PNG',
  libGoImageName: 'golang.org/x/image & Go fonts',
};

export const t = translator(de, en);
export type AboutKey = keyof typeof de;
