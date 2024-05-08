package views_controls

// Wrapper-Objekt für die drei Werte r,g,b im RGB-Farbmodell.
//
// Konstruktor: F(r, g, b uint8)
type Farbe interface {
	// Getter für die drei Werte r, g, b der Farbe.
	RGB() (r, g, b uint8)
}
