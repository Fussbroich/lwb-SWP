package hilf

import "time"

type Routine interface {
	SetzeAusnahmeHandler(func())
	StarteLoop(time.Duration)
	StarteRate(uint64)
	Starte()
	GibRate() uint64
	GibName() string
	Stoppe()
	Laeuft() bool
	Einmal()
}
