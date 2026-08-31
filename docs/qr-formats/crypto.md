# Cryptocurrency Payment QR Format

Payment URIs for cryptocurrency wallets.

Standards: **BIP 21** (Bitcoin and compatible chains) and **EIP-681** for Ethereum.

## Parameters

| CLI Flag | MCP Parameter | Description |
|----------|---------------|-------------|
| `-crypto-coin` | `coin` | Coin: `bitcoin`, `ethereum`, `solana`, `dogecoin`, `litecoin` (default: `bitcoin`). |
| `-crypto-addr` | `address` | **Required.** The wallet address. |
| `-crypto-amount` | `amount` | Optional amount, in the coin's own unit. |
| (N/A) | `label` | Optional recipient label. Available through MCP only. |
| `-crypto-msg` | `message` | Optional payment message or memo. |

## Usage Examples

### CLI Example
```bash
./bin/barcode -crypto-coin bitcoin -crypto-addr "bc1qexampleaddress" \
  -crypto-amount 0.005 -crypto-msg "Invoice 1002" -out pay.png
```

### MCP Tool Example
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

## Notes

**The address is never validated.** Every chain has its own format and checksum, and a wrong address means funds are gone with no way back. Check the address before you print the code, not after.

`label` exists as an MCP parameter but has no CLI flag — the asymmetry is deliberate rather than an oversight only insofar as nobody has needed it on the command line yet.
