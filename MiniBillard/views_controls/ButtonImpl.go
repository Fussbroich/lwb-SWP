package views_controls

import (
	"gfx"

	"../fonts"
)

type button struct {
	text string
	widget
}

// Buttons haben einen Text in der Mitte
func NewButton(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8, ra uint16) *button {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy,
		hg: hg, vg: vg, transparenz: 0, eckradius: ra}
	return &button{text: t, widget: fenster}
}

func (f *button) Zeichne() {
	f.ZeichneRand()
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()

	font := fonts.LiberationMonoRegular(int(höhe) * 3 / 5)
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)

	d := (höhe - uint16(font.GibSchriftgroesse())) / 2

	gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgroesse())
	gfx.SchreibeFont(f.startX+(breite/2)-uint16(len(f.text)*font.GibSchriftgroesse()*7/24), f.startY+d, f.text)
}
