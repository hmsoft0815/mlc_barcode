// "My details": who the user is and where, kept on this device only
// (localStorage). They fill in new event, place and contact codes: the
// organizer of an event, its location, coordinates and time zone.
import { writable } from 'svelte/store';

export interface Profile {
  lastName: string;
  firstName: string;
  email: string;
  phone: string;
  location: string; // e.g. "Praxis Dr. Meier, Hauptstr. 1, Berlin"
  latitude: string; // decimal degrees, kept as typed
  longitude: string;
  timeZone: string; // IANA, e.g. Europe/Berlin; empty = the device's
}

const KEY = 'mlc_profile';
const EMPTY: Profile = { lastName: '', firstName: '', email: '', phone: '', location: '', latitude: '', longitude: '', timeZone: '' };

function load(): Profile {
  try {
    const raw = localStorage.getItem(KEY);
    return raw ? { ...EMPTY, ...JSON.parse(raw) } : { ...EMPTY };
  } catch {
    return { ...EMPTY };
  }
}

export const profile = writable<Profile>(load());

profile.subscribe((p) => {
  try {
    localStorage.setItem(KEY, JSON.stringify(p));
  } catch {
    // storage blocked: the details still work for this session
  }
});

// The organizer as calendars show it: "Lechner, Michael".
export function organizerName(p: Profile): string {
  return [p.lastName.trim(), p.firstName.trim()].filter(Boolean).join(', ');
}

// The stored coordinates, or null when not set or not numbers in range.
export function profileCoords(p: Profile): { latitude: number; longitude: number } | null {
  const lat = Number(p.latitude.replace(',', '.'));
  const lon = Number(p.longitude.replace(',', '.'));
  if (p.latitude.trim() === '' || p.longitude.trim() === '' || !isFinite(lat) || !isFinite(lon)) return null;
  if (Math.abs(lat) > 90 || Math.abs(lon) > 180) return null;
  return { latitude: lat, longitude: lon };
}
