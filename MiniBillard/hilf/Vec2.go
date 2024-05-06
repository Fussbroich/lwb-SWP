package hilf

// Ein Vec2 ist ein Wrapper-Objekt für zwei float64-Werte,
// die allgemein als 2D-Vektor (Ortsvektor oder relativ
// je nach Kontext) interpretiert werden.
type Vec2 interface {
	// Getter für die x-Komponente.
	X() float64
	// Getter für die y-Komponente.
	Y() float64
	// Feststellen, ob beide Komponenten 0 sind.
	IstNull() bool
	// Getter für die Länge des Vektors.
	Betrag() float64
	// Getter für den normierten Vektor.
	Normiert() Vec2
	// Skalarprodukt mit einem anderen Vec2.
	Punkt(Vec2) float64
	// Summe mit einem anderen Vec2.
	Plus(Vec2) Vec2
	// Differenz zu einem anderen Vec2.
	Minus(Vec2) Vec2
	// Vektorprodukt mit einem anderen Vec2.
	Mal(float64) Vec2
	// Projektion auf einen anderen Vec2.
	ProjiziertAuf(Vec2) Vec2
}
