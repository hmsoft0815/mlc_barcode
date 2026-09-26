// Scan history of the mobile UI: the last scans, newest first, kept in
// localStorage on the device (per viewer, nothing leaves the device).
import { writable } from 'svelte/store';
import type { DecodedCode } from '../../../bindings/github.com/mlcmcp/mlc_barcode/internal/gui/models';

export interface HistoryEntry {
  id: string;
  at: number; // ms since epoch
  code: DecodedCode;
}

const KEY = 'mlc_scan_history';
const MAX_ENTRIES = 200;
// The same code scanned again within this time is one scan, not two.
const REPEAT_MS = 10_000;

function load(): HistoryEntry[] {
  try {
    const raw = localStorage.getItem(KEY);
    const list = raw ? JSON.parse(raw) : [];
    return Array.isArray(list) ? list : [];
  } catch {
    return [];
  }
}

export const history = writable<HistoryEntry[]>(load());

history.subscribe((list) => {
  try {
    localStorage.setItem(KEY, JSON.stringify(list));
  } catch {
    // storage full or blocked: the list still works for this session
  }
});

export function addScan(code: DecodedCode) {
  history.update((list) => {
    const last = list[0];
    if (last && last.code.type === code.type && last.code.text === code.text && Date.now() - last.at < REPEAT_MS) {
      return list;
    }
    const entry = { id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`, at: Date.now(), code };
    return [entry, ...list].slice(0, MAX_ENTRIES);
  });
}

export function removeScan(id: string) {
  history.update((list) => list.filter((e) => e.id !== id));
}

export function clearHistory() {
  history.set([]);
}
