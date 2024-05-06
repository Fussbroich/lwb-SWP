package views_controls

// Wrapper-Objekt für die drei Werte r,g,b im RGB-Farbmodell.
type Farbe interface {
	// Getter für die drei Werte r, g, b der Farbe.
	RGB() (r, g, b uint8)
}
