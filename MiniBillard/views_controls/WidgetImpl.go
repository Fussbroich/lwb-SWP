package views_controls

import (
	"gfx"

	"../hilf"
	//	"../modelle"
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

func (f *widget) SetzeKoordinaten(startx, starty, stopx, stopy uint16) {
	f.startX, f.startY, f.stopY, f.stopY = startx, starty, stopx, stopy
}

func (f *widget) SetzeFarben(hg, vg Farbe) {
	f.hg, f.vg = hg, vg
}

func (f *widget) SetzeTransparenz(tr uint8) {
	f.transparenz = tr
}

func (f *widget) SetzeEckradius(ra uint16) {
	f.eckradius = ra
}

func (f *widget) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *widget) GibGroesse() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

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
