package views

import (
	"gfx"
)

type textfenster struct {
	text string
	fenster
}

func NewTextFenster(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8, ra uint16) *textfenster {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
	return &textfenster{text: t, fenster: fenster}
}

func (f *textfenster) Zeichne() {
	f.fenster.Zeichne()
	//breite, höhe := f.GibGröße()
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	schriftgröße := 24
	gfx.SetzeFont(fp, schriftgröße)
	gfx.SchreibeFont(f.startX+f.eckradius, f.startY+f.eckradius, f.text)
}
