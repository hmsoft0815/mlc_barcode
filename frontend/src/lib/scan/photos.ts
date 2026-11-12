// "From Photos": the system photo picker (macOS 13+, via internal/photokit).
// The photo comes back as a File, so it takes the same path as a file the
// user picked or dropped.
import { PhotoLibraryAvailable, PickFromPhotos } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';

// Asked once; false where there is no photo library picker.
export const photoLibrary: Promise<boolean> = PhotoLibraryAvailable().catch(() => false);

// null when the user closed the picker without choosing.
export async function pickFromPhotos(): Promise<File | null> {
  const photo = await PickFromPhotos();
  if (!photo?.dataUrl) return null;
  const [head, b64] = photo.dataUrl.split(',', 2);
  const mime = /data:([^;]+)/.exec(head)?.[1] ?? 'image/jpeg';
  const bin = atob(b64);
  const bytes = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
  return new File([bytes], photo.name || 'photo.jpg', { type: mime });
}
