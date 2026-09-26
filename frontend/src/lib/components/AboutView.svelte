<script lang="ts">
  // Braces would be Svelte syntax in the markup, so the example lives here.
  const cliAndMcpExample = `# Kommandozeile:
barcode -type qr -data "https://mlcgo.eu" -out qr.svg

# MCP-Server: startet Ihr KI-Client selbst.
# Eintrag z. B. in claude_desktop_config.json (Windows-Setup):
"mcpServers": {
  "mlc-barcode": {
    "command": "C:\\\\Program Files\\\\MLC Barcode\\\\mcp-barcode-server.exe"
  }
}`;
  import licenseText from '../../../../LICENSE?raw';
  import { BARCODE_TYPES } from '../types';
  export let appVersion: string = __APP_VERSION__;
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
        <span class="badge bg-success-subtle text-success-emphasis border">MIT with Attribution</span>
        <span class="badge bg-secondary-subtle text-secondary border">Wails v3 + Svelte + Go</span>
      </div>
      <div class="d-flex justify-content-center gap-2">
        <a
          href="https://github.com/hmsoft0815/mlc_barcode"
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
                  <span class="badge {t.category !== '1D Linear' ? 'bg-info-subtle text-info-emphasis' : 'bg-body-secondary text-body border'}">
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
            <li><strong>9 QR-Sonderformate:</strong> GiroCode (SEPA-Überweisung), Krypto-Wallets, Maps (Geo), Telefon, SMS, E-Mail, WLAN, vCard 3.0 und Kalender-Events.</li>
            <li><strong>Prüfen:</strong> Barcodes aus Bildern lesen (Datei, Drag &amp; Drop, <kbd>Strg</kbd>+<kbd>V</kbd>) – Format, Inhalt und zerlegte Felder, beim GiroCode mit IBAN-Prüfung. Jeder erzeugte Code wird automatisch gegengelesen („Lesbar geprüft“).</li>
            <li><strong>Beschriftung:</strong> Freitext unter dem Barcode mit einstellbarer Schriftgröße, in SVG und PNG.</li>
            <li><strong>Batch-Generierung:</strong> Import von <code>.txt</code>/<code>.csv</code> (Format siehe unten) und Stapel-Export in einen Zielordner.</li>
            <li><strong>Etiketten-Druckbogen:</strong> DIN-A4-Layouts nach gängigen Avery-Zweckform-Rastern; Etiketten aus dem Batch-Generator übernehmen oder direkt aus TXT/CSV mit eigenem Etikett-Text importieren.</li>
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
          <pre class="bg-body-secondary p-2 rounded mb-0 font-monospace small"><code>{cliAndMcpExample}</code></pre>
          <a
            href="https://mlcgo.eu/products/mlc-barcode/de/mcp/"
            target="_blank"
            rel="noreferrer"
            class="btn btn-sm btn-outline-primary mt-2"
          >
            <i class="bi bi-plug me-1"></i> Anleitung: MCP-Server in Claude, Gemini &amp; Cursor einbinden
          </a>
        </div>
      </div>
    </div>
  </div>

  <!-- File import format -->
  <div class="card shadow-sm border mb-3">
    <div class="card-header bg-body border-bottom py-2">
      <h6 class="mb-0 fw-semibold text-body">
        <i class="bi bi-filetype-csv me-1 text-primary"></i> Datei-Import (Batch-Generator & Etiketten-Druck)
      </h6>
    </div>
    <div class="card-body small text-body-secondary">
      <ul class="ps-3 mb-2">
        <li><strong>TXT:</strong> Jede Zeile ist ein Barcode-Inhalt. Es wird nichts aufgeteilt — Kommas und Semikolons gehören zum Inhalt.</li>
        <li><strong>CSV:</strong> Spalte 1 = Inhalt (z.&nbsp;B. für den QR-Code), Spalte 2 = <em>optionaler</em> Etikett-Text. Der Batch-Generator nutzt nur Spalte 1.</li>
        <li><strong>Trennzeichen:</strong> <code>;</code> oder Tab werden erkannt (Excel speichert deutsch mit <code>;</code>). Ein Komma gilt nur als Trennzeichen, wenn jede Zeile gleich viele hat. Ohne erkennbares Trennzeichen ist jede Zeile ein Inhalt.</li>
        <li>Enthält der Inhalt selbst das Trennzeichen (z.&nbsp;B. <code>WIFI:…;…;</code>), das Feld in Anführungszeichen setzen.</li>
        <li>Leere Zeilen werden übersprungen, eine Kopfzeile lässt sich im Etiketten-Druck abschalten.</li>
      </ul>
      <pre class="bg-body-secondary p-2 rounded mb-0 font-monospace small"><code>ART-1001;Schraube M4
ART-1002
"WIFI:T:WPA;S:Gast;P:geheim;;";Gäste-WLAN</code></pre>
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
            © 2026 Michael Lechner
          </p>
          <p class="mb-2">
            Veröffentlicht als Open-Source-Software unter der <strong>MIT-Lizenz mit Namensnennung</strong>
            (MIT with Attribution). Der Quellcode darf frei genutzt, verändert und weitergegeben werden —
            auch kommerziell. Produkte, die ihn verwenden, müssen den Autor „Michael Lechner“ sichtbar
            nennen (z.&nbsp;B. in der Dokumentation oder einem Info-Dialog). Eine kommerzielle Lizenz
            ohne Namensnennung ist auf Anfrage erhältlich.
          </p>
          <div class="font-monospace p-2 bg-body-secondary rounded border small">
            Lizenz: MIT with Attribution (siehe unten)<br />
            Git Repository: <a href="https://github.com/hmsoft0815/mlc_barcode" target="_blank" rel="noreferrer" class="text-decoration-none">https://github.com/hmsoft0815/mlc_barcode</a>
          </div>
          <details class="mt-2">
            <summary class="small text-primary" style="cursor: pointer;">Lizenztext anzeigen</summary>
            <pre class="license-text bg-body-secondary border rounded p-2 mt-2 mb-0">{licenseText}</pre>
          </details>
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
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>MCP Go-SDK</strong> (Apache-2.0 / MIT) — Model Context Protocol
            </li>
            <li>
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>Bootstrap 5.3 &amp; Bootstrap Icons</strong> (MIT) — UI Styling & Icons
            </li>
            <li class="mt-1">
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>Inter</strong> (SIL OFL 1.1) — Schrift der Oberfläche
            </li>
            <li class="mt-1">
              <i class="bi bi-box-seam me-1 text-primary"></i> <strong>golang.org/x/image &amp; Go-Schriften</strong> (BSD-3-Clause) — Beschriftung im PNG
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .license-text {
    white-space: pre-wrap;
    font-size: 0.75rem;
    max-height: 22rem;
    overflow-y: auto;
  }
</style>
