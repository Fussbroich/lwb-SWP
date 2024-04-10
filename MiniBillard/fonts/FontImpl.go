package fonts

type font struct {
	name         string
	dateipfad    string
	schriftgröße int
}

func (f *font) GibName() string {
	return f.name
}

func (f *font) GibDateipfad() string {
	return f.dateipfad
}

func (f *font) GibSchriftgröße() int {
	return f.schriftgröße
}

func (f *font) SetzeSchriftgröße(größe int) {
	f.schriftgröße = größe
}
