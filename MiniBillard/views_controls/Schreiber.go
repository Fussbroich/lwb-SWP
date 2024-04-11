package views_controls

type Schreiber interface {
	SetzeSchriftgroesse(int)
	GibSchriftgroesse() int
	Schreibe(x, y uint16, text string)
}
