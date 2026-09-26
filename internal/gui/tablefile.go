package gui

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// TableFile is an imported TXT/CSV file. Every row has at least one cell:
// cell 0 is the barcode content, cell 1 (optional) a caption.
type TableFile struct {
	Path      string     `json:"path"`
	Separator string     `json:"separator"` // "" = one column per line
	Rows      [][]string `json:"rows"`
}

// PickTableFile opens a TXT/CSV file and splits it into rows. An empty Path
// means the user cancelled.
func (a *BarcodeApp) PickTableFile() (TableFile, error) {
	dialog := application.Get().Dialog.OpenFile()
	dialog.SetMessage("Text- oder CSV-Datei auswählen")
	dialog.AddFilter("Text & CSV Dateien (*.txt, *.csv)", "*.txt;*.csv")
	dialog.AddFilter("Alle Dateien (*.*)", "*.*")

	path, err := dialog.PromptForSingleSelection()
	if err != nil || path == "" {
		return TableFile{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return TableFile{Path: path}, fmt.Errorf("could not read file: %w", err)
	}
	sep, rows := parseTable(path, data)
	return TableFile{Path: path, Separator: sep, Rows: rows}, nil
}

// parseTable splits a CSV file into cells; any other file is one column
// per line. Empty lines are dropped, a UTF-8 BOM is removed.
func parseTable(name string, data []byte) (string, [][]string) {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	lines := nonEmptyLines(string(data))

	sep := ""
	if strings.EqualFold(filepath.Ext(name), ".csv") {
		sep = detectSeparator(lines)
	}
	if sep != "" {
		if rows, ok := readCSV(lines, sep); ok {
			return sep, rows
		}
	}

	rows := make([][]string, len(lines))
	for i, l := range lines {
		rows[i] = []string{l}
	}
	return "", rows
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, l := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// detectSeparator looks at the first lines. ';' and tab count when no line
// has more than one of them (content + optional caption) or all lines have
// the same number. A comma is common inside plain content, so it only counts
// when every line has the same, non-zero number of them.
func detectSeparator(lines []string) string {
	sample := lines[:min(len(lines), 20)]
	for _, sep := range []string{";", "\t", ","} {
		counts := make([]int, len(sample))
		maxCount, equal := 0, true
		for i, l := range sample {
			counts[i] = countOutsideQuotes(l, sep[0])
			maxCount = max(maxCount, counts[i])
			equal = equal && counts[i] == counts[0]
		}
		switch {
		case maxCount == 0:
			continue
		case equal, sep != "," && maxCount == 1:
			return sep
		}
	}
	return ""
}

func countOutsideQuotes(s string, sep byte) int {
	n, quoted := 0, false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			quoted = !quoted
		case sep:
			if !quoted {
				n++
			}
		}
	}
	return n
}

func readCSV(lines []string, sep string) ([][]string, bool) {
	r := csv.NewReader(strings.NewReader(strings.Join(lines, "\n")))
	r.Comma = rune(sep[0])
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	records, err := r.ReadAll()
	if err != nil {
		return nil, false
	}

	rows := make([][]string, 0, len(records))
	for _, rec := range records {
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}
		if len(rec) > 0 && rec[0] != "" {
			rows = append(rows, rec)
		}
	}
	return rows, true
}
