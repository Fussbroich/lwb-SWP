package modelle

import "time"

// Wrapper-Objekt für einen bestimmten Zeitraum,
// der durch Methoden manipulierbar ist.
type Countdown interface {
	// Getter für den Stand des Countdowns.
	GibRestzeit() time.Duration
	// Setter für den Stand.
	Setze(d time.Duration)
	// Zieht einen bestimmten Zeitraum vom aktuellen Stand ab.
	ZieheAb(d time.Duration)
	// Gibt an, ob der Stand bei 0 liegt.
	IstAbgelaufen() bool
	// Sperrt das Abziehen (Herunterzählen).
	Halt()
	// Entsperrt das Abziehen (Herunterzählen) wieder.
	Weiter()
}
