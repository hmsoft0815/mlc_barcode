// Texts of the PZN lookup button (PznLookup.svelte).
import { translator } from '../lang';

const de = {
  button: 'Beim BfArM nachschlagen',
  copied:
    'PZN {pzn} kopiert. In der AMIce-Suche die Bedingungen akzeptieren, bei „in“ „Pharmazentralnummer“ wählen und die PZN einfügen.',
};

const en: Record<keyof typeof de, string> = {
  button: 'Look up at BfArM',
  copied:
    'PZN {pzn} copied. In the AMIce search accept the terms, choose "Pharmazentralnummer" under "in" and paste the PZN.',
};

export const t = translator(de, en);
