// Date and time helpers for events, shared by the desktop generator and
// the mobile "create" view. Values are held as <input type="datetime-local">
// strings (YYYY-MM-DDTHH:MM, local time).

const pad2 = (n: number) => String(n).padStart(2, '0');

export function toLocalInput(d: Date): string {
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}T${pad2(d.getHours())}:${pad2(d.getMinutes())}`;
}

export function nextFullHour(): Date {
  const d = new Date();
  d.setHours(d.getHours() + 1, 0, 0, 0);
  return d;
}

export function shiftLocal(value: string, ms: number): string {
  const d = new Date(value);
  return isNaN(d.getTime()) ? value : toLocalInput(new Date(d.getTime() + ms));
}

// YYYY-MM-DDTHH:MM → YYYYMMDDTHHMM00
export function icalDateTime(value: string): string {
  return value ? value.replace(/[-:]/g, '').slice(0, 13) + '00' : '';
}

// YYYY-MM-DD… → YYYYMMDD
export function icalDate(value: string): string {
  return value ? value.slice(0, 10).replace(/-/g, '') : '';
}

// Date inputs only change the date part and keep the time.
export function withDate(current: string, date: string): string {
  return date ? date + current.slice(10) : current;
}

export const HOUR_MS = 3_600_000;
export const DAY_MS = 24 * HOUR_MS;

// The device's IANA time zone, e.g. Europe/Berlin.
export function deviceTimeZone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone || '';
  } catch {
    return '';
  }
}
