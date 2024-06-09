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

// praktische Farben
func Weiss() Farbe {
	return &rgb{r: 255, g: 255, b: 255}
}

func Schwarz() Farbe {
	return &rgb{}
}

func Rot() Farbe {
	return &rgb{r: 255}
}

func Gruen() Farbe {
	return &rgb{g: 255}
}

func Blau() Farbe {
	return &rgb{b: 255}
}
