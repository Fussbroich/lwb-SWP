package hilf

import "time"

type Routine interface {
	SetzeAusnahmeHandler(func())
	StarteLoop(time.Duration)
	StarteRate(uint64)
	Starte()
	StarteHier() // Läuft ohne go-Routine
	GibRate() uint64
	GibName() string
	Stoppe()
	Laeuft() bool
}
