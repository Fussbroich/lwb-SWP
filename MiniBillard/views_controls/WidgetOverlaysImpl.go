package views_controls

import (
	"gfx"

	"../fonts"
)

type text_overlay struct {
	text string
	widget
}

type infotext struct {
	text string
	widget
}

// TextOverlay zeigt den Hintergrund
func NewTextOverlay(startx, starty, stopx, stopy uint16, t string, hg, vg Farbe, tr uint8) *text_overlay {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &text_overlay{text: t, widget: fenster}
}

func (f *text_overlay) Zeichne() {
	f.widget.Zeichne()
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	font := fonts.LiberationMonoBoldItalic(int(f.stopY-f.startY) / 5)
	gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgröße())
	gfx.SchreibeFont((f.stopX-f.startX)/3, (f.stopY-f.startY)/4, f.text)
}

// InfoText ist immer Transparent
func NewInfoText(startx, starty, stopx, stopy uint16, t string, vg Farbe) *infotext {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy,
		hg: Weiß(), vg: vg, transparenz: 255}
	return &infotext{text: t, widget: fenster}
}

func (f *infotext) Zeichne() {
	f.widget.Zeichne()
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)

	_, höhe := f.GibGröße()
	font := fonts.LiberationMonoBoldItalic(int(höhe) * 3 / 5)
	d := (höhe - uint16(font.GibSchriftgröße())) / 2

	gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgröße())
	gfx.SchreibeFont(f.startX+d, f.startY+d, f.text)
}
