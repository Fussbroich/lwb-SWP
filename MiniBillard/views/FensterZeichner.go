package views

import (
	"gfx"

	"../hilf"
	"../welt"
)

type FensterZeichner struct {
	startX, startY uint16
	stopX, stopY   uint16
	maßstab        float64
}

func NewFensterZeichner(startx, starty, stopx, stopy uint16, maßstab float64) *FensterZeichner {
	return &FensterZeichner{
		startX: startx, startY: starty,
		stopX: stopx, stopY: stopy,
		maßstab: maßstab}
}

func (f *FensterZeichner) ZeichneMiniBillardSpiel(spiel welt.MiniBillardSpiel) {
	gfx.Cls()
	l, b := spiel.GibGröße()
	// zeichne den Belag
	f.ZeichneVollRechteck(hilf.V2(0, 0), l, b, 60, 179, 113)
	// zeichne die Taschen
	for _, t := range spiel.GibTaschen() {
		f.ZeichneVollKreis(t.GibPos(), t.GibRadius(), 0, 0, 0)
	}
	// zeichne die Kugeln
	for _, k := range spiel.GibAktiveKugeln() {
		r, g, b := k.GibFarbe()
		f.ZeichneVollKreis(k.GibPos(), k.GibRadius(), 0, 0, 0)
		f.ZeichneVollKreis(k.GibPos(), k.GibRadius()-1, r, g, b)
	}
}

func (f *FensterZeichner) FülleFläche(r, g, b uint8) {
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
}

func (f *FensterZeichner) ZeichneVollKreis(pos hilf.Vec2, radius float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		f.startX+uint16(pos.X()*f.maßstab), f.startY+uint16(pos.Y()*f.maßstab),
		uint16(radius*f.maßstab))
}

func (f *FensterZeichner) ZeichneVollRechteck(pos hilf.Vec2, breite, höhe float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(
		f.startX+uint16(pos.X()*f.maßstab), f.startY+uint16(pos.Y()*f.maßstab),
		uint16(breite*f.maßstab), uint16(höhe*f.maßstab))
}

func (f *FensterZeichner) ZeichneVollDreieck(p1, p2, p3 hilf.Vec2) {
	gfx.Volldreieck(
		f.startX+uint16(p1.X()*f.maßstab), f.startY+uint16(p1.Y()*f.maßstab),
		f.startX+uint16(p2.X()*f.maßstab), f.startY+uint16(p2.Y()*f.maßstab),
		f.startX+uint16(p3.X()*f.maßstab), f.startY+uint16(p3.Y()*f.maßstab))
}

func (f *FensterZeichner) ZeichneBreiteLinie(pV, pN hilf.Vec2, breite float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)

	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	f.ZeichneVollDreieck(pA, pB, pC)
	f.ZeichneVollDreieck(pA, pC, pD)
	gfx.Vollkreis(
		f.startX+uint16(pV.X()*f.maßstab), f.startY+uint16(pV.Y()*f.maßstab),
		uint16(f.maßstab*breite/2))
	gfx.Vollkreis(
		f.startX+uint16(pN.X()*f.maßstab), f.startX+uint16(pN.Y()*f.maßstab),
		uint16(f.maßstab*breite/2))
}
