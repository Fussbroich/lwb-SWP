package views_controls

import (
	"strings"
	"unicode/utf8"
)

// Eine simple Box, die Text in Zeilen umbricht und anzeigt.
type textbox struct {
	text string
	widget
}

func NewTextBox(t string) *textbox { return &textbox{text: t, widget: *NewFenster()} }

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
	schreiber := f.monoRegularSchreiber()
	schreiber.SetzeSchriftgroesse(24)
	// Schriftgroesse automatisch anpasssen bzgl. Gesamtfläche der Box
	//	int(math.Min(24,
	//		math.Sqrt(float64(B*H)/float64(utf8.RuneCountInString(f.text))*12/7*5/6))))

	zMax, cMax := int(H)/schreiber.GibSchriftgroesse(), int(B)/(7*schreiber.GibSchriftgroesse()/12)
	for z, zeile := range teileTextInZeilen(f.text, cMax) {
		if z > zMax {
			break
		}
		schreiber.Schreibe(f.startX+f.eckra, f.startY+f.eckra+uint16(z*(schreiber.GibSchriftgroesse()*6/5)), zeile)
	}
}
