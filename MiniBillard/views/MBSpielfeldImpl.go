package views

import (
	"fmt"
	"gfx"

	"../hilf"
	"../welt"
)

type miniBSpielfeld struct {
	spiel          welt.MiniBillardSpiel
	startX, startY uint16
	stopX, stopY   uint16
	hg, vg         Farbe
}

func NewMBSpielfeldFenster(spiel welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16) *miniBSpielfeld {
	return &miniBSpielfeld{
		spiel:  spiel,
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		hg: Weiß(), vg: Schwarz()}
}

func (f *miniBSpielfeld) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *miniBSpielfeld) GibGröße() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

func (f *miniBSpielfeld) Zeichne() {
	fp := fontDateipfad("LiberationMono-Regular.ttf")
	l, b := f.spiel.GibGröße()
	// zeichne das Tuch
	gfxVollRechteck(f.startX, f.startY, hilf.V2(0, 0), l, b, F(60, 179, 113))
	// zeichne die Taschen
	ts := f.spiel.GibTaschen()
	gfxVollKreissektor(f.startX, f.startY, ts[0].GibPos(), ts[0].GibRadius(), 270, 0, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[1].GibPos(), ts[1].GibRadius(), 0, 90, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[2].GibPos(), ts[2].GibRadius(), 0, 180, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[3].GibPos(), ts[3].GibRadius(), 90, 180, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[4].GibPos(), ts[4].GibRadius(), 180, 270, Schwarz())
	gfxVollKreissektor(f.startX, f.startY, ts[5].GibPos(), ts[5].GibRadius(), 180, 360, Schwarz())
	// zeichne die Kugeln
	for _, k := range f.spiel.GibAktiveKugeln() {
		zeichneKugel(f.startX, f.startY, k.GibPos(), k)
	}
	if f.spiel.IstStillstand() && !f.spiel.GibStoßkugel().IstEingelocht() {
		kS := f.spiel.GibStoßkugel()
		pK := kS.GibPos()
		// Zeichne Peillinie bis zum Rand
		gfx.Stiftfarbe(100, 100, 100)
		zielP := pK.Plus(f.spiel.GibVStoß().Normiert().Mal(float64(f.stopX - f.startX)))
		gfx.Linie(f.startX+uint16(pK.X()), f.startY+uint16(pK.Y()), f.startX+uint16(zielP.X()), f.startY+uint16(zielP.Y()))
		// zeichne die Stoßrichtung und -stärke bezogen auf Kugelradien
		stärke := f.spiel.GibVStoß().Betrag()
		var farbe Farbe
		if stärke > 10 {
			farbe = F(249, 73, 68)
		} else {
			farbe = F(250, 175, 50)
		}
		gfxBreiteLinie(f.startX, f.startY,
			pK, pK.Plus(f.spiel.GibVStoß().Mal(kS.GibRadius())),
			5, farbe)
		// Schreibe den Wert der Stärke daneben
		gfx.Stiftfarbe(100, 100, 100)
		schriftgröße := int(f.spiel.GibStoßkugel().GibRadius()*0.67 + 0.5)
		gfx.SetzeFont(fp, schriftgröße)
		pStärke := pK.Plus(f.spiel.GibVStoß().Mal(f.spiel.GibStoßkugel().GibRadius() / 2))
		gfx.SchreibeFont(uint16(pStärke.X()), uint16(pStärke.Y()), fmt.Sprintf("Stärke: %d", uint16(stärke+0.5)))
	}
	// debugging
	if f.spiel.IstZeitlupe() {
		// zeichne Geschwindigkeiten
		for _, k := range f.spiel.GibAktiveKugeln() {
			if !k.GibV().IstNull() {
				gfxBreiteLinie(f.startX, f.startY,
					k.GibPos(), k.GibPos().Plus(k.GibV().Mal(k.GibRadius())),
					2, F(250, 175, 50))
			}
		}
		gfx.Stiftfarbe(100, 100, 100)
		gfx.SetzeFont(fp, int(f.spiel.GibStoßkugel().GibRadius()+0.5))
		gfx.SchreibeFont(4*(f.stopX-f.startX)/5, f.startY+5, "Zeitlupe")
	}
}
