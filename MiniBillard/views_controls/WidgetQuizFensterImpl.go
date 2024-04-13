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
	return &quizfenster{quiz: quiz, widget: *NewFenster()}
}

func (f *quizfenster) MausklickBei(mausX, mausY uint16) {
	if !f.IstAktiv() {
		return
	}
	for i, a := range f.antworten {
		if a.ImFenster(mausX, mausY) {
			f.quiz.GibAktuelleFrage().Gewaehlt(i)
		}
	}
}

func (f *quizfenster) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.ZeichneRand()
	f.widget.Zeichne()
	f.stiftfarbe(f.vg)

	breite, höhe := f.GibGroesse()
	ra := f.eckra
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

	f.frage.SetzeFarben(Fquiz(), Ftext())
	a0.SetzeFarben(FquizA0(), Ftext())
	a1.SetzeFarben(FquizA1(), Ftext())
	a2.SetzeFarben(FquizA2(), Ftext())
	a3.SetzeFarben(FquizA3(), Ftext())

	a0.SetzeTransparenz(f.trans)
	a1.SetzeTransparenz(f.trans)
	a2.SetzeTransparenz(f.trans)
	a3.SetzeTransparenz(f.trans)

	f.antworten = [4]Widget{a0, a1, a2, a3}
	f.frage.Zeichne()
	for _, af := range f.antworten {
		af.ZeichneRand()
		af.Zeichne()
	}
}
