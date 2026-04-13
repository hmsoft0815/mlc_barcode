# MLC Barcode

A powerful utility for generating barcodes, available both as a standalone Command Line Interface (CLI) and a Model Context Protocol (MCP) Server.

## What it does

MLC Barcode is designed to seamlessly integrate barcode generation into your daily workflows and automated pipelines. Whether you need to generate a simple QR code from your terminal, automate labels, or equip your favorite Large Language Models with the ability to create barcodes directly in the chat, MLC Barcode has you covered.

## Key features

- **Versatile Barcode Support**: Generate QR, DataMatrix, Code128, Code39, EAN-13, EAN-8, UPC-A, and ITF formats.
- **LLM Integration via MCP**: The built-in MCP server provides the `generate_barcode` tool to AI clients like Claude Desktop, enabling natural language generation of barcodes.
- **Artifact Integration**: Generated files can be directly saved to the `mlcartifact` service, keeping your local workspace clean.
- **Smart QR Code Templates**: Built-in support for generating specialized QR codes, such as automatically formatted WiFi payloads, vCards, and Calendar events.
- **High-Quality Output**: Render pixel-perfect vector graphics (`SVG`) or standard raster images (`PNG`) with adjustable sizing and transparent backgrounds.
