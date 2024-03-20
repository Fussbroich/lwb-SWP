package views

import (
	"gfx"

	"../welt"
)

type quizfenster struct {
	quiz  welt.Quiz
	frage Fenster
	as    [4]Fenster
	fenster
}

// TextOverlay zeigt den Hintergrund
func NewQuizFenster(quiz welt.Quiz, startx, starty, stopx, stopy uint16, hg, vg Farbe, ra uint16) *quizfenster {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: 0, eckradius: ra}
	return &quizfenster{quiz: quiz, fenster: fenster}
}

func (f *quizfenster) MausklickBei(mausX, mausY uint16) {
	for i, a := range f.as {
		if a.ImFenster(mausX, mausY) {
			println("Antwort:", i)
			f.quiz.Antwort(i)
		}
	}
}

func (f *quizfenster) Zeichne() {
	f.ZeichneRand()
	f.fenster.Zeichne()
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)

	breite, höhe := f.GibGröße()
	ra := f.eckradius
	breite, höhe = breite-2*ra, höhe-2*ra
	sx, sy := f.startX+ra, f.startY+ra
	var d uint16 = 3
	f.frage = NewTextFenster(sx, sy, sx+breite, sy+höhe*3/7-d,
		f.quiz.GibAktuelleFrage().GibFrage(), f.hg, f.vg, f.transparenz, 0)
	f.as = [4]Fenster{
		NewTextFenster(sx, sy+höhe*3/7, sx+breite/2-d, sy+höhe*5/7-d,
			f.quiz.GibAktuelleFrage().GibAntworten()[0], F(155, 155, 0), f.vg, f.transparenz, 0),
		NewTextFenster(sx+breite/2+d, sy+höhe*3/7, sx+breite, sy+höhe*5/7-d,
			f.quiz.GibAktuelleFrage().GibAntworten()[1], F(255, 255, 0), f.vg, f.transparenz, 0),
		NewTextFenster(sx, sy+höhe*5/7+d, sx+breite/2-d, sy+höhe,
			f.quiz.GibAktuelleFrage().GibAntworten()[2], F(0, 255, 255), f.vg, f.transparenz, 0),
		NewTextFenster(sx+breite/2+d, sy+höhe*5/7+d, sx+breite, sy+höhe,
			f.quiz.GibAktuelleFrage().GibAntworten()[3], F(255, 0, 255), f.vg, f.transparenz, 0)}
	f.frage.Zeichne()
	for _, af := range f.as {
		af.ZeichneRand()
		af.Zeichne()
	}
}
