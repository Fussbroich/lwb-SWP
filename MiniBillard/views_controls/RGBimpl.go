package views_controls

type rgb struct {
	r, g, b uint8
}

func F(r, g, b uint8) *rgb {
	return &rgb{r: r, g: g, b: b}
}

func (r *rgb) RGB() (uint8, uint8, uint8) {
	return r.r, r.g, r.b
}

func Weiß() *rgb {
	return &rgb{r: 255, g: 255, b: 255}
}

func Schwarz() *rgb {
	return &rgb{}
}
