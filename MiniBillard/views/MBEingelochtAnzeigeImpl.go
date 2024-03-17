package views

import (
	"fmt"
	"gfx"

	"../welt"
)

type miniBEingelochte struct {
	billard welt.MiniBillardSpiel
	fenster
}

func NewMBEingelochteFenster(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBEingelochte {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &miniBEingelochte{billard: billard, fenster: fenster}
}

func (f *miniBEingelochte) Zeichne() {
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	//schreibe Stößezahl
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	höhe := f.stopY - f.startY
	schriftgröße := höhe / 3
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Eingelocht", len(f.billard.GibEingelochteKugeln())))
}
