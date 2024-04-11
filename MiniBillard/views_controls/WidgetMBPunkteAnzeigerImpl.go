package views_controls

import (
	"fmt"
	"gfx"
	"math"

	"../fonts"

	"../modelle"
)

type miniBSpielinfo struct {
	billard modelle.MiniBillardSpiel
	widget
}

func NewMBPunkteAnzeiger(billard modelle.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBSpielinfo {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &miniBSpielinfo{billard: billard, widget: fenster}
}

func (f *miniBSpielinfo) Zeichne() {
	f.widget.Zeichne()
	font := fonts.LiberationMonoBold(24)
	vr, vg, vb := f.vg.RGB()
	hr, hg, hb := f.hg.RGB()
	breite, höhe := f.GibGröße()
	ra := höhe / 2

	var zeilenhöhe uint16
	// dieses Widget zeigt Treffer und Strafen an
	var anzKugeln uint8 = uint8(len(f.billard.GibKugeln()) - 1)
	tr, st := f.billard.GibTreffer(), math.Min(float64(anzKugeln), float64(f.billard.GibStrafpunkte()))
	zeilenhöhe = höhe / 2
	font.SetzeSchriftgröße(int(zeilenhöhe) * 3 / 5)
	d := (zeilenhöhe - uint16(font.GibSchriftgröße())) / 2
	var bBalken, xSBalken uint16 = breite - 2*ra - 5*uint16(font.GibSchriftgröße()), f.startX + 2*ra + 5*uint16(font.GibSchriftgröße())

	gfx.Stiftfarbe(hr, hg, hb)
	gfx.Vollkreis(f.startX+ra, f.startY+ra, ra)
	gfx.Vollrechteck(f.startX+ra, f.startY, xSBalken-f.startX-ra, höhe)
	gfx.Stiftfarbe(vr, vg, vb)
	gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgröße())
	gfx.SchreibeFont(f.startX+2*ra+d, f.startY+d, "Treffer")
	gfx.SchreibeFont(f.startX+2*ra+d, f.startY+zeilenhöhe+d, "Fouls")

	// zeichne beide Fortschritts-Balken
	gfx.Stiftfarbe(243, 186, 0) // Treffer gelb
	gfx.Vollrechteck(xSBalken, f.startY+1, bBalken*uint16(tr)/uint16(anzKugeln), zeilenhöhe-2)
	gfx.Stiftfarbe(219, 80, 0) // Fouls in Warnrot
	gfx.Vollrechteck(xSBalken, f.startY+zeilenhöhe+1, bBalken*uint16(st)/uint16(anzKugeln), zeilenhöhe-2)

	// Kreis links zeigt Treffer an
	font.SetzeSchriftgröße(int(ra * 4 / 3))
	gfx.Stiftfarbe(vr, vg, vb)
	gfx.Kreis(ra+f.startX, ra+f.startY, ra)
	gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgröße())
	var x, y uint16
	if tr > 9 {
		x = ra + f.startX - uint16(font.GibSchriftgröße())*2/5
		y = ra + f.startY - uint16(font.GibSchriftgröße())/2
	} else {
		x = ra + f.startX - uint16(font.GibSchriftgröße())/4
		y = ra + f.startY - uint16(font.GibSchriftgröße())/2
	}
	gfx.SchreibeFont(x, y, fmt.Sprintf("%d", tr))
}
