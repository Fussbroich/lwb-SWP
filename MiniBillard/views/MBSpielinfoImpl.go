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

func NewMBSpielinfoFenster(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBSpielinfo {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, trans: tr}
	return &miniBSpielinfo{billard: billard, fenster: fenster}
}

func (f *miniBSpielinfo) Zeichne() {
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	//schreibe Stößezahl
	cr, cg, cb := f.vg.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	schriftgröße := (f.stopY - f.startY) / 4
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Stöße/%d Strafen", f.billard.GibStößeBisher(), f.billard.GibStrafpunkte()))
}
