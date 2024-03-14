package views

type Farbe interface {
	RGB() (r, g, b uint8)
}
