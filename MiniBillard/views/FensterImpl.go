package views

import (
	"errors"
	"fmt"
	"gfx"
	"os"
	"path/filepath"

	"../hilf"
	"../welt"
)

type fenster struct {
	startX, startY uint16
	stopX, stopY   uint16
	hg, vg         Farbe
	trans          uint8
}

func NewFenster(startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *fenster {
	return &fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, trans: tr}
}

func (f *fenster) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *fenster) GibGröße() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

func (f *fenster) ZeichneLayout() {
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(f.trans)
	gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
	gfx.Transparenz(0)
}

func (f *fenster) Zeichne() {
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(f.trans)
	gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
	gfx.Transparenz(0)
}

// ######## Hilfsfunktionen #######################################################################

var (
	kugelPalette *[16]Farbe
)

func mBKugelPalette() *[16]Farbe {
	if kugelPalette == nil {
		kugelPalette = &[16]Farbe{
			F(252, 253, 242), // weiß
			F(255, 201, 78),  // gelb
			F(34, 88, 175),   // blau
			F(249, 73, 68),   // hellrot
			F(84, 73, 149),   // violett
			F(255, 139, 33),  // orange
			F(47, 159, 52),   // grün
			F(155, 53, 30),   // dunkelrot
			F(48, 49, 54),    // schwarz
			F(255, 201, 78),  // gelb
			F(34, 88, 175),   // blau
			F(249, 73, 68),   // hellrot
			F(84, 73, 149),   // violett
			F(255, 139, 33),  // orange
			F(47, 159, 52),   // grün
			F(155, 53, 30)}   // dunkelrot
	}
	return kugelPalette
}

func fontDateipfad(filename string) string {
	fontsDir := "MiniBillard/fonts"
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	workDir := filepath.Dir(wd)
	fp := filepath.Join(workDir, fontsDir, filename)
	//	println("wdir:", wdir)
	//	println("fontsDir:", fontsDir)
	//	println("filename", filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return fp
}

func zeichneKugel(startX, startY uint16, p hilf.Vec2, k welt.Kugel) {
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	schriftgröße := int(k.GibRadius()) - 3
	gfxVollKreis(startX, startY, p, k.GibRadius(), F(48, 49, 54))
	gfxVollKreis(startX, startY, p, k.GibRadius()-1, F(252, 253, 242))
	c := mBKugelPalette()[k.GibWert()]
	if k.GibWert() <= 8 {
		gfxVollKreis(startX, startY, p, k.GibRadius()-1, c)
	} else {
		r, g, b := c.RGB()
		gfx.Stiftfarbe(r, g, b)
		gfx.Vollrechteck(startX+uint16(p.X()-k.GibRadius()*0.818+0.5),
			startY+uint16(p.Y()-k.GibRadius()*0.61+0.5),
			uint16(2*0.818*k.GibRadius()+0.5), uint16(2*0.61*k.GibRadius()+0.5))
		gfxVollKreissektor(startX, startY, p, k.GibRadius()-1, 325, 35, c)
		gfxVollKreissektor(startX, startY, p, k.GibRadius()-1, 145, 215, c)
	}
	if k.GibWert() != 0 {
		gfxVollKreis(startX, startY, p, (k.GibRadius()-1)/2, F(252, 253, 242))
		gfx.Stiftfarbe(0, 0, 0)
		gfx.SetzeFont(fp, schriftgröße)
		if k.GibWert() < 10 {
			gfx.SchreibeFont(
				startX-uint16(schriftgröße)/4+uint16(p.X()+0.5),
				startY-uint16(schriftgröße)/2+uint16(p.Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		} else {
			gfx.SchreibeFont(
				startX-uint16(schriftgröße)/2+uint16(p.X()+0.5),
				startY-uint16(schriftgröße)/2+uint16(p.Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		}
	}
}

func gfxVollKreis(startX, startY uint16, pos hilf.Vec2, radius float64, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		startX+uint16(0.5+pos.X()), startY+uint16(0.5+pos.Y()),
		uint16(0.5+radius))
}

func gfxVollKreissektor(startX, startY uint16, pos hilf.Vec2, radius float64, wVon, wBis uint16, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreissektor(
		startX+uint16(0.5+pos.X()), startY+uint16(0.5+pos.Y()),
		uint16(0.5+radius), wVon, wBis)
}

func gfxVollRechteck(startX, startY uint16, pos hilf.Vec2, breite, höhe float64, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(
		startX+uint16(0.5+pos.X()), startY+uint16(0.5+pos.Y()),
		uint16(0.5+breite), uint16(0.5+höhe))
}

func gfxVollDreieck(startX, startY uint16, p1, p2, p3 hilf.Vec2, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Volldreieck(
		startX+uint16(0.5+p1.X()), startY+uint16(0.5+p1.Y()),
		startX+uint16(0.5+p2.X()), startY+uint16(0.5+p2.Y()),
		startX+uint16(0.5+p3.X()), startY+uint16(0.5+p3.Y()))
}

func gfxBreiteLinie(startX, startY uint16, pV, pN hilf.Vec2, breite float64, c Farbe) {
	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	gfxVollDreieck(startX, startY, pA, pB, pC, c)
	gfxVollDreieck(startX, startY, pA, pC, pD, c)
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		startX+uint16(0.5+pV.X()), startY+uint16(0.5+pV.Y()), uint16(0.5+breite/2))
	gfx.Vollkreis(
		startX+uint16(0.5+pN.X()), startX+uint16(0.5+pN.Y()), uint16(0.5+breite/2))
}
