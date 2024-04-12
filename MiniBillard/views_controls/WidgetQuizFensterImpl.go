package views_controls

import (
	"../modelle"
)

type quizfenster struct {
	quiz      modelle.Quiz
	frage     Widget
	antworten [4]Widget
	widget
}

// TextOverlay zeigt den Hintergrund
func NewQuizFenster(quiz modelle.Quiz) *quizfenster {
	return &quizfenster{quiz: quiz, widget: widget{}}
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
	f.Stiftfarbe(f.vg)

	breite, höhe := f.GibGroesse()
	ra := f.eckradius
	breite, höhe = breite-2*ra, höhe-2*ra
	sx, sy := f.startX+ra, f.startY+ra
	var d uint16 = 3
	f.frage = NewTextBox(f.quiz.GibAktuelleFrage().GibFrage())
	f.frage.SetzeKoordinaten(sx, sy, sx+breite, sy+höhe*3/7-d)
	a0 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[0])
	a0.SetzeKoordinaten(sx, sy+höhe*3/7, sx+breite/2-d, sy+höhe*5/7-d)
	a1 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[1])
	a1.SetzeKoordinaten(sx+breite/2+d, sy+höhe*3/7, sx+breite, sy+höhe*5/7-d)
	a2 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[2])
	a2.SetzeKoordinaten(sx, sy+höhe*5/7+d, sx+breite/2-d, sy+höhe)
	a3 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[3])
	a3.SetzeKoordinaten(sx+breite/2+d, sy+höhe*5/7+d, sx+breite, sy+höhe)

	f.frage.SetzeFarben(f.hg, f.vg)
	a0.SetzeFarben(F(155, 155, 0), f.vg)
	a1.SetzeFarben(F(255, 255, 0), f.vg)
	a2.SetzeFarben(F(0, 255, 255), f.vg)
	a3.SetzeFarben(F(255, 0, 255), f.vg)

	a0.SetzeTransparenz(f.transparenz)
	a1.SetzeTransparenz(f.transparenz)
	a2.SetzeTransparenz(f.transparenz)
	a3.SetzeTransparenz(f.transparenz)

	f.antworten = [4]Widget{a0, a1, a2, a3}
	f.frage.Zeichne()
	for _, af := range f.antworten {
		af.ZeichneRand()
		af.Zeichne()
	}
}
