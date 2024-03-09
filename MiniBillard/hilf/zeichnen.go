package hilf

import (
	"gfx"
)

func ZeichneVollDreieck(p1, p2, p3 Vec2) {
	gfx.Volldreieck(uint16(p1.X()), uint16(p1.Y()),
		uint16(p2.X()), uint16(p2.Y()),
		uint16(p3.X()), uint16(p3.Y()))
}

func ZeichneBreiteLinie(pV Vec2, pN Vec2, breite float64) {
	v := pN.Minus(pV).Normiert()
	d := V2(v.Y(), -v.X())

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	ZeichneVollDreieck(pA, pB, pC)
	ZeichneVollDreieck(pA, pC, pD)
	gfx.Vollkreis(uint16(pV.X()), uint16(pV.Y()), uint16(breite/2))
	gfx.Vollkreis(uint16(pN.X()), uint16(pN.Y()), uint16(breite/2))
}
