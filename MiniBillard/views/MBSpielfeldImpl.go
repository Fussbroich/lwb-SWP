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

func NewMBSpieltischFenster(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8, ra uint16) *miniBSpielfeld {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
	return &miniBSpielfeld{billard: billard, fenster: fenster}
}

/*
	func (f *fenster) findeSchnittPunktMitRand(u, r hilf.Vec2) hilf.Vec2 {
		mx := func(x float64) float64 { return (x - u.X()) / r.X() }
		my := func(y float64) float64 { return (y - u.Y()) / r.Y() }
		xm := func(m float64) float64 { return u.X() + m*r.X() }
		ym := func(m float64) float64 { return u.Y() + m*r.Y() }

		var m, x, y float64
		if m = my(float64(f.startY)); m >= 0 {
			x = xm(m)
			if f.startX <= uint16(x+0.5) && f.stopX >= uint16(x+0.5) {
				return hilf.V2(x, float64(f.startY))
			}
		}
		if m = my(float64(f.stopY)); m >= 0 {
			x = xm(m)
			if f.startX <= uint16(x+0.5) && f.stopX >= uint16(x+0.5) {
				return hilf.V2(x, float64(f.stopY))
			}
		}
		if m = mx(float64(f.startX)); m >= 0 {
			y = ym(m)
			if f.startY <= uint16(y+0.5) && f.stopY >= uint16(y+0.5) {
				return hilf.V2(float64(f.startX), y)
			}
		}
		if m = mx(float64(f.stopX)); m >= 0 {
			y = ym(m)
			if f.startY <= uint16(y+0.5) && f.stopY >= uint16(y+0.5) {
				return hilf.V2(float64(f.stopX), y)
			}
		}
		return u.Plus(r.Mal(1000))
	}
*/

func (f *miniBSpielfeld) Zeichne() {
	// zeichne das Tuch
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	// zeichne die Taschen

	breite, _ := f.GibGröße()
	kS := f.billard.GibStoßkugel()
	ra := kS.GibRadius()

	for _, t := range f.billard.GibTaschen() {
		gfxVollKreis(f.startX, f.startY, t.GibPos(), ra, Schwarz())
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
