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

	breite, höhe := f.GibGröße()
	ra := f.eckradius
	breite, höhe = breite-2*ra, höhe-2*ra
	sx, sy := f.startX+ra, f.startY+ra
	zeilenhöhe := höhe / 5
	schriftgröße := zeilenhöhe * 2 / 5
	//d := (zeilenhöhe - schriftgröße) / 2
	gfx.SetzeFont(fp, int(schriftgröße))
	var d uint16 = 3
	frage := NewTextFenster(sx, sy, sx+breite, sy+höhe*3/7-d,
		f.frage.GibFrage(), f.hg, f.vg, f.transparenz, 0)
	aA := NewTextFenster(sx, sy+höhe*3/7, sx+breite/2-d, sy+höhe*5/7-d,
		f.frage.GibAntworten()[0], F(155, 155, 0), f.vg, f.transparenz, ra/2)
	bA := NewTextFenster(sx+breite/2+d, sy+höhe*3/7, sx+breite, sy+höhe*5/7-d,
		f.frage.GibAntworten()[1], F(255, 255, 0), f.vg, f.transparenz, ra/2)
	cA := NewTextFenster(sx, sy+höhe*5/7+d, sx+breite/2-d, sy+höhe,
		f.frage.GibAntworten()[2], F(0, 255, 255), f.vg, f.transparenz, ra/2)
	dA := NewTextFenster(sx+breite/2+d, sy+höhe*5/7+d, sx+breite, sy+höhe,
		f.frage.GibAntworten()[3], F(255, 0, 255), f.vg, f.transparenz, ra/2)
	frage.Zeichne()
	aA.ZeichneRand()
	aA.Zeichne()
	bA.ZeichneRand()
	bA.Zeichne()
	cA.ZeichneRand()
	cA.Zeichne()
	dA.ZeichneRand()
	dA.Zeichne()
}
