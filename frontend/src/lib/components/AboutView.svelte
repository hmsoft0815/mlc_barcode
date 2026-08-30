<script lang="ts">
  import { BARCODE_TYPES } from '../types';
  export let appVersion: string = '1.3.0';
</script>

<div class="container-fluid py-3" style="max-width: 900px;">
  <div class="card shadow-sm border-0 mb-3">
    <div class="card-body p-4 text-center">
      <div class="display-6 text-primary mb-2">
        <i class="bi bi-qr-code"></i> MLC Barcode
      </div>
      <p class="lead text-muted mb-2">
        Leichtgewichtiges, schnelles Desktop-Tool zur Generierung und Stapelverarbeitung von Barcodes & QR-Codes.
      </p>
      <div>
        <span class="badge bg-primary me-2">Version {appVersion}</span>
        <span class="badge bg-secondary-subtle text-secondary border">Wails v3 + Svelte + Go</span>
      </div>
    </div>
  </div>

  <div class="card shadow-sm border-0 mb-3">
    <div class="card-header bg-white border-bottom py-2">
      <h6 class="mb-0 fw-semibold text-secondary">
        <i class="bi bi-card-checklist me-1"></i> Unterstützte Symbologien & Formate
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
                <td class="fw-medium">{t.name}</td>
                <td>
                  <span class="badge {t.category === '2D Matrix' ? 'bg-info-subtle text-info-emphasis' : 'bg-light text-dark border'}">
                    {t.category}
                  </span>
                </td>
                <td class="small text-muted">{t.description}</td>
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
      <div class="card shadow-sm border-0 h-100">
        <div class="card-header bg-white border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-secondary">
            <i class="bi bi-lightning-charge me-1"></i> Features
          </h6>
        </div>
        <div class="card-body small text-muted">
          <ul class="mb-0 ps-3">
            <li><strong>Live SVG Preview:</strong> Direkte Vektoranzeige in Echtzeit bei jeder Tastatureingabe.</li>
            <li><strong>Strukturierte QR-Formate:</strong> Masken für WLAN, vCard 3.0 und iCal Termine.</li>
            <li><strong>Batch-Generierung:</strong> Import von Textdateien (Zeile = Barcode) mit direktem Ordner-Export.</li>
            <li><strong>Etikettendruck:</strong> DIN A4 Raster für Standard-Etikettenbögen (Avery/Zweckform).</li>
            <li><strong>Clipboard:</strong> 1-Click Kopieren als SVG-Code oder direkt als PNG-Bild.</li>
          </ul>
        </div>
      </div>
    </div>

    <div class="col-md-6">
      <div class="card shadow-sm border-0 h-100">
        <div class="card-header bg-white border-bottom py-2">
          <h6 class="mb-0 fw-semibold text-secondary">
            <i class="bi bi-terminal me-1"></i> CLI & MCP Server
          </h6>
        </div>
        <div class="card-body small text-muted">
          <p class="mb-2">
            Neben dieser Desktop-GUI enthält das Projekt auch:
          </p>
          <ul class="mb-0 ps-3">
            <li><code>barcode</code>: CLI-Tool für Shell-Skripte & Automatisierung.</li>
            <li><code>mcp-barcode-server</code>: Model Context Protocol Server für LLM- und Agenten-Integrationen.</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>
