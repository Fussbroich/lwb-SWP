package views

import (
	"fmt"
	"gfx"

	"../welt"
)

type miniBSpielinfo struct {
	billard welt.MiniBillardSpiel
	fenster
}

func NewMBPunkteAnzeiger(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBSpielinfo {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &miniBSpielinfo{billard: billard, fenster: fenster}
}

func (f *miniBSpielinfo) Zeichne() {
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	vr, vg, vb := f.vg.RGB()
	hr, hg, hb := f.hg.RGB()
	breite, höhe := f.GibGröße()
	ra := höhe / 2

	var zeilenhöhe uint16
	var schriftgröße int

	// Display rechts zeigt Treffer und Strafen an
	tr, st := f.billard.GibTreffer(), f.billard.GibStrafpunkte()
	zeilenhöhe = höhe / 2
	schriftgröße = int(zeilenhöhe) * 3 / 5
	d := (zeilenhöhe - uint16(schriftgröße)) / 2

	gfx.Stiftfarbe(hr, hg, hb)
	gfx.Vollkreis(f.startX+ra, f.startY+ra, ra)
	gfx.Vollrechteck(f.startX+ra, f.startY, breite-ra, höhe)
	gfx.Stiftfarbe(vr, vg, vb)
	gfx.SetzeFont(fp, schriftgröße)
	gfx.SchreibeFont(f.startX+2*ra+d, f.startY+d, "Treffer")
	gfx.SchreibeFont(f.startX+2*ra+d, f.startY+zeilenhöhe+d, "Fouls")

	// zeichne Fortschritts-Balken
	var bBalken, xSBalken uint16 = breite - 2*ra - 5*uint16(schriftgröße), f.startX + 2*ra + 5*uint16(schriftgröße)
	var anzKugeln uint8 = uint8(len(f.billard.GibKugeln()) - 1)

	gfx.Stiftfarbe(34, 88, 175)
	gfx.Vollrechteck(xSBalken, f.startY+1, bBalken*uint16(tr)/uint16(anzKugeln), zeilenhöhe-2)
	gfx.Stiftfarbe(255, 201, 78)
	gfx.Vollrechteck(xSBalken, f.startY+zeilenhöhe+1, bBalken*uint16(st)/uint16(anzKugeln), zeilenhöhe-2)

	// Kreis links zeigt Treffer an
	schriftgröße = int(ra * 4 / 3)
	gfx.Stiftfarbe(vr, vg, vb)
	gfx.Kreis(ra+f.startX, ra+f.startY, ra)
	gfx.SetzeFont(fp, schriftgröße)
	var x, y uint16
	if tr > 9 {
		x = ra + f.startX - uint16(schriftgröße)*2/5
		y = ra + f.startY - uint16(schriftgröße)/2
	} else {
		x = ra + f.startX - uint16(schriftgröße)/4
		y = ra + f.startY - uint16(schriftgröße)/2
	}
	gfx.SchreibeFont(x, y, fmt.Sprintf("%d", tr))
}
