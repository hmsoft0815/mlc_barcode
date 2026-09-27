import type * as Models from '../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';

export type BarcodeType =
  | 'qr'
  | 'datamatrix'
  | 'aztec'
  | 'pdf417'
  | 'code128'
  | 'code39'
  | 'ean13'
  | 'ean8'
  | 'upca'
  | 'itf';

export interface BarcodeTypeOption {
  id: BarcodeType;
  name: string;
  category: '2D Matrix' | '2D Stacked' | '1D Linear';
  sample: string;
}

export const BARCODE_TYPES: BarcodeTypeOption[] = [
  { id: 'qr', name: 'QR Code', category: '2D Matrix', sample: 'https://mlcgo.eu' },
  { id: 'datamatrix', name: 'DataMatrix', category: '2D Matrix', sample: 'MLC-DM-12345' },
  { id: 'aztec', name: 'Aztec', category: '2D Matrix', sample: 'TICKET-ICE-599-FRA-MUC' },
  { id: 'pdf417', name: 'PDF417', category: '2D Stacked', sample: 'BOARDING PASS LH123 FRA-JFK' },
  { id: 'code128', name: 'Code 128', category: '1D Linear', sample: 'MLC-128-ABC' },
  { id: 'code39', name: 'Code 39', category: '1D Linear', sample: 'CODE39' },
  { id: 'ean13', name: 'EAN-13 / GTIN-13', category: '1D Linear', sample: '4012345678901' },
  { id: 'ean8', name: 'EAN-8 / GTIN-8', category: '1D Linear', sample: '40123455' },
  { id: 'upca', name: 'UPC-A', category: '1D Linear', sample: '012345678905' },
  { id: 'itf', name: 'ITF (Interleaved 2 of 5)', category: '1D Linear', sample: '12345678' }
];

// One label on the print sheet; data is the text printed under the code.
export interface PrintItem {
  data: string;
  svg: string;
  type: string;
}

export interface PrintLabelConfig {
  columns: number;
  rows: number;
  labelWidthMm: number;
  labelHeightMm: number;
  fontSizePt: number;
  showText: boolean;
}
