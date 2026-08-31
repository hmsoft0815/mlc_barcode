<script lang="ts">
  import { BARCODE_TYPES } from '../types';
  export let appVersion: string = '1.3.0';
</script>

<div class="container-fluid py-3" style="max-width: 900px;">
  <div class="card shadow-sm border mb-3">
    <div class="card-body p-4 text-center">
      <div class="display-6 text-primary mb-2">
        <i class="bi bi-qr-code"></i> MLC Barcode
      </div>
      <p class="lead text-body-secondary mb-2">
        Leichtgewichtiges, schnelles Desktop-Tool zur Generierung und Stapelverarbeitung von Barcodes & QR-Codes.
      </p>
      <div>
        <span class="badge bg-primary me-2">Version {appVersion}</span>
        <span class="badge bg-secondary-subtle text-secondary border">Wails v3 + Svelte + Go</span>
      </div>
    </div>
  </div>

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

  <div class="row g-3">
    <div class="col-md-6">
      <div class="card shadow-sm border h-100">
        <div class="card-header bg-body border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-body">
            <i class="bi bi-gear-fill me-1 text-primary"></i> Funktionen im Überblick
          </h6>
        </div>
        <div class="card-body small text-body-secondary">
          <ul class="mb-0 ps-3">
            <li><strong>Live-Vorschau:</strong> Direkte Vektordarstellung in SVG und PNG.</li>
            <li><strong>QR-Sonderformate:</strong> Assistenten für WLAN-Zugangsdaten, vCards und Kalender-Events.</li>
            <li><strong>Batch-Generierung:</strong> Import von `.txt`/`.csv` Textdateien und Stapel-Export in Zielordner.</li>
            <li><strong>Etiketten-Druck:</strong> DIN A4 Bogenlayouts nach gängigen Avery-Zweckform-Rastern.</li>
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
            MLC Barcode ist auch als stand-alone Kommandozeilentool und MCP-Server für KI-Assistenten verfügbar:
          </p>
          <pre class="bg-body-secondary p-2 rounded mb-0 font-monospace small"><code># CLI Aufruf:
barcode -type qr -data "https://mlcgo.eu" -out qr.svg

# MCP Server:
mcp-barcode-server</code></pre>
        </div>
      </div>
    </div>
  </div>
</div>
