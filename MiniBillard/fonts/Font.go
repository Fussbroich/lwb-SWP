package fonts

type Font interface {
	GibName() string
	GibDateipfad() string
	SetzeSchriftgröße(int)
	GibSchriftgröße() int
}
