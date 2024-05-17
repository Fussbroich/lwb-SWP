package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
	"../modelle"
)

type miniBSpielfeld struct {
	billard       modelle.MiniBillardSpiel
	kugelZeichner *kugelZeichner
	tpsInfo       Widget
	schreiber     *schreiber
	widget
}

func NewMBSpieltisch(billard modelle.MiniBillardSpiel) *miniBSpielfeld {
	f := miniBSpielfeld{billard: billard, widget: *NewFenster()}
	f.schreiber = f.newSchreiber(Regular)
	f.tpsInfo = NewInfoText(func() string { return fmt.Sprintf("%04d TPS", billard.GetTicksPS()*10/10) })
	return &f
}

// stoßen
func (f *miniBSpielfeld) MausklickBei(mausX, mausY uint16) {
	if !f.IstAktiv() {
		return
	}
	if !f.billard.Laeuft() || !f.billard.IstStillstand() {
		return
	}
	f.billard.Stosse()
}

// zielen
func (f *miniBSpielfeld) MausBei(mausX, mausY uint16) {
	if !f.IstAktiv() {
		return
	}
	if !f.billard.Laeuft() || !f.billard.IstStillstand() {
		return
	}
	xs, ys := f.GibStartkoordinaten()
	f.billard.SetzeStossRichtung(hilf.V2(float64(mausX), float64(mausY)).
		Minus(f.billard.GibSpielkugel().GibPos()).
		Minus(hilf.V2(float64(xs), float64(ys))))
}

// stärker
func (f *miniBSpielfeld) MausScrolltHoch() {
	if !f.IstAktiv() {
		return
	}
	if !f.billard.Laeuft() || !f.billard.IstStillstand() {
		return
	}
	f.billard.SetzeStosskraft(f.billard.GibVStoss().Betrag() + 1)
}

// schwächer
func (f *miniBSpielfeld) MausScrolltRunter() {
	if !f.IstAktiv() {
		return
	}
	if !f.billard.Laeuft() || !f.billard.IstStillstand() {
		return
	}
	f.billard.SetzeStosskraft(f.billard.GibVStoss().Betrag() - 1)
}

func (f *miniBSpielfeld) zeichneDiamant(x, y, d uint16) {
	gfx.Volldreieck(x-d/2, y, x+d/2, y, x, y-d/2)
	gfx.Volldreieck(x-d/2, y, x+d/2, y, x, y+d/2)
}

func (f *miniBSpielfeld) ZeichneLayout() {
	if !f.IstAktiv() {
		return
	}
	f.widget.ZeichneLayout()
	// Ticks je Sekunde
	breite, höhe := f.GibGroesse()
	f.tpsInfo.SetzeKoordinaten(f.startX+breite/20, f.startY, f.startX+breite/10, f.startY+höhe/10)
	f.tpsInfo.SetzeFarben(Fbillardtuch, Frot)
	f.tpsInfo.Zeichne()
}

func (f *miniBSpielfeld) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	if f.kugelZeichner == nil {
		f.kugelZeichner = newKugelZeichnerIn(&f.widget)
	}
	// zeichne das Tuch
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()
	kS := f.billard.GibSpielkugel()
	ra := kS.GibRadius()

	// zeichne Diamanten
	f.stiftfarbe(f.vg)
	for _, i := range []uint16{1, 2, 3, 5, 6, 7} {
		f.zeichneDiamant(f.startX+i*breite/8, f.startY-uint16(ra+0.5), uint16(ra/3+0.5))
		f.zeichneDiamant(f.startX+i*breite/8, f.startY+höhe+uint16(ra+0.5), uint16(ra/3+0.5))
	}
	for _, i := range []uint16{1, 2, 3} {
		f.zeichneDiamant(f.startX-uint16(ra+0.5), f.startY+i*höhe/4, uint16(ra/3+0.5))
		f.zeichneDiamant(f.startX+breite+uint16(ra+0.5), f.startY+i*höhe/4, uint16(ra/3+0.5))
	}
	// zeichne die Taschen
	for _, t := range f.billard.GibTaschen() {
		f.vollKreis(t.GibPos(), ra*1.3, Schwarz())
	}

	// zeichne die Kugeln
	for _, k := range f.billard.GibAktiveKugeln() {
		f.kugelZeichner.SetzeKugel(k)
		f.kugelZeichner.Zeichne()
	}
	if f.billard.IstStillstand() && !f.billard.GibSpielkugel().IstEingelocht() {
		pK := kS.GibPos()
		// zeichne die Stoßrichtung und -stärke bezogen auf Kugelradien
		stärke := f.billard.GibVStoss().Betrag()
		var farbe Farbe
		if stärke < 7 {
			farbe = F(47, 159, 52)
		} else if stärke > 9 {
			farbe = F(249, 73, 68)
		} else {
			farbe = F(250, 175, 50)
		}
		f.breiteLinie(pK, pK.Plus(f.billard.GibVStoss().Mal(ra)), ra/4, farbe)
		// Schreibe den Wert der Stärke daneben
		f.stiftfarbe(F(100, 100, 100))
		f.schreiber.SetzeSchriftgroesse(int(ra*0.67 + 0.5))
		pStärke := pK.Plus(f.billard.GibVStoss().Mal(ra * 3 / 4))
		f.schreiber.Schreibe(f.startX+uint16(pStärke.X()), f.startY+uint16(pStärke.Y()-2*ra), fmt.Sprintf("Stärke: %d", uint16(stärke+0.5)))
	}
}
