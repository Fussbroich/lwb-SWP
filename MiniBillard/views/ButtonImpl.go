package views

import "gfx"

type button struct {
	text string
	fenster
}

// Buttons haben einen Text in der Mitte
func NewButton(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8, ra uint16) *button {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy,
		hg: hg, vg: vg, transparenz: 0, eckradius: ra}
	return &button{text: t, fenster: fenster}
}

func (f *button) Zeichne() {
	f.ZeichneRand()
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)

	breite, höhe := f.GibGröße()
	schriftgröße := int(höhe) * 3 / 5
	d := (höhe - uint16(schriftgröße)) / 2

	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX+(breite/2)-uint16(len(f.text)*schriftgröße*7/24), f.startY+d, f.text)
}
