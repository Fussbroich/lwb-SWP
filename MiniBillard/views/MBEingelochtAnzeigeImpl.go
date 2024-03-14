package views

import (
	"fmt"
	"gfx"

	"../welt"
)

type miniBEingelochte struct {
	spiel          welt.MiniBillardSpiel
	startX, startY uint16
	stopX, stopY   uint16
	hg, vg         Farbe
}

func NewMBEingelochteFenster(spiel welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg Farbe) *miniBEingelochte {
	return &miniBEingelochte{
		spiel:  spiel,
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		hg: hg, vg: Schwarz()}
}

func (f *miniBEingelochte) Zeichne() {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	// zeichne den Hintergrund
	gfx.Stiftfarbe(80, 80, 80)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	gfx.Stiftfarbe(180, 180, 180)
	schriftgröße := höhe / 3
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Eingelocht", len(f.spiel.GibEingelochteKugeln())))
}
