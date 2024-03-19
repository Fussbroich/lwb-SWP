package views

import (
	"gfx"

	"../welt"
)

type quizfenster struct {
	frage welt.QuizFrage
	fenster
}

// TextOverlay zeigt den Hintergrund
func NewQuizFenster(frage welt.QuizFrage, startx, starty, stopx, stopy uint16, hg, vg Farbe, ra uint16) *quizfenster {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: 0, eckradius: ra}
	return &quizfenster{frage: frage, fenster: fenster}
}

func (f *quizfenster) Zeichne() {
	f.ZeichneRand()
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationSerif-BoldItalic.ttf")
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)

	_, höhe := f.GibGröße()
	zeilenhöhe := höhe / 5
	schriftgröße := zeilenhöhe * 2 / 5
	d := (zeilenhöhe - schriftgröße) / 2
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY+d, f.frage.GibFrage())
	gfx.SchreibeFont(f.startX, f.startY+d+zeilenhöhe, f.frage.GibAntworten()[0])
	gfx.SchreibeFont(f.startX, f.startY+d+2*zeilenhöhe, f.frage.GibAntworten()[1])
	gfx.SchreibeFont(f.startX, f.startY+d+3*zeilenhöhe, f.frage.GibAntworten()[2])
	gfx.SchreibeFont(f.startX, f.startY+d+4*zeilenhöhe, f.frage.GibAntworten()[3])
}
