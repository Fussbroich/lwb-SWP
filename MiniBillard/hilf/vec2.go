package hilf

import "math"

type Vec2 struct {
	x, y float64
}

func (v Vec2) X() float64 {
	return v.x
}

func (v Vec2) Y() float64 {
	return v.y
}

type Gerade2 struct {
	x Vec2
	r Vec2
}

func G2(x, r Vec2) Gerade2 {
	return Gerade2{x: x, r: r}
}

func (g Gerade2) GibPosFür(t float64) Vec2 {
	return g.x.Plus(g.r.Mal(t))
}

func V2(x, y float64) Vec2 {
	return Vec2{x: x, y: y}
}

func (v Vec2) IstNull() bool {
	return v.x == 0 && v.y == 0
}

func (v Vec2) Betrag() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v Vec2) Normiert() Vec2 {
	len := v.Betrag()
	if len == 0 {
		return Vec2{x: v.x, y: v.y}
	}
	return Vec2{x: v.x / len, y: v.y / len}
}

func (a Vec2) Punkt(b Vec2) float64 {
	return a.x*b.x + a.y*b.y
}

func (a Vec2) Plus(b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func (a Vec2) Minus(b Vec2) Vec2 {
	return Vec2{x: a.x - b.x, y: a.y - b.y}
}

func (v Vec2) Mal(f float64) Vec2 {
	return Vec2{x: v.x * f, y: v.y * f}
}

func (v Vec2) ProjiziertAuf(u Vec2) Vec2 {
	uNorm := u.Normiert()
	return uNorm.Mal(uNorm.Punkt(v))
}
