package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerPrompts adds workflow prompts: they tell the model which tool to
// call with which arguments for the two jobs people ask for most.
func registerPrompts(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Name:        "payment_qr_from_invoice",
		Title:       "GiroCode from an invoice",
		Description: "Read an invoice and create the matching SEPA payment QR code (GiroCode / EPC QR)",
		Arguments: []*mcp.PromptArgument{
			{Name: "invoice", Description: "Invoice text (pasted or extracted)", Required: true},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		invoice := strings.TrimSpace(req.Params.Arguments["invoice"])
		if invoice == "" {
			return nil, missingArgument("invoice")
		}
		return userPrompt("Payment QR code from invoice", `From the invoice below, extract:
- beneficiary name (the company or person to be paid),
- IBAN (remove spaces), BIC if printed,
- amount in EUR (the total due, not a subtotal),
- payment reference (invoice number, or the reference the invoice asks for).

Then call generate_epc_qr with name, iban, bic, amount and reference, and caption:"<amount> € · <reference>".
If the IBAN or the amount is missing or ambiguous, ask instead of guessing — a wrong IBAN sends money to the wrong account.

Invoice:
`+invoice), nil
	})

	s.AddPrompt(&mcp.Prompt{
		Name:        "product_labels",
		Title:       "Product labels",
		Description: "Create one labelled barcode per product line, e.g. for shelf or stock labels",
		Arguments: []*mcp.PromptArgument{
			{Name: "items", Description: "One product per line: code, optionally followed by ';' and a label text", Required: true},
			{Name: "type", Description: "Symbology: ean13 (default), ean8, upca, code128, qr, datamatrix"},
		},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		items := strings.TrimSpace(req.Params.Arguments["items"])
		if items == "" {
			return nil, missingArgument("items")
		}
		btype := strings.TrimSpace(req.Params.Arguments["type"])
		if btype == "" {
			btype = "ean13"
		}
		return userPrompt("Product labels", fmt.Sprintf(`For every line below call generate_barcode with type:%s and data set to the code.
If the line has a label text after ';', pass it as caption, otherwise set text:true.
EAN/UPC codes may lack the check digit — it is computed; report the completed code from structuredContent.encoded_data.
If a code is rejected, list the line and the reason at the end instead of stopping.

Lines:
%s`, btype, items)), nil
	})
}

func userPrompt(description, text string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: description,
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}
}

// missingArgument is the JSON-RPC invalid-params error the spec asks for
// when a required prompt argument is missing.
func missingArgument(name string) error {
	return &jsonrpc.Error{Code: jsonrpc.CodeInvalidParams, Message: "missing required argument: " + name}
}
