// Native actions with fallbacks: the Go side (Wails' iOS/Android bridges)
// first, then the browser's own means, then the clipboard.
import { Browser } from '@wailsio/runtime';
import {
  CopyToClipboard,
  Haptic,
  OpenExternal,
  ShareImage,
  ShareText,
} from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/barcodeapp';

export async function openLink(url: string) {
  if (await OpenExternal(url).catch(() => false)) return;
  await Browser.OpenURL(url).catch(() => window.open(url, '_blank'));
}

export async function copyText(text: string) {
  if (!(await CopyToClipboard(text).catch(() => false))) await navigator.clipboard.writeText(text);
}

// Returns what happened, so the UI can say "kopiert" when sharing was not
// possible.
export async function shareText(text: string): Promise<'shared' | 'copied'> {
  if (await ShareText(text).catch(() => false)) return 'shared';
  if (navigator.share) {
    try {
      await navigator.share({ text });
      return 'shared';
    } catch {
      // cancelled or not allowed: fall through
    }
  }
  await copyText(text);
  return 'copied';
}

export async function shareImage(pngDataUrl: string, name: string, text: string): Promise<'shared' | 'copied'> {
  if (await ShareImage(pngDataUrl, name).catch(() => false)) return 'shared';
  try {
    const blob = await (await fetch(pngDataUrl)).blob();
    const file = new File([blob], `${name}.png`, { type: 'image/png' });
    if (navigator.canShare?.({ files: [file] })) {
      await navigator.share({ files: [file] });
      return 'shared';
    }
  } catch {
    // fall through to sharing the content as text
  }
  return shareText(text);
}

export function haptic(kind: 'success' | 'warning' | 'error' | 'selection') {
  Haptic(kind).catch(() => navigator.vibrate?.(kind === 'success' ? 30 : 60));
}

// Links that the result screen can open: web pages, mail, phone, maps.
export function linkOf(text: string): string | null {
  const t = text.trim();
  return /^(https?:\/\/|mailto:|tel:|geo:)\S+$/i.test(t) ? t : null;
}
