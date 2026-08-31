package qrformats

import (
	"fmt"
	"net/url"
	"strings"
)

// CryptoOptions holds configuration for cryptocurrency payment URIs.
type CryptoOptions struct {
	Coin    string  // "bitcoin", "ethereum", "solana", "dogecoin", "litecoin", etc.
	Address string  // Wallet address
	Amount  float64 // Optional amount
	Label   string  // Optional label / recipient name
	Message string  // Optional message / memo
}

// FormatCrypto returns standard URI schemes (BIP 21 for Bitcoin, EIP-681 for Ethereum, etc.).
func FormatCrypto(opts CryptoOptions) string {
	coin := strings.ToLower(strings.TrimSpace(opts.Coin))
	if coin == "" || coin == "btc" {
		coin = "bitcoin"
	} else if coin == "eth" {
		coin = "ethereum"
	} else if coin == "sol" {
		coin = "solana"
	} else if coin == "doge" {
		coin = "dogecoin"
	} else if coin == "ltc" {
		coin = "litecoin"
	}

	address := strings.TrimSpace(opts.Address)
	if address == "" {
		return ""
	}

	params := url.Values{}
	if opts.Amount > 0 {
		// Format amount without scientific notation
		params.Set("amount", fmt.Sprintf("%g", opts.Amount))
	}
	if opts.Label != "" {
		params.Set("label", opts.Label)
	}
	if opts.Message != "" {
		params.Set("message", opts.Message)
	}

	queryString := params.Encode()
	if queryString != "" {
		return fmt.Sprintf("%s:%s?%s", coin, address, queryString)
	}
	return fmt.Sprintf("%s:%s", coin, address)
}
