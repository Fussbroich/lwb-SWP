package hilf

import "time"

// Ein Wrapper-Objekt für eine anonyme Funktion
// und einigen Methoden zur Steuerung. Jede Routine wrappt
// eine Funktion, die als go-Routine in einer Endlosschleife
// ausgeführt wird. Die Funktion wird zusätzlich zu einem Namen
// für die Routine dem Konstruktor übergeben:
// NewRoutine(name string, f_run func())
type Routine interface {
	// Optionale Funktion, die mit defer vor starten einer
	// Endlosschleife gerufen wird. Die Funktion kann ein recover
	// enthalten zum Behandeln von panics.
	SetzeAusnahmeHandler(func())
	// Starte die Funktion als go-Routine in einem Loop mit einem vorgegebenen Zeittakt.
	StarteLoop(time.Duration)
	// Starte die Funktion als go-Routine in einem Loop mit einer vorgegebenen Soll-Frequenz.
	// Die Sollfrequenz wird nur durch eine geregelte Verzögerung angepasst, so
	// dass sie gar nicht erreicht werden kann, falls die Funktion zu lange braucht.
	StarteRate(uint64)
	// Starte die Funktion als go-Routine in einem Loop. Die Funktion wird so schnell
	// wie möglich wiederholt aufgerufen.
	Starte()
	// Starte die Funktion in einem lokalen Loop. Die Funktion wird so schnell
	// wie möglich wiederholt aufgerufen. Die Methode blockiert, so dass sie als
	// letzter Aufruf in einem Prozess stehen sollte.
	StarteHier() // Läuft ohne go-Routine
	// Getter für die aktuelle Ausführrate (Frequenz) der Funktion.
	GibRate() uint64
	// Getter für den Namen der Routine.
	GibName() string
	// Stoppe die Ausführung der Funktion.
	Stoppe()
	// Getter für den Zustand der Routine.
	Laeuft() bool
}
