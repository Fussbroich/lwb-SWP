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

func NewMBEingelochteFenster(spiel welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe) *miniBEingelochte {
	return &miniBEingelochte{
		spiel:  spiel,
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		hg: hg, vg: vg}
}

func (f *miniBEingelochte) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *miniBEingelochte) GibGröße() (uint16, uint16) {
	return f.stopX - f.startX, f.stopY - f.startY
}

func (f *miniBEingelochte) Zeichne() {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	// zeichne den Hintergrund
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	r, g, b = f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	schriftgröße := höhe / 3
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Eingelocht", len(f.spiel.GibEingelochteKugeln())))
}
