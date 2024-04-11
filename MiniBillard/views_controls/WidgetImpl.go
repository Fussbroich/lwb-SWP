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
	f.Stiftfarbe(f.hg)
	f.Transparenz(f.transparenz)
	br, ho := f.GibGroesse()
	if f.eckradius > 0 {
		f.VollRechteckGFX(f.eckradius, 0, br-2*f.eckradius, ho)
		f.VollRechteckGFX(0, f.eckradius, br, ho-2*f.eckradius)
		f.VollKreisGFX(f.eckradius, f.eckradius, f.eckradius)
		f.VollKreisGFX(f.eckradius, ho-f.eckradius, f.eckradius)
		f.VollKreisGFX(br-f.eckradius, ho-f.eckradius, f.eckradius)
		f.VollKreisGFX(br-f.eckradius, f.eckradius, f.eckradius)
	} else {
		f.VollRechteckGFX(0, 0, br, ho)
	}
	f.Stiftfarbe(f.vg)
	f.Transparenz(0)
}

func (f *widget) Zeichne() {
	f.Stiftfarbe(f.hg)
	f.Transparenz(f.transparenz)
	br, ho := f.GibGroesse()
	if f.eckradius > 0 {
		f.VollRechteckGFX(f.eckradius, 0, br-2*f.eckradius, ho)
		f.VollRechteckGFX(0, f.eckradius, br, ho-2*f.eckradius)
		f.VollKreisGFX(f.eckradius, f.eckradius, f.eckradius)
		f.VollKreisGFX(f.eckradius, ho-f.eckradius, f.eckradius)
		f.VollKreisGFX(br-f.eckradius, ho-f.eckradius, f.eckradius)
		f.VollKreisGFX(br-f.eckradius, f.eckradius, f.eckradius)
	} else {
		f.VollRechteckGFX(0, 0, br, ho)
	}
	f.Stiftfarbe(f.vg)
	f.Transparenz(0)
}

func (f *widget) ZeichneRand() {
	f.Stiftfarbe(f.vg)
	f.Transparenz(0)
	br, ho := f.GibGroesse()
	if f.eckradius > 0 {
		f.KreissektorGFX(f.eckradius, f.eckradius, f.eckradius, 90, 180)
		f.KreissektorGFX(f.eckradius, ho-f.eckradius, f.eckradius, 180, 270)
		f.KreissektorGFX(br-f.eckradius, ho-f.eckradius, f.eckradius, 270, 0)
		f.KreissektorGFX(br-f.eckradius, f.eckradius, f.eckradius, 0, 90)
		f.LinieGFX(f.eckradius, 0, br-f.eckradius, 0)
		f.LinieGFX(f.eckradius, 0, br-f.eckradius, 0)
		f.LinieGFX(f.eckradius, ho, br-f.eckradius, ho)
		f.LinieGFX(f.eckradius, ho, br-f.eckradius, ho)
	} else {
		f.LinieGFX(0, 0, br, 0)
		f.LinieGFX(0, ho, br, ho)
		f.LinieGFX(0, 0, 0, ho)
		f.LinieGFX(br, 0, br, ho)
	}
}

// ######## Hilfsfunktionen zum Zeichnen ############################################################

func (f *widget) Stiftfarbe(c Farbe) {
	r, g, b := c.RGB()
	f.StiftfarbeGFX(r, g, b)
}

func (f *widget) StiftfarbeGFX(r, g, b uint8) {
	gfx.Stiftfarbe(r, g, b)
}

func (f *widget) Transparenz(tr uint8) {
	gfx.Transparenz(tr)
}

func (f *widget) LinieGFX(xV, yV, xN, yN uint16) {
	gfx.Linie(f.startX+xV, f.startY+yV, f.startX+xN, f.startY+yN)
}

func (f *widget) VollRechteckGFX(xV, yV, b, h uint16) {
	gfx.Vollrechteck(f.startX+xV, f.startY+yV, b, h)
}

func (f *widget) VollKreis(pos hilf.Vec2, radius float64, c Farbe) {
	f.Stiftfarbe(c)
	f.VollKreisGFX(uint16(0.5+pos.X()), uint16(0.5+pos.Y()), uint16(0.5+radius))
}

func (f *widget) VollKreisGFX(x, y, ra uint16) {
	gfx.Vollkreis(f.startX+x, f.startY+y, ra)
}

func (f *widget) Kreissektor(pos hilf.Vec2, radius float64, wVon, wBis uint16, c Farbe) {
	f.Stiftfarbe(c)
	f.KreissektorGFX(uint16(0.5+pos.X()), uint16(0.5+pos.Y()), uint16(0.5+radius), wVon, wBis)
}

func (f *widget) KreissektorGFX(x, y, ra, wVon, wBis uint16) {
	gfx.Kreissektor(f.startX+x, f.startY+y, ra, wVon, wBis)
}

func (f *widget) VollKreissektor(pos hilf.Vec2, radius float64, wVon, wBis uint16, c Farbe) {
	f.Stiftfarbe(c)
	f.VollKreissektorGFX(uint16(0.5+pos.X()), uint16(0.5+pos.Y()), uint16(0.5+radius), wVon, wBis)
}

func (f *widget) VollKreissektorGFX(x, y, ra, wVon, wBis uint16) {
	gfx.Vollkreissektor(f.startX+x, f.startY+y, ra, wVon, wBis)
}

func (f *widget) VollDreieckGFX(x1, y1, x2, y2, x3, y3 uint16) {
	gfx.Volldreieck(
		f.startX+x1, f.startY+y1,
		f.startX+x2, f.startY+y2,
		f.startX+x3, f.startY+y3)
}

func (f *widget) BreiteLinie(pV, pN hilf.Vec2, breite float64, c Farbe) {
	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())
	f.Stiftfarbe(c)

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))

	gfx.Volldreieck(f.startX+uint16(0.5+pA.X()), f.startY+uint16(0.5+pA.Y()),
		f.startX+uint16(0.5+pB.X()), f.startY+uint16(0.5+pB.Y()),
		f.startX+uint16(0.5+pC.X()), f.startY+uint16(0.5+pC.Y()))
	gfx.Volldreieck(f.startX+uint16(0.5+pA.X()), f.startY+uint16(0.5+pA.Y()),
		f.startX+uint16(0.5+pC.X()), f.startY+uint16(0.5+pC.Y()),
		f.startX+uint16(0.5+pD.X()), f.startY+uint16(0.5+pD.Y()))
}
