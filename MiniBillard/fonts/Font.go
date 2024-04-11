package fonts

type Font interface {
	GibName() string
	GibDateipfad() string
	SetzeSchriftgroesse(int)
	GibSchriftgroesse() int
}
