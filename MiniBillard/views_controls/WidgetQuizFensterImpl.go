package views_controls

import (
	"gfx"

	"../modelle"
)

type quizfenster struct {
	quiz      modelle.Quiz
	frage     Widget
	antworten [4]Widget
	widget
}

// TextOverlay zeigt den Hintergrund
func NewQuizFenster(quiz modelle.Quiz, startx, starty, stopx, stopy uint16, hg, vg Farbe, ra uint16) *quizfenster {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: 0, eckradius: ra}
	return &quizfenster{quiz: quiz, widget: fenster}
}

func (f *quizfenster) MausklickBei(mausX, mausY uint16) {
	for i, a := range f.antworten {
		if a.ImFenster(mausX, mausY) {
			f.quiz.GibAktuelleFrage().Gewaehlt(i)
		}
	}
}

func (f *quizfenster) Zeichne() {
	f.ZeichneRand()
	f.widget.Zeichne()
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)

	breite, höhe := f.GibGroesse()
	ra := f.eckradius
	breite, höhe = breite-2*ra, höhe-2*ra
	sx, sy := f.startX+ra, f.startY+ra
	var d uint16 = 3
	f.frage = NewTextFenster(sx, sy, sx+breite, sy+höhe*3/7-d,
		f.quiz.GibAktuelleFrage().GibFrage(), f.hg, f.vg, f.transparenz, 0)
	f.antworten = [4]Widget{
		NewTextFenster(sx, sy+höhe*3/7, sx+breite/2-d, sy+höhe*5/7-d,
			f.quiz.GibAktuelleFrage().GibAntworten()[0], F(155, 155, 0), f.vg, f.transparenz, 0),
		NewTextFenster(sx+breite/2+d, sy+höhe*3/7, sx+breite, sy+höhe*5/7-d,
			f.quiz.GibAktuelleFrage().GibAntworten()[1], F(255, 255, 0), f.vg, f.transparenz, 0),
		NewTextFenster(sx, sy+höhe*5/7+d, sx+breite/2-d, sy+höhe,
			f.quiz.GibAktuelleFrage().GibAntworten()[2], F(0, 255, 255), f.vg, f.transparenz, 0),
		NewTextFenster(sx+breite/2+d, sy+höhe*5/7+d, sx+breite, sy+höhe,
			f.quiz.GibAktuelleFrage().GibAntworten()[3], F(255, 0, 255), f.vg, f.transparenz, 0)}
	f.frage.Zeichne()
	for _, af := range f.antworten {
		af.ZeichneRand()
		af.Zeichne()
	}
}
