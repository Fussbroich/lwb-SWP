package views_controls

import (
	"fmt"
	"gfx"

	"../modelle"
)

type MBSpielView interface {
	Widget
}

type miniBSpielfeld struct {
	billard       modelle.MiniBillardSpiel
	kugelZeichner *KugelZeichner
	widget
}

func NewMBSpieltisch(billard modelle.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8, ra uint16) *miniBSpielfeld {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
	kzeichner := &KugelZeichner{}
	return &miniBSpielfeld{billard: billard, kugelZeichner: kzeichner, widget: fenster}
}

func (f *miniBSpielfeld) zeichneDiamant(x, y, d uint16) {
	gfx.Volldreieck(x-d/2, y, x+d/2, y, x, y-d/2)
	gfx.Volldreieck(x-d/2, y, x+d/2, y, x, y+d/2)
}

func (f *miniBSpielfeld) Zeichne() {
	breite, höhe := f.GibGroesse()
	kS := f.billard.GibSpielkugel()
	ra := kS.GibRadius()
	schreiber := LiberationMonoRegular(24)
	// zeichne das Tuch
	f.widget.Zeichne()

	// zeichne Diamanten
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
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
		gfxVollKreis(f.startX, f.startY, t.GibPos(), ra*1.3, Schwarz())
	}
	// zeichne die Kugeln
	for _, k := range f.billard.GibAktiveKugeln() {
		f.kugelZeichner.ZeichneKugel(f.startX, f.startY, k)
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
		gfxBreiteLinie(f.startX, f.startY,
			pK, pK.Plus(f.billard.GibVStoss().Mal(ra)),
			4, farbe)
		// Schreibe den Wert der Stärke daneben
		gfx.Stiftfarbe(100, 100, 100)
		schreiber.SetzeSchriftgroesse(int(ra*0.67 + 0.5))
		pStärke := pK.Plus(f.billard.GibVStoss().Mal(ra * 3 / 4))
		schreiber.Schreibe(f.startX+uint16(pStärke.X()), f.startY+uint16(pStärke.Y()-2*ra), fmt.Sprintf("Stärke: %d", uint16(stärke+0.5)))
	}
	// debugging
	if !f.billard.Laeuft() {
		// Pause
		gfx.Stiftfarbe(100, 100, 100)
		schreiber.SetzeSchriftgroesse(int(f.billard.GibSpielkugel().GibRadius() + 0.5))
		schreiber.Schreibe(4*breite/5, f.startY+5, "Pause")
	} else if f.billard.IstZeitlupe() {
		// zeichne Geschwindigkeiten
		for _, k := range f.billard.GibAktiveKugeln() {
			if !k.GibV().IstNull() {
				gfxBreiteLinie(f.startX, f.startY,
					k.GibPos(), k.GibPos().Plus(k.GibV().Mal(k.GibRadius())),
					2, F(250, 175, 50))
			}
		}
		gfx.Stiftfarbe(100, 100, 100)
		schreiber.SetzeSchriftgroesse(int(f.billard.GibSpielkugel().GibRadius() + 0.5))
		schreiber.Schreibe(4*breite/5, f.startY+5, "Zeitlupe")
	}
}
