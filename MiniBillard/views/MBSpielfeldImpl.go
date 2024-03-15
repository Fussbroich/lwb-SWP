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
	// zeichne die Stoßstärke und -richtung bezogen auf Kugelradien
	if f.spiel.IstStillstand() && !f.spiel.GibStoßkugel().IstEingelocht() {
		kS := f.spiel.GibStoßkugel()
		stärke := f.spiel.GibVStoß().Betrag()
		gfxBreiteLinie(f.startX, f.startY,
			kS.GibPos(), kS.GibPos().Plus(f.spiel.GibVStoß().Mal(kS.GibRadius())),
			5, F(250, 175, 50))

		schriftgröße := int(f.spiel.GibStoßkugel().GibRadius()*0.67 + 0.5)
		gfx.Stiftfarbe(100, 100, 100)
		gfx.SetzeFont(fp, schriftgröße)
		// Schreibe den Wert der Stärke daneben
		pStärke := kS.GibPos().Plus(f.spiel.GibVStoß().Mal(f.spiel.GibStoßkugel().GibRadius() / 2))
		gfx.SchreibeFont(uint16(pStärke.X()), uint16(pStärke.Y()), fmt.Sprintf("Stärke: %d", uint16(stärke+0.5)))
	}
	// debugging: zeichne Geschwindigkeiten
	if f.spiel.IstDebugMode() {
		for _, k := range f.spiel.GibAktiveKugeln() {
			if !k.GibV().IstNull() {
				gfxBreiteLinie(f.startX, f.startY,
					k.GibPos(), k.GibPos().Plus(k.GibV().Mal(k.GibRadius())),
					2, F(250, 175, 50))
			}
		}
	}
}
