package views

import (
	"fmt"
	"gfx"

	"../welt"
)

type miniBSpielinfo struct {
	spiel          welt.MiniBillardSpiel
	startX, startY uint16
	stopX, stopY   uint16
	hg, vg         Farbe
}

func NewMBSpielinfoFenster(spiel welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe) *miniBSpielinfo {
	return &miniBSpielinfo{
		spiel:  spiel,
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		hg: hg, vg: vg}
}

func (f *miniBSpielinfo) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *miniBSpielinfo) GibGröße() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

func (f *miniBSpielinfo) Zeichne() {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	// zeichne den Hintergrund
	cr, cg, cb := f.hg.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	cr, cg, cb = f.vg.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	schriftgröße := höhe / 4
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Stöße/%d Strafen", f.spiel.GibStößeBisher(), f.spiel.GibStrafpunkte()))
}
