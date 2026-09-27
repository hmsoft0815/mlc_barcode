package gui

import "testing"

// Every text of the Go side exists in German and English.
func TestUITextsComplete(t *testing.T) {
	for lang, texts := range uiTexts {
		for other, otherTexts := range uiTexts {
			for key := range texts {
				if otherTexts[key] == "" {
					t.Errorf("%q is in %s but missing in %s", key, lang, other)
				}
			}
		}
	}
	a := NewBarcodeApp()
	a.SetLanguage("en")
	if got := tr("saveTitle"); got != "Save barcode as …" {
		t.Errorf("en: %q", got)
	}
	a.SetLanguage("xx") // unknown: stays English
	if got := tr("saveTitle"); got != "Save barcode as …" {
		t.Errorf("unknown language changed the texts: %q", got)
	}
	a.SetLanguage("de")
}
