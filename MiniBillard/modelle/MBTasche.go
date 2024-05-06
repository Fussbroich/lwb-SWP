package modelle

import "../hilf"

// Modelliert eine kreisrunde Tasche auf dem Billardtisch. Eine MBKugel
// gilt als eingelocht, wenn ihr Mittelpunkt den Rand des Kreises
// überschreitet.
// Konstruktor: NewTasche(pos hilf.Vec2, r float64)
type MBTasche interface {
	// Getter für die Position der Tasche (zur Darstellung)
	GibPos() hilf.Vec2
	// Getter für den Radius der Tasche (zur Darstellung)
	GibRadius() float64
}
