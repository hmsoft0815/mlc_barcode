# Krypto-Zahlungs-QR-Format

Zahlungs-URIs für Krypto-Wallets.

Standards: **BIP 21** (Bitcoin und kompatible Ketten) sowie **EIP-681** für Ethereum.

## Parameter

| CLI-Flag | MCP-Parameter | Beschreibung |
|----------|---------------|-------------|
| `-crypto-coin` | `coin` | Währung: `bitcoin`, `ethereum`, `solana`, `dogecoin`, `litecoin` (Standard: `bitcoin`). |
| `-crypto-addr` | `address` | **Pflicht.** Die Wallet-Adresse. |
| `-crypto-amount` | `amount` | Optionaler Betrag in der Einheit der jeweiligen Währung. |
| (N/A) | `label` | Optionale Empfängerbezeichnung. Nur über MCP verfügbar. |
| `-crypto-msg` | `message` | Optionale Zahlungsnachricht oder Memo. |

## Anwendungsbeispiele

### CLI-Beispiel
```bash
./bin/barcode -crypto-coin bitcoin -crypto-addr "bc1qexampleaddress" \
  -crypto-amount 0.005 -crypto-msg "Rechnung 1002" -out zahlung.png
```

### MCP-Tool-Beispiel
```json
{
  "name": "generate_crypto_qr",
  "arguments": {
    "coin": "ethereum",
    "address": "0xExampleAddress",
    "amount": 0.25
  }
}
```

## Hinweise

**Die Adresse wird nicht geprüft.** Jede Kette hat ihr eigenes Format und ihre eigene Prüfsumme, und eine falsche Adresse bedeutet unwiederbringlich verlorene Mittel. Die Adresse gehört vor dem Druck kontrolliert, nicht danach.

`label` gibt es als MCP-Parameter, aber nicht als CLI-Flag — bislang hat es auf der Kommandozeile schlicht niemand gebraucht.
