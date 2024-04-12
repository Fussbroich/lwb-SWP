package views_controls

import (
	"math"
	"unicode"
	"unicode/utf8"
)

// Eine simple Box, die Text in Zeilen umbricht und anzeigt.
type textbox struct {
	text string
	widget
}

func NewTextBox(t string) *textbox { return &textbox{text: t, widget: widget{}} }

func worteInZeilen(worte []string, lMax int) (zeilen []string) {
	var zeile string
	// Baue Zeilen aus Worten.
	for _, wort := range worte {
		if utf8.RuneCountInString(wort) > lMax {
			zeilen = append(zeilen, wort)
			continue
		}
		if (utf8.RuneCountInString(zeile) + utf8.RuneCountInString(wort) + 1) <= lMax {
			if zeile != "" {
				zeile += " "
			}
			zeile += wort
		} else {
			zeilen = append(zeilen, zeile)
			zeile = wort
		}
	}
	if zeile != "" {
		zeilen = append(zeilen, zeile)
	}
	return
}

func textInWorte(text string) (worte []string) {
	var wort string
	// Zerlege den Text in Worte.
	for _, char := range text {
		if unicode.IsSpace(char) {
			if wort != "" {
				worte = append(worte, wort)
				wort = ""
			}
		} else {
			wort += string(char)
		}
	}
	if wort != "" {
		worte = append(worte, wort)
	}
	return
}

func (f *textbox) Zeichne() {
	f.widget.Zeichne()
	B, H := f.GibGroesse()
	f.Stiftfarbe(f.vg)
	// Schriftgroesse automatisch anpasssen bzgl. Gesamtfläche der Box
	schreiber := f.LiberationMonoRegularSchreiber()
	schreiber.SetzeSchriftgroesse(
		int(math.Min(
			24,
			math.Sqrt(float64(B*H)/float64(utf8.RuneCountInString(f.text))*12/7*5/6))))

	zMax, cMax := int(H)/schreiber.GibSchriftgroesse(), int(B)/(7*schreiber.GibSchriftgroesse()/12)
	for z, zeile := range worteInZeilen(textInWorte(f.text), cMax) {
		if z > zMax {
			break
		}
		schreiber.Schreibe(f.startX+f.eckradius, f.startY+f.eckradius+uint16(z*(schreiber.GibSchriftgroesse()*6/5)), zeile)
	}
}
