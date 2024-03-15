package views

import "gfx"

type text_overlay struct {
	startX, startY uint16
	stopX, stopY   uint16
	text           string
	transparenz    uint8
	hg, vg         Farbe
}

type infotext struct {
	startX, startY uint16
	stopX, stopY   uint16
	text           string
	hg, vg         Farbe
}

func NewTextOverlay(startx, starty, stopx, stopy uint16, t string, tr uint8, hg, vg Farbe) *text_overlay {
	return &text_overlay{
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		text: t, transparenz: tr,
		hg: Weiß(), vg: vg}
}

func (f *text_overlay) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *text_overlay) GibGröße() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

func (f *text_overlay) Zeichne() {
	fp := fontDateipfad("LiberationSerif-BoldItalic.ttf")
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(f.transparenz)
	gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
	r, g, b = f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(0)
	schriftgröße := (f.stopY - f.startY) / 5
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont((f.stopX-f.startX)/3, (f.stopY-f.startY)/4, f.text)
}

func NewInfoText(startx, starty, stopx, stopy uint16, t string, vg Farbe) *infotext {
	return &infotext{
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		text: t,
		hg:   Weiß(), vg: vg}
}

func (f *infotext) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *infotext) GibGröße() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

func (f *infotext) Zeichne() {
	fp := fontDateipfad("LiberationSerif-BoldItalic.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	schriftgröße := (f.stopY - f.startY) / 5
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont((f.stopX-f.startX)/3, (f.stopY-f.startY)/4, f.text)
}
