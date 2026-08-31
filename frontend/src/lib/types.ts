import type * as Models from '../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';

export type BarcodeType = 'qr' | 'datamatrix' | 'code128' | 'code39' | 'ean13' | 'ean8' | 'upca' | 'itf';

export interface BarcodeTypeOption {
  id: BarcodeType;
  name: string;
  category: '2D Matrix' | '1D Linear';
  description: string;
  sample: string;
}

export const BARCODE_TYPES: BarcodeTypeOption[] = [
  { id: 'qr', name: 'QR Code', category: '2D Matrix', description: '2D Matrixcode für Text, URLs, WLAN, vCard, Events', sample: 'https://mlcgo.eu' },
  { id: 'datamatrix', name: 'DataMatrix', category: '2D Matrix', description: 'Kompakter 2D-Code für Industrie und Bauteile', sample: 'MLC-DM-12345' },
  { id: 'code128', name: 'Code 128', category: '1D Linear', description: 'Universeller 1D-Barcode für alle ASCII-Zeichen', sample: 'MLC-128-ABC' },
  { id: 'code39', name: 'Code 39', category: '1D Linear', description: 'Alphanumerischer Barcode (Großbuchstaben, Ziffern)', sample: 'CODE39' },
  { id: 'ean13', name: 'EAN-13 / GTIN-13', category: '1D Linear', description: '13-stelliger Einzelhandels-Barcode (12 Ziffern + Prüfziffer)', sample: '4012345678901' },
  { id: 'ean8', name: 'EAN-8 / GTIN-8', category: '1D Linear', description: '8-stelliger kompakter Handels-Barcode', sample: '40123455' },
  { id: 'upca', name: 'UPC-A', category: '1D Linear', description: '12-stelliger US-Einzelhandels-Barcode', sample: '012345678905' },
  { id: 'itf', name: 'ITF (Interleaved 2 of 5)', category: '1D Linear', description: 'Kompakter numerischer Barcode (nur Ziffern, gerade Anzahl)', sample: '12345678' }
];

export interface PrintLabelConfig {
  columns: number;
  rows: number;
  labelWidthMm: number;
  labelHeightMm: number;
  fontSizePt: number;
  showText: boolean;
}
