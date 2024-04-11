package views_controls

type EingabeProzess interface {
	StarteRate(uint64)
	Starte()
	Stoppe()
}
