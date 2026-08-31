<script lang="ts">
  import { BARCODE_TYPES } from '../types';
  export let appVersion: string = '1.3.0';
</script>

<div class="container-fluid py-3" style="max-width: 900px;">
  <!-- Header Card -->
  <div class="card shadow-sm border mb-3">
    <div class="card-body p-4 text-center">
      <div class="mb-3">
        <img src="/appicon.png" width="72" height="72" class="rounded-3 shadow" alt="MLC Barcode Logo" />
      </div>
      <h3 class="fw-bold text-body mb-1">MLC Barcode</h3>
      <p class="lead text-body-secondary mb-2">
        Plattformübergreifendes Desktop-Werkzeug zur Generierung, Stapelverarbeitung und dem Druck von Barcodes & QR-Codes.
      </p>
      <div class="d-flex flex-wrap justify-content-center gap-2 mb-3">
        <span class="badge bg-primary">Version {appVersion}</span>
        <span class="badge bg-success-subtle text-success-emphasis border">MIT Lizenz</span>
        <span class="badge bg-secondary-subtle text-secondary border">Wails v3 + Svelte + Go</span>
      </div>
      <div class="d-flex justify-content-center gap-2">
        <a
          href="https://github.com/mlcmcp/mlc_barcode"
          target="_blank"
          rel="noreferrer"
          class="btn btn-sm btn-outline-secondary"
        >
          <i class="bi bi-github me-1"></i> GitHub Repository
        </a>
        <a
          href="https://mlcgo.eu/products/mlc-barcode/"
          target="_blank"
          rel="noreferrer"
          class="btn btn-sm btn-outline-primary"
        >
          <i class="bi bi-globe me-1"></i> Produktseite (mlcgo.eu)
        </a>
      </div>
    </div>
  </div>

  <!-- Supported Symbologies -->
  <div class="card shadow-sm border mb-3">
    <div class="card-header bg-body border-bottom py-2">
      <h6 class="mb-0 fw-semibold text-body">
        <i class="bi bi-card-checklist me-1 text-primary"></i> Unterstützte Symbologien & Formate
      </h6>
    </div>
    <div class="card-body p-0">
      <div class="table-responsive">
        <table class="table table-hover table-sm mb-0">
          <thead class="table-light">
            <tr>
              <th style="width: 140px;">Typ</th>
              <th style="width: 110px;">Kategorie</th>
              <th>Beschreibung & Spezifikation</th>
              <th>Beispiel</th>
            </tr>
          </thead>
          <tbody>
            {#each BARCODE_TYPES as t}
              <tr>
                <td class="fw-medium text-body">{t.name}</td>
                <td>
                  <span class="badge {t.category === '2D Matrix' ? 'bg-info-subtle text-info-emphasis' : 'bg-body-secondary text-body border'}">
                    {t.category}
                  </span>
                </td>
                <td class="small text-body-secondary">{t.description}</td>
                <td><code class="small">{t.sample}</code></td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <!-- Features & CLI -->
  <div class="row g-3 mb-3">
    <div class="col-md-6">
      <div class="card shadow-sm border h-100">
        <div class="card-header bg-body border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-gear-fill me-1 text-primary"></i> Funktionen im Überblick
          </h6>
        </div>
        <div class="card-body small text-body-secondary">
          <ul class="mb-0 ps-3">
            <li><strong>Live-Vorschau:</strong> Direkte Vektordarstellung in SVG und hochauflösendem PNG.</li>
            <li><strong>10 QR-Sonderformate:</strong> GiroCode (SEPA-Überweisung), Krypto-Wallets, Maps (Geo), Telefon, SMS, E-Mail, WLAN, vCard 3.0 und Kalender-Events.</li>
            <li><strong>Batch-Generierung:</strong> Import von `.txt`/`.csv` Textdateien und Stapel-Export in Zielordner.</li>
            <li><strong>Etiketten-Druckbogen:</strong> DIN A4 Layouts nach gängigen Avery-Zweckform-Rastern mit Sofortdruck.</li>
          </ul>
        </div>
      </div>
    </div>

    <div class="col-md-6">
      <div class="card shadow-sm border h-100">
        <div class="card-header bg-body border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-terminal me-1 text-primary"></i> CLI & MCP Server
          </h6>
        </div>
        <div class="card-body small text-body-secondary">
          <p class="mb-2">
            MLC Barcode ist auch als stand-alone Kommandozeilentool und MCP-Server für KI-Assistenten (Claude Desktop, Cursor, Gemini-CLI) verfügbar:
          </p>
          <pre class="bg-body-secondary p-2 rounded mb-0 font-monospace small"><code># CLI Aufruf:
barcode -type qr -data "https://mlcgo.eu" -out qr.svg

# MCP Server (Stdio/SSE):
mcp-barcode-server</code></pre>
        </div>
      </div>
    </div>
  </div>

  <!-- Copyright, License & Credits -->
  <div class="card shadow-sm border mb-3">
    <div class="card-header bg-body border-bottom py-2">
      <h6 class="mb-0 fw-semibold text-body">
        <i class="bi bi-shield-check me-1 text-success"></i> Lizenz & Copyright
      </h6>
    </div>
    <div class="card-body small text-body-secondary">
      <div class="row g-3 align-items-center">
        <div class="col-md-7">
          <p class="fw-semibold text-body mb-1">
            © 2026 Michael Lechner · Alle Rechte vorbehalten.
          </p>
          <p class="mb-2">
            Veröffentlicht als Open-Source-Software unter der <strong>MIT Lizenz</strong>. Der Quellcode ist frei zugänglich, darf modifiziert und weiterverteilt werden.
          </p>
          <div class="font-monospace p-2 bg-body-secondary rounded border small">
            SPDX-License-Identifier: MIT<br />
            Git Repository: <a href="https://github.com/mlcmcp/mlc_barcode" target="_blank" rel="noreferrer" class="text-decoration-none">https://github.com/mlcmcp/mlc_barcode</a>
          </div>
        </div>
        <div class="col-md-5 border-start-md ps-md-3">
          <h6 class="fw-semibold text-body small mb-2">Verwendete Open-Source Bibliotheken:</h6>
          <ul class="list-unstyled mb-0" style="font-size: 0.82rem;">
            <li class="mb-1">
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>Wails v3</strong> (MIT) — Desktop App Framework
            </li>
            <li class="mb-1">
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>Svelte 5</strong> (MIT) — Reactive UI Framework
            </li>
            <li class="mb-1">
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>boombuler/barcode</strong> (MIT) — Barcode Engine
            </li>
            <li class="mb-1">
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>MCP Go-SDK</strong> (MIT) — Model Context Protocol
            </li>
            <li>
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>Bootstrap 5.3</strong> (MIT) — UI Styling & Icons
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>
