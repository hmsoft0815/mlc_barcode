package gui

import (
	"reflect"
	"testing"
)

func TestParseTable(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		data    string
		wantSep string
		want    [][]string
	}{
		{"txt never splits", "codes.txt", "https://a.de/x,y\nWIFI:S:net;T:WPA;;\n", "",
			[][]string{{"https://a.de/x,y"}, {"WIFI:S:net;T:WPA;;"}}},
		{"csv with only text", "codes.csv", "ART-1\nART-2\r\n\nART-3\n", "",
			[][]string{{"ART-1"}, {"ART-2"}, {"ART-3"}}},
		{"csv content;caption, caption optional", "codes.csv", "ART-1;Schraube M4\nART-2\n", ";",
			[][]string{{"ART-1", "Schraube M4"}, {"ART-2"}}},
		{"csv quoted content with separator", "codes.csv", "\"WIFI:S:net;T:WPA;;\";Gäste-WLAN\n", ";",
			[][]string{{"WIFI:S:net;T:WPA;;", "Gäste-WLAN"}}},
		{"single commas stay content", "codes.csv", "Hallo, Welt\nZweite Zeile\n", "",
			[][]string{{"Hallo, Welt"}, {"Zweite Zeile"}}},
		{"consistent commas split", "codes.csv", "A1,Erster\nA2,Zweiter\n", ",",
			[][]string{{"A1", "Erster"}, {"A2", "Zweiter"}}},
		{"tab separated", "codes.csv", "A1\tErster\nA2\tZweiter\n", "\t",
			[][]string{{"A1", "Erster"}, {"A2", "Zweiter"}}},
		{"excel bom", "codes.csv", "\xef\xbb\xbfA1;Erster\n", ";",
			[][]string{{"A1", "Erster"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sep, rows := parseTable(tt.file, []byte(tt.data))
			if sep != tt.wantSep {
				t.Errorf("separator %q, want %q", sep, tt.wantSep)
			}
			if !reflect.DeepEqual(rows, tt.want) {
				t.Errorf("rows %q, want %q", rows, tt.want)
			}
		})
	}
}
