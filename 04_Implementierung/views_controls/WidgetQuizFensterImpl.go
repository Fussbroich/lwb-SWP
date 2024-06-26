package views_controls

import (
	"../modelle"
)

type quizfenster struct {
	quiz      modelle.Quiz
	frage     Widget
	antworten [4]Widget
	richtig   func()
	falsch    func()
	widget
}

// TextOverlay zeigt den Hintergrund
func NewQuizFenster(quiz modelle.Quiz, richtig func(), falsch func()) *quizfenster {
	return &quizfenster{quiz: quiz, richtig: richtig, falsch: falsch, widget: *NewFenster()}
}

func (f *quizfenster) MausklickBei(mausX, mausY uint16) {
	if !f.IstAktiv() {
		return
	}
	var richtig bool
	for i, a := range f.antworten {
		if a.ImFenster(mausX, mausY) {
			if f.quiz.GibAktuelleFrage().IstRichtig(i) {
				richtig = true
				break
			}
		}
	}
	if richtig {
		f.richtig()
	} else {
		f.falsch()
	}
}

func (f *quizfenster) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	f.stiftfarbe(f.vg)

	breite, höhe := f.GibGroesse()
	ra := f.eckra
	breite, höhe = breite-2*ra, höhe-2*ra
	sx, sy := f.startX+ra, f.startY+ra
	var d uint16 = 3
	// zeichne eine Box
	// die Frage steht im oberen Drittel
	// der Rest ist wieder in 4 Boxen unterteilt (2x2), mit den möglichen Antworten
	f.frage = NewTextBox(f.quiz.GibAktuelleFrage().GibFrage(), Regular, 24)
	f.frage.SetzeKoordinaten(sx, sy, sx+breite, sy+höhe*2/8-d)
	a0 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[0], Regular, 22)
	a0.SetzeKoordinaten(sx, sy+höhe*2/8, sx+breite/2-d, sy+höhe*5/8-d) // oben links
	a1 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[1], Regular, 22)
	a1.SetzeKoordinaten(sx+breite/2+d, sy+höhe*2/8, sx+breite, sy+höhe*5/8-d) // oben rechts
	a2 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[2], Regular, 22)
	a2.SetzeKoordinaten(sx, sy+höhe*5/8+d, sx+breite/2-d, sy+höhe) // unten links
	a3 := NewTextBox(f.quiz.GibAktuelleFrage().GibAntworten()[3], Regular, 22)
	a3.SetzeKoordinaten(sx+breite/2+d, sy+höhe*5/8+d, sx+breite, sy+höhe) // unten rechts

	f.frage.SetzeFarben(Fquiz, Ftext)
	a0.SetzeFarben(FquizA0, Ftext)
	a1.SetzeFarben(FquizA1, Ftext)
	a2.SetzeFarben(FquizA2, Ftext)
	a3.SetzeFarben(FquizA3, Ftext)

	a0.SetzeTransparenz(f.trans)
	a1.SetzeTransparenz(f.trans)
	a2.SetzeTransparenz(f.trans)
	a3.SetzeTransparenz(f.trans)

	f.antworten = [4]Widget{a0, a1, a2, a3}
	f.frage.Zeichne()
	for _, af := range f.antworten {
		af.Zeichne()
	}
}
