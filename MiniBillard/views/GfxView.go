package views

import (
	"gfx"

	"../hilf"
)

type GfxView struct {
	startX, startY, stopX uint16
	maßstab               float64
}

func NewGfxView(startx, stopx, starty, gesamt uint16) *GfxView {
	return &GfxView{
		startX: startx,
		stopX:  stopx, startY: starty,
		maßstab: float64(stopx-startx) / float64(gesamt)}
}

func (v *GfxView) ZeichneVollKreis(pos hilf.Vec2, radius float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		v.startX+uint16(pos.X()*v.maßstab), v.startY+uint16(pos.Y()*v.maßstab),
		uint16(radius*v.maßstab))
}

func (v *GfxView) ZeichneVollRechteck(pos hilf.Vec2, breite, höhe float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(
		v.startX+uint16(pos.X()*v.maßstab), v.startY+uint16(pos.Y()*v.maßstab),
		uint16(breite*v.maßstab), uint16(höhe*v.maßstab))
}

func (v *GfxView) ZeichneVollDreieck(p1, p2, p3 hilf.Vec2) {
	gfx.Volldreieck(
		v.startX+uint16(p1.X()*v.maßstab), v.startY+uint16(p1.Y()*v.maßstab),
		v.startX+uint16(p2.X()*v.maßstab), v.startY+uint16(p2.Y()*v.maßstab),
		v.startX+uint16(p3.X()*v.maßstab), v.startY+uint16(p3.Y()*v.maßstab))
}

func (v *GfxView) ZeichneBreiteLinie(pV, pN hilf.Vec2, breite float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)

	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	v.ZeichneVollDreieck(pA, pB, pC)
	v.ZeichneVollDreieck(pA, pC, pD)
	gfx.Vollkreis(
		v.startX+uint16(pV.X()*v.maßstab), v.startY+uint16(pV.Y()*v.maßstab),
		uint16(v.maßstab*breite/2))
	gfx.Vollkreis(
		v.startX+uint16(pN.X()*v.maßstab), v.startX+uint16(pN.Y()*v.maßstab),
		uint16(v.maßstab*breite/2))
}
