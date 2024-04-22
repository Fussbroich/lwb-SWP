package views_controls

import (
	"strings"
	"unicode/utf8"
)

// Eine simple Box, die Text in Zeilen umbricht und anzeigt.
type textbox struct {
	text           string
	schriftgroesse int
	schreiber      *schreiber
	widget
}

func NewTextBox(t string, g int) *textbox {
	w := textbox{text: t, schriftgroesse: g,
		widget: *NewFenster()}
	w.schreiber = w.monoRegularSchreiber()
	w.schreiber.SetzeSchriftgroesse(g)
	return &w
}

func teileTextInZeilen(text string, nMax int) (zeilen []string) {
	// Brich Zeilen um, die länger als nMax Zeichen sind
	for _, line := range strings.Split(text, "\n") {
		if utf8.RuneCountInString(line) <= nMax {
			zeilen = append(zeilen, strings.TrimSpace(line))
		} else {
			// Teile lange Zeilen in Teile von maximal n Zeichen auf
			for len(line) > 0 {
				if utf8.RuneCountInString(line) > nMax {
					// Suche nach einem Leerzeichen oder Bindestrich, um die Zeile zu unterbrechen
					breakIndex := nMax
					for breakIndex > 0 && !isWhitespaceOrHyphen(line[breakIndex-1]) {
						// Überprüfe, ob das Zeichen ein  Rune ist
						if utf8.RuneStart(line[breakIndex-1]) {
							_, size := utf8.DecodeLastRuneInString(line[:breakIndex])
							breakIndex -= size
						} else {
							breakIndex--
						}
					}
					if breakIndex == 0 {
						// Kein Trennzeichen gefunden, daher verwende nMax Zeichen
						zeilen = append(zeilen, strings.TrimSpace(line[:nMax]))
						line = line[nMax:]
					} else {
						// Trennzeichen gefunden, daher verwende breakIndex
						zeilen = append(zeilen, strings.TrimSpace(line[:breakIndex]))
						line = line[breakIndex:]
					}
				} else {
					zeilen = append(zeilen, line)
					break
				}
			}
		}
	}
	return
}

// Hilfsfunktion, um zu überprüfen, ob ein Zeichen ein Leerzeichen oder Bindestrich ist
func isWhitespaceOrHyphen(char byte) bool {
	return char == ' ' || char == '-'
}

func (f *textbox) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	B, H := f.GibGroesse()
	f.stiftfarbe(f.vg)
	zMax, cMax := int(H)/f.schriftgroesse, int(B)/(2*f.schriftgroesse/3)
	for z, zeile := range teileTextInZeilen(f.text, cMax) {
		if z > zMax {
			break
		}
		f.schreiber.Schreibe(
			f.startX+f.eckra,
			f.startY+f.eckra+uint16(z*(f.schriftgroesse*6/5)),
			zeile)
	}
}
