package hilf

import "math"

type Vec2 interface {
	X() float64
	Y() float64
	IstNull() bool
	Betrag() float64
	Normiert() Vec2
	Punkt(Vec2) float64
	Plus(Vec2) Vec2
	Minus(Vec2) Vec2
	Mal(float64) Vec2
	ProjiziertAuf(Vec2) Vec2
}

type v2 struct {
	x, y float64
}

func V2(x, y float64) *v2 {
	return &v2{x: x, y: y}
}

func V2null() *v2 {
	return &v2{}
}

func (v *v2) X() float64 {
	return v.x
}

func (v *v2) Y() float64 {
	return v.y
}

func (v *v2) IstNull() bool {
	return v.x == 0 && v.y == 0
}

func (v *v2) Betrag() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v *v2) Normiert() Vec2 {
	len := v.Betrag()
	if len == 0 {
		return &v2{x: v.x, y: v.y}
	}
	return &v2{x: v.x / len, y: v.y / len}
}

func (a *v2) Punkt(b Vec2) float64 {
	return a.x*b.X() + a.y*b.Y()
}

func (a *v2) Plus(b Vec2) Vec2 {
	return &v2{x: a.x + b.X(), y: a.y + b.Y()}
}

func (a *v2) Minus(b Vec2) Vec2 {
	return &v2{x: a.x - b.X(), y: a.y - b.Y()}
}

func (v *v2) Mal(f float64) Vec2 {
	return &v2{x: v.x * f, y: v.y * f}
}

func (v *v2) ProjiziertAuf(u Vec2) Vec2 {
	uNorm := u.Normiert()
	return uNorm.Mal(uNorm.Punkt(v))
}
