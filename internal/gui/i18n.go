package gui

import "sync"

// Texts the Go side shows itself: native file dialogs and batch export
// errors. The frontend reports its language with SetLanguage; errors of
// the barcode engine keep going out as codes (ErrorInfo) and are worded by
// the frontend.

var (
	langMu  sync.RWMutex
	uiLang  = "de"
	uiTexts = map[string]map[string]string{
		"de": {
			"saveTitle":      "Barcode speichern unter …",
			"filterSVG":      "SVG-Vektorgrafik (*.svg)",
			"filterPNG":      "PNG-Bild (*.png)",
			"exportFolder":   "Zielordner für Barcode-Dateien auswählen",
			"tableTitle":     "Text- oder CSV-Datei auswählen",
			"filterTable":    "Text- und CSV-Dateien (*.txt, *.csv)",
			"filterAll":      "Alle Dateien (*.*)",
			"noExportFolder": "Zielordner ist nicht angegeben",
			"cannotMakeDir":  "Ordner konnte nicht erstellt werden: %v",
		},
		"en": {
			"saveTitle":      "Save barcode as …",
			"filterSVG":      "SVG vector graphic (*.svg)",
			"filterPNG":      "PNG image (*.png)",
			"exportFolder":   "Choose the target folder for the barcode files",
			"tableTitle":     "Choose a text or CSV file",
			"filterTable":    "Text and CSV files (*.txt, *.csv)",
			"filterAll":      "All files (*.*)",
			"noExportFolder": "No target folder given",
			"cannotMakeDir":  "Folder could not be created: %v",
		},
	}
)

// SetLanguage tells the Go side the UI language ("de" or "en") for the
// native dialogs and the messages it words itself.
func (a *BarcodeApp) SetLanguage(lang string) {
	if _, ok := uiTexts[lang]; !ok {
		return
	}
	langMu.Lock()
	uiLang = lang
	langMu.Unlock()
}

// tr is the text for key in the current UI language.
func tr(key string) string {
	langMu.RLock()
	defer langMu.RUnlock()
	return uiTexts[uiLang][key]
}
