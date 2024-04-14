package hilf

import "time"

type Routine interface {
	StarteLoop(time.Duration)
	StarteRate(uint64)
	Starte()
	GibRate() uint64
	Stoppe()
	Laeuft() bool
}
