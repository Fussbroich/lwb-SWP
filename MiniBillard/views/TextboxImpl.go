package views

import (
	"gfx"
)

type textbox struct {
	text string
	widget
}

func NewTextFenster(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8, ra uint16) *textbox {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
	return &textbox{text: t, widget: fenster}
}

func (f *textbox) Zeichne() {
	f.widget.Zeichne()
	//breite, höhe := f.GibGröße()
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	schriftgröße := 24
	gfx.SetzeFont(fp, schriftgröße)
	gfx.SchreibeFont(f.startX+f.eckradius, f.startY+f.eckradius, f.text)
}
