package views_controls

import (
	"fmt"

	"../modelle"
)

type kugelZeichner struct {
	kugel   modelle.MBKugel
	english bool
	widget
}

var (
	standardPoolPalette = [16]Farbe{
		F(252, 253, 242), // weiß
		F(255, 201, 78),  // gelb
		F(34, 88, 175),   // blau
		F(249, 73, 68),   // hellrot
		F(84, 73, 149),   // violett
		F(255, 139, 33),  // orange
		F(47, 159, 52),   // grün
		F(155, 53, 30),   // dunkelrot
		F(48, 49, 54),    // schwarz
		F(255, 201, 78),  // gelb
		F(34, 88, 175),   // blau
		F(249, 73, 68),   // hellrot
		F(84, 73, 149),   // violett
		F(255, 139, 33),  // orange
		F(47, 159, 52),   // grün
		F(155, 53, 30)}   // dunkelrot

	englishPoolPalette = [16]Farbe{
		F(252, 253, 242), // weiß
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(48, 49, 54),    // schwarz
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78)}  // gelb

)

func (w *kugelZeichner) GibKugelPalette() *[16]Farbe {
	if w.english {
		return &englishPoolPalette
	}
	return &standardPoolPalette
}

func (w *kugelZeichner) SetzeEnglish() {
	w.english = true
}

func (w *kugelZeichner) SetzeKugel(k modelle.MBKugel) {
	w.kugel = k
}

func (w *kugelZeichner) Zeichne() {
	k := w.kugel
	schreiber := LiberationMonoBold(int(k.GibRadius()) - 3)
	w.VollKreis(k.GibPos(), k.GibRadius(), F(48, 49, 54))
	w.VollKreis(k.GibPos(), k.GibRadius()-1, F(252, 253, 242))
	c := w.GibKugelPalette()[k.GibWert()]
	if k.GibWert() <= 8 || w.english {
		w.VollKreis(k.GibPos(), k.GibRadius()-1, c)
	} else {
		w.Stiftfarbe(c)
		w.VollRechteckGFX(
			uint16(k.GibPos().X()-k.GibRadius()*0.75+0.5),
			uint16(k.GibPos().Y()-k.GibRadius()*0.6+0.5),
			uint16(2*0.75*k.GibRadius()+0.5),
			uint16(2*0.6*k.GibRadius()+0.5))
		w.VollKreissektor(k.GibPos(), k.GibRadius()-1, 325, 35, c)
		w.VollKreissektor(k.GibPos(), k.GibRadius()-1, 145, 215, c)
	}
	// Die weiße erhält keine Nummer.
	if k.GibWert() != 0 && !w.english {
		w.VollKreis(k.GibPos(), (k.GibRadius()-1)/2, F(252, 253, 242))
		w.Stiftfarbe(Schwarz())
		if k.GibWert() < 10 {
			schreiber.Schreibe(
				w.startX-uint16(schreiber.GibSchriftgroesse())/4+uint16(k.GibPos().X()+0.5),
				w.startY-uint16(schreiber.GibSchriftgroesse())/2+uint16(k.GibPos().Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		} else {
			schreiber.Schreibe(
				w.startX-uint16(schreiber.GibSchriftgroesse())/2+uint16(k.GibPos().X()+0.5),
				w.startY-uint16(schreiber.GibSchriftgroesse())/2+uint16(k.GibPos().Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		}
	}
}
