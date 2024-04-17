package hilf

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
