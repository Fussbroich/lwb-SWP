package views

import (
	"gfx"
	"math"
	"unicode"
	"unicode/utf8"
)

// Eine simple Box, die Text in Zeilen umbricht und anzeigt.
type textbox struct {
	text string
	widget
}

func NewTextFenster(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8, ra uint16) *textbox {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
	return &textbox{text: t, widget: fenster}
}

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
	B, H := f.GibGröße()
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	// Schriftgröße automatisch anpasssen bzgl. Gesamtfläche der Box
	var sh, sb int
	sh = int(math.Min(24, math.Sqrt(float64(B*H)/float64(utf8.RuneCountInString(f.text))*12/7*5/6)))
	sb = 7 * sh / 12
	gfx.SetzeFont(fp, sh)
	zMax, cMax := int(H)/sh, int(B)/sb
	for z, zeile := range worteInZeilen(textInWorte(f.text), cMax) {
		if z > zMax {
			break
		}
		gfx.SchreibeFont(f.startX+f.eckradius, f.startY+f.eckradius+uint16(z*(sh*6/5)), zeile)
	}
}
