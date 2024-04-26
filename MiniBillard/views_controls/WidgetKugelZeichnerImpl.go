package views_controls

import (
	"fmt"

	"../modelle"
)

type kugelZeichner struct {
	kugel     modelle.MBKugel
	schreiber *schreiber
	widget
}

// Kugelzeichner übernimmt die Koordinaten eines anderen Widgets
func newKugelZeichnerIn(w *widget) *kugelZeichner {
	k := kugelZeichner{widget: *w}
	k.schreiber = w.monoBoldSchreiber()
	return &k
}

func (w *kugelZeichner) SetzeKugel(k modelle.MBKugel) {
	w.kugel = k
}

func (w *kugelZeichner) Zeichne() {
	if !w.IstAktiv() {
		return
	}
	k := w.kugel
	w.schreiber.SetzeSchriftgroesse(int(k.GibRadius()) - 3)
	c := gibKugelFarbe(k.GibWert())
	if k.GibWert() <= 8 || w.schlicht {
		w.vollKreis(k.GibPos(), k.GibRadius()-1, c)
	} else {
		w.vollKreis(k.GibPos(), k.GibRadius()-1, gibKugelFarbe(0))
		w.stiftfarbe(c)
		w.vollRechteckGFX(
			uint16(k.GibPos().X()-k.GibRadius()*0.75+0.5),
			uint16(k.GibPos().Y()-k.GibRadius()*0.6+0.5),
			uint16(2*0.75*k.GibRadius()+0.5),
			uint16(2*0.6*k.GibRadius()+0.5))
		w.vollKreissektor(k.GibPos(), k.GibRadius()-1, 325, 35, c)
		w.vollKreissektor(k.GibPos(), k.GibRadius()-1, 145, 215, c)
	}
	// Die weiße erhält keine Nummer.
	if k.GibWert() != 0 && !w.schlicht {
		w.vollKreis(k.GibPos(), (k.GibRadius()-1)/2, gibKugelFarbe(0))
		w.stiftfarbe(Schwarz())
		if k.GibWert() < 10 {
			w.schreiber.Schreibe(
				w.startX-uint16(w.schreiber.GibSchriftgroesse())/4+uint16(k.GibPos().X()+0.5),
				w.startY-uint16(w.schreiber.GibSchriftgroesse())/2+uint16(k.GibPos().Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		} else {
			w.schreiber.Schreibe(
				w.startX-uint16(w.schreiber.GibSchriftgroesse())/2+uint16(k.GibPos().X()+0.5),
				w.startY-uint16(w.schreiber.GibSchriftgroesse())/2+uint16(k.GibPos().Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		}
	}
}
