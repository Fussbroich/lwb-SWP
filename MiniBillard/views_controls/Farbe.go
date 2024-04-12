package views_controls

type Farbe interface {
	RGB() (r, g, b uint8)
}
