// Texts of the "My details" form (ProfileForm.svelte), desktop and mobile.
import { translator } from '../lang';

const de = {
  title: 'Meine Angaben',
  intro: 'Werden beim Erzeugen von Terminen, Orten und Kontakten vorausgefüllt. Sie bleiben auf diesem Gerät.',
  lastName: 'Nachname',
  firstName: 'Vorname',
  email: 'E-Mail',
  phone: 'Telefon',
  organizerPreview: 'Als Organisator: {name}',
  place: 'Standort (optional)',
  location: 'Adresse oder Ortsname',
  locationPh: 'z. B. Praxis Dr. Meier, Hauptstr. 1, Berlin',
  latitude: 'Breitengrad',
  longitude: 'Längengrad',
  useMyPosition: 'Aktuellen Standort übernehmen',
  locating: 'Standort wird bestimmt …',
  positionFailed: 'Standort nicht verfügbar ({reason}).',
  timeZone: 'Zeitzone',
  timeZoneDevice: 'Gerät: {tz}',
  clear: 'Angaben löschen',
  saved: 'Gespeichert',
  done: 'Fertig',
};

const en: Record<keyof typeof de, string> = {
  title: 'My details',
  intro: 'Filled in when you create events, places and contacts. They stay on this device.',
  lastName: 'Last name',
  firstName: 'First name',
  email: 'Email',
  phone: 'Phone',
  organizerPreview: 'As organizer: {name}',
  place: 'Location (optional)',
  location: 'Address or place name',
  locationPh: 'e.g. Dr Meier practice, 1 Main St, Berlin',
  latitude: 'Latitude',
  longitude: 'Longitude',
  useMyPosition: 'Use my current position',
  locating: 'Finding your position …',
  positionFailed: 'Position not available ({reason}).',
  timeZone: 'Time zone',
  timeZoneDevice: 'Device: {tz}',
  clear: 'Clear details',
  saved: 'Saved',
  done: 'Done',
};

export const t = translator(de, en);
