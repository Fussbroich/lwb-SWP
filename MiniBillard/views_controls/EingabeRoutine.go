package views_controls

type EingabeRoutine interface {
	StarteRate(uint64)
	Starte()
	Stoppe()
}
