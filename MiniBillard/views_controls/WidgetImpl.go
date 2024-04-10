package views_controls

import (
	"fmt"
	"gfx"

	"../fonts"
	"../hilf"
	"../modelle"
)

type widget struct {
	startX, startY uint16
	stopX, stopY   uint16
	hg, vg         Farbe
	transparenz    uint8
	eckradius      uint16
}

func NewFenster(startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8, ra uint16) *widget {
	return &widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr, eckradius: ra}
}

func (f *widget) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *widget) GibGröße() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

func (f *widget) ImFenster(x, y uint16) bool {
	xs, ys := f.startX+f.eckradius*3/10, f.startY+f.eckradius*3/10
	b, h := f.stopX-f.startX-f.eckradius*6/10, f.stopY-f.startY-f.eckradius*6/10
	return x > xs && x < xs+b && y > ys && y < ys+h
}

func (f *widget) MausklickBei(x, y uint16) {}

func (f *widget) ZeichneLayout() {
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(f.transparenz)
	if f.eckradius > 0 {
		gfx.Vollrechteck(f.startX+f.eckradius, f.startY, f.stopX-f.startX-2*f.eckradius, f.stopY-f.startY)
		gfx.Vollrechteck(f.startX, f.startY+f.eckradius, f.stopX-f.startX, f.stopY-f.startY-2*f.eckradius)
		gfx.Vollkreis(f.startX+f.eckradius, f.startY+f.eckradius, f.eckradius)
		gfx.Vollkreis(f.startX+f.eckradius, f.stopY-f.eckradius, f.eckradius)
		gfx.Vollkreis(f.stopX-f.eckradius, f.stopY-f.eckradius, f.eckradius)
		gfx.Vollkreis(f.stopX-f.eckradius, f.startY+f.eckradius, f.eckradius)
	} else {
		gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
	}
	r, g, b = f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(0)
}

func (f *widget) Zeichne() {
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(f.transparenz)
	if f.eckradius > 0 {
		gfx.Vollrechteck(f.startX+f.eckradius, f.startY, f.stopX-f.startX-2*f.eckradius, f.stopY-f.startY)
		gfx.Vollrechteck(f.startX, f.startY+f.eckradius, f.stopX-f.startX, f.stopY-f.startY-2*f.eckradius)
		gfx.Vollkreis(f.startX+f.eckradius, f.startY+f.eckradius, f.eckradius)
		gfx.Vollkreis(f.startX+f.eckradius, f.stopY-f.eckradius, f.eckradius)
		gfx.Vollkreis(f.stopX-f.eckradius, f.stopY-f.eckradius, f.eckradius)
		gfx.Vollkreis(f.stopX-f.eckradius, f.startY+f.eckradius, f.eckradius)
	} else {
		gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
	}
	r, g, b = f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(0)
}

