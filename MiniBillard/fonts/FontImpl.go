package fonts

type font struct {
	name           string
	dateipfad      string
	schriftgroesse int
}

func (f *font) GibName() string {
	return f.name
}

func (f *font) GibDateipfad() string {
	return f.dateipfad
}

func (f *font) GibSchriftgroesse() int {
	return f.schriftgroesse
}

func (f *font) SetzeSchriftgroesse(größe int) {
	f.schriftgroesse = größe
}
