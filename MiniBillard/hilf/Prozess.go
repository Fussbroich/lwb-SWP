package hilf

import "time"

type Prozess interface {
	StarteLoop(time.Duration)
	StarteRate(uint64)
	Starte()
	GibRate() uint64
	Stoppe()
	Läuft() bool
}
