package views

import (
	"fmt"
	"gfx"

	"../welt"
)

type miniBSpielfeld struct {
	billard welt.MiniBillardSpiel
	fenster
}

func NewMBSpielfeldFenster(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBSpielfeld {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, trans: tr}
	return &miniBSpielfeld{billard: billard, fenster: fenster}
}

func (f *miniBSpielfeld) Zeichne() {
	// zeichne das Tuch
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	// zeichne die Taschen
	ts := f.billard.GibTaschen()
	gfxVollKreissektor(f.startX, f.startY, ts[0].GibPos(), ts[0].GibRadius(), 270, 0, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[1].GibPos(), ts[1].GibRadius(), 0, 90, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[2].GibPos(), ts[2].GibRadius(), 0, 180, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[3].GibPos(), ts[3].GibRadius(), 90, 180, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[4].GibPos(), ts[4].GibRadius(), 180, 270, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[5].GibPos(), ts[5].GibRadius(), 180, 360, Schwarz())
	// zeichne die Kugeln
	for _, k := range f.billard.GibAktiveKugeln() {
		zeichneKugel(f.startX, f.startY, k.GibPos(), k)
	}
	if f.billard.IstStillstand() && !f.billard.GibStoßkugel().IstEingelocht() {
		kS := f.billard.GibStoßkugel()
		pK := kS.GibPos()
		// Zeichne Peillinie bis zum Rand
		gfx.Stiftfarbe(100, 100, 100)
		zielP := pK.Plus(f.billard.GibVStoß().Normiert().Mal(float64(f.stopX - f.startX)))
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
			pK, pK.Plus(f.billard.GibVStoß().Mal(kS.GibRadius())),
			5, farbe)
		// Schreibe den Wert der Stärke daneben
		gfx.Stiftfarbe(100, 100, 100)
		schriftgröße := int(f.billard.GibStoßkugel().GibRadius()*0.67 + 0.5)
		gfx.SetzeFont(fp, schriftgröße)
		pStärke := pK.Plus(f.billard.GibVStoß().Mal(f.billard.GibStoßkugel().GibRadius() / 2))
		gfx.SchreibeFont(uint16(pStärke.X()), uint16(pStärke.Y()), fmt.Sprintf("Stärke: %d", uint16(stärke+0.5)))
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
		gfx.SchreibeFont(4*(f.stopX-f.startX)/5, f.startY+5, "Zeitlupe")
	}
}