func (f *widget) ZeichneRand() {
	r, g, b := f.vg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Transparenz(0)
	if f.eckradius > 0 {
		gfx.Kreissektor(f.startX+f.eckradius, f.startY+f.eckradius, f.eckradius+1, 90, 180)
		gfx.Kreissektor(f.startX+f.eckradius, f.startY+f.eckradius, f.eckradius+2, 90, 180)
		gfx.Kreissektor(f.startX+f.eckradius, f.stopY-f.eckradius, f.eckradius+1, 180, 270)
		gfx.Kreissektor(f.startX+f.eckradius, f.stopY-f.eckradius, f.eckradius+2, 180, 270)
		gfx.Kreissektor(f.stopX-f.eckradius, f.stopY-f.eckradius, f.eckradius+1, 270, 0)
		gfx.Kreissektor(f.stopX-f.eckradius, f.stopY-f.eckradius, f.eckradius+2, 270, 0)
		gfx.Kreissektor(f.stopX-f.eckradius, f.startY+f.eckradius, f.eckradius+2, 0, 90)
		gfx.Linie(f.startX+f.eckradius, f.startY-1, f.stopX-f.eckradius, f.startY-1)
		gfx.Linie(f.startX+f.eckradius, f.startY-2, f.stopX-f.eckradius, f.startY-2)
		gfx.Linie(f.startX+f.eckradius, f.stopY+1, f.stopX-f.eckradius, f.stopY+1)
		gfx.Linie(f.startX+f.eckradius, f.stopY+2, f.stopX-f.eckradius, f.stopY+2)
		gfx.Linie(f.startX-1, f.startY+f.eckradius, f.startX-1, f.stopY-f.eckradius)
		gfx.Linie(f.startX-2, f.startY+f.eckradius, f.startX-2, f.stopY-f.eckradius)
		gfx.Linie(f.stopX+1, f.startY+f.eckradius, f.stopX+1, f.stopY-f.eckradius)
		gfx.Linie(f.stopX+2, f.startY+f.eckradius, f.stopX+2, f.stopY-f.eckradius)
	} else {
		gfx.Linie(f.startX-1, f.startY-1, f.stopX+1, f.startY-1)
		gfx.Linie(f.startX-1, f.stopY+1, f.stopX+1, f.stopY+1)
		gfx.Linie(f.startX-1, f.startY-1, f.startX-1, f.stopY+1)
		gfx.Linie(f.stopX+1, f.startY-1, f.stopX+1, f.stopY+1)
	}
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

func zeichneKugel(startX, startY uint16, p hilf.Vec2, k modelle.MBKugel) {
	font := fonts.LiberationMonoBold(int(k.GibRadius()) - 3)
	gfxVollKreis(startX, startY, p, k.GibRadius(), F(48, 49, 54))
	gfxVollKreis(startX, startY, p, k.GibRadius()-1, F(252, 253, 242))
	c := mBKugelPalette()[k.GibWert()]
	if k.GibWert() <= 8 {
		gfxVollKreis(startX, startY, p, k.GibRadius()-1, c)
	} else {
		r, g, b := c.RGB()
		gfx.Stiftfarbe(r, g, b)
		gfx.Vollrechteck(startX+uint16(p.X()-k.GibRadius()*0.75+0.5), startY+uint16(p.Y()-k.GibRadius()*0.6+0.5),
			uint16(2*0.75*k.GibRadius()+0.5), uint16(2*0.6*k.GibRadius()+0.5))
		gfxVollKreissektor(startX, startY, p, k.GibRadius()-1, 325, 35, c)
		gfxVollKreissektor(startX, startY, p, k.GibRadius()-1, 145, 215, c)
	}
	if k.GibWert() != 0 {
		gfxVollKreis(startX, startY, p, (k.GibRadius()-1)/2, F(252, 253, 242))
		gfx.Stiftfarbe(0, 0, 0)
		gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgröße())
		if k.GibWert() < 10 {
			gfx.SchreibeFont(
				startX-uint16(font.GibSchriftgröße())/4+uint16(p.X()+0.5),
				startY-uint16(font.GibSchriftgröße())/2+uint16(p.Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		} else {
			gfx.SchreibeFont(
				startX-uint16(font.GibSchriftgröße())/2+uint16(p.X()+0.5),
				startY-uint16(font.GibSchriftgröße())/2+uint16(p.Y()+0.5),
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
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	gfx.Volldreieck(startX+uint16(0.5+pA.X()), startY+uint16(0.5+pA.Y()), startX+uint16(0.5+pB.X()), startY+uint16(0.5+pB.Y()), startX+uint16(0.5+pC.X()), startY+uint16(0.5+pC.Y()))
	gfx.Volldreieck(startX+uint16(0.5+pA.X()), startY+uint16(0.5+pA.Y()), startX+uint16(0.5+pC.X()), startY+uint16(0.5+pC.Y()), startX+uint16(0.5+pD.X()), startY+uint16(0.5+pD.Y()))
}
