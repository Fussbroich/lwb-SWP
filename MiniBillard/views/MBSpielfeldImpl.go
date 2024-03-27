package views

import (
	"fmt"
	"gfx"

	"../welt"
)

type MBSpielView interface {
	Widget
}

type miniBSpielfeld struct {
	billard welt.MiniBillardSpiel
	widget
}

func NewMBSpieltisch(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8, ra uint16) *miniBSpielfeld {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
	return &miniBSpielfeld{billard: billard, widget: fenster}
}

func (f *miniBSpielfeld) Zeichne() {
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	breite, _ := f.GibGröße()
	kS := f.billard.GibStoßkugel()
	ra := kS.GibRadius()
	// zeichne das Tuch
	f.widget.Zeichne()

	// zeichne die Taschen
	for _, t := range f.billard.GibTaschen() {
		gfxVollKreis(f.startX, f.startY, t.GibPos(), ra*1.3, Schwarz())
	}
	// zeichne die Kugeln
	for _, k := range f.billard.GibAktiveKugeln() {
		zeichneKugel(f.startX, f.startY, k.GibPos(), k)
	}
	if f.billard.IstStillstand() && !f.billard.GibStoßkugel().IstEingelocht() {
		pK := kS.GibPos()
		// Zeichne Peillinie bis zum Rand
		gfx.Stiftfarbe(100, 100, 100)
		zielP := pK.Plus(f.billard.GibVStoß().Normiert().Mal(float64(2 * breite)))
		gfx.Linie(f.startX+uint16(pK.X()), f.startY+uint16(pK.Y()), f.startX+uint16(zielP.X()), f.startY+uint16(zielP.Y()))
		// zeichne die Stoßrichtung und -stärke bezogen auf Kugelradien
		stärke := f.billard.GibVStoß().Betrag()
		var farbe Farbe
		if stärke < 7 {
			farbe = F(47, 159, 52)
		} else if stärke > 9 {
			farbe = F(249, 73, 68)
		} else {
			farbe = F(250, 175, 50)
		}
		gfxBreiteLinie(f.startX, f.startY,
			pK, pK.Plus(f.billard.GibVStoß().Mal(ra)),
			4, farbe)
		// Schreibe den Wert der Stärke daneben
		gfx.Stiftfarbe(100, 100, 100)
		schriftgröße := int(ra*0.67 + 0.5)
		gfx.SetzeFont(fp, schriftgröße)
		pStärke := pK.Plus(f.billard.GibVStoß().Mal(ra * 3 / 4))
		gfx.SchreibeFont(f.startX+uint16(pStärke.X()), f.startY+uint16(pStärke.Y()-2*ra), fmt.Sprintf("Stärke: %d", uint16(stärke+0.5)))
	}
	// debugging
	if f.billard.IstZeitlupe() {
		// zeichne Geschwindigkeiten
		for _, k := range f.billard.GibAktiveKugeln() {
			if !k.GibV().IstNull() {
				gfxBreiteLinie(f.startX, f.startY,
					k.GibPos(), k.GibPos().Plus(k.GibV().Mal(k.GibRadius())),
					2, F(250, 175, 50))
			}
		}
		gfx.Stiftfarbe(100, 100, 100)
		gfx.SetzeFont(fp, int(f.billard.GibStoßkugel().GibRadius()+0.5))
		gfx.SchreibeFont(4*breite/5, f.startY+5, "Zeitlupe")
	}
}
