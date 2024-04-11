package views_controls

import (
	"fmt"
	"gfx"
	"math"

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
	schreiber := LiberationMonoBold(24)
	vr, vg, vb := f.vg.RGB()
	hr, hg, hb := f.hg.RGB()
	breite, höhe := f.GibGroesse()
	ra := höhe / 2

	var zeilenhöhe uint16
	// dieses Widget zeigt Treffer und Strafen an
	var anzKugeln uint8 = uint8(len(f.billard.GibKugeln()) - 1)
	tr, st := f.billard.GibTreffer(), math.Min(float64(anzKugeln), float64(f.billard.GibStrafpunkte()))
	zeilenhöhe = höhe / 2
	schreiber.SetzeSchriftgroesse(int(zeilenhöhe) * 3 / 5)
	d := (zeilenhöhe - uint16(schreiber.GibSchriftgroesse())) / 2
	var bBalken, xSBalken uint16 = breite - 2*ra - 5*uint16(schreiber.GibSchriftgroesse()), f.startX + 2*ra + 5*uint16(schreiber.GibSchriftgroesse())

	gfx.Stiftfarbe(hr, hg, hb)
	gfx.Vollkreis(f.startX+ra, f.startY+ra, ra)
	gfx.Vollrechteck(f.startX+ra, f.startY, xSBalken-f.startX-ra, höhe)
	gfx.Stiftfarbe(vr, vg, vb)

	schreiber.Schreibe(f.startX+2*ra+d, f.startY+d, "Treffer")
	schreiber.Schreibe(f.startX+2*ra+d, f.startY+zeilenhöhe+d, "Fouls")

	// zeichne beide Fortschritts-Balken
	gfx.Stiftfarbe(243, 186, 0) // Treffer gelb
	gfx.Vollrechteck(xSBalken, f.startY+1, bBalken*uint16(tr)/uint16(anzKugeln), zeilenhöhe-2)
	gfx.Stiftfarbe(219, 80, 0) // Fouls in Warnrot
	gfx.Vollrechteck(xSBalken, f.startY+zeilenhöhe+1, bBalken*uint16(st)/uint16(anzKugeln), zeilenhöhe-2)

	// Kreis links zeigt Treffer an
	schreiber.SetzeSchriftgroesse(int(ra * 4 / 3))
	gfx.Stiftfarbe(vr, vg, vb)
	gfx.Kreis(ra+f.startX, ra+f.startY, ra)
	var x, y uint16
	if tr > 9 {
		x = ra + f.startX - uint16(schreiber.GibSchriftgroesse())*2/5
		y = ra + f.startY - uint16(schreiber.GibSchriftgroesse())/2
	} else {
		x = ra + f.startX - uint16(schreiber.GibSchriftgroesse())/4
		y = ra + f.startY - uint16(schreiber.GibSchriftgroesse())/2
	}
	schreiber.Schreibe(x, y, fmt.Sprintf("%d", tr))
}
