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
	cr, cg, cb := f.vg.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	var x, y uint16
	var schriftgröße int

	// Kreis links zeigt Treffer an
	ra := (f.stopY - f.startY) / 2
	schriftgröße = int(ra * 4 / 3)
	gfx.Kreis(ra+f.startX, ra+f.startY, ra)
	gfx.SetzeFont(fp, schriftgröße)
	tr := f.billard.GibTreffer()
	if tr > 9 {
		x = ra + f.startX - uint16(schriftgröße)*2/5
		y = ra + f.startY - uint16(schriftgröße)/2
	} else {
		x = ra + f.startX - uint16(schriftgröße)/4
		y = ra + f.startY - uint16(schriftgröße)/2
	}
	gfx.SchreibeFont(x, y, fmt.Sprintf("%d", tr))

	// Display rechts zeigt Treffer und Strafen an
	//st := f.billard.GibStrafpunkte()
	//fp = fontDateipfad("LiberationMono-Regular.ttf")
	x, y = f.startX+2*ra, f.startY
	zeilenhöhe := (f.stopY - f.startY) / 2
	d := zeilenhöhe / 10
	schriftgröße = int(zeilenhöhe) * 3 / 5
	gfx.Rechteck(x, y, 6*uint16(schriftgröße), 2*zeilenhöhe)
	gfx.SetzeFont(fp, schriftgröße)
	gfx.SchreibeFont(x+d, y+d, "Treffer")
	gfx.SchreibeFont(x+d, y+zeilenhöhe+d, "Fouls")

	//TODO

}
