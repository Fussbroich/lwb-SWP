package views

import "gfx"

type text_overlay struct {
	text string
	fenster
}

type infotext struct {
	text string
	fenster
}

func NewTextOverlay(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8) *text_overlay {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &text_overlay{text: t, fenster: fenster}
}

func (f *text_overlay) Zeichne() {
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationSerif-BoldItalic.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	schriftgröße := (f.stopY - f.startY) / 5
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont((f.stopX-f.startX)/3, (f.stopY-f.startY)/4, f.text)
}

func NewInfoText(startx, starty, stopx, stopy uint16, t string, vg Farbe) *infotext {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy,
		hg: Weiß(), vg: vg, transparenz: 255}
	return &infotext{text: t, fenster: fenster}
}

func (f *infotext) Zeichne() {
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationSerif-BoldItalic.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	schriftgröße := (f.stopY - f.startY) / 5
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY, f.text)
}
