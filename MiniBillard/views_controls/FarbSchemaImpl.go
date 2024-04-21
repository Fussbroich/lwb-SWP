package views_controls

type FarbSchema interface {
	// Getter für eine bestimmte Farbe aus einem Schema.
	// Hinweis: Schema-Farben haben Namen (Strings). Wird eine unbekannte
	// Farbe angefordert, so gibt die Methode rot zurück. Für die bekannten
	// Farbnamen der hier vorinstallierten Schemata gibt es die Spezialfunktion F...()
	GibFarbe(string) Farbe
}

type schema struct {
	farben map[string]Farbe
}

var (
	standardfarbschema schema = schema{
		farben: map[string]Farbe{
			"rot":         Rot(),
			"weiss":       Weiss(),
			"schwarz":     Schwarz(),
			"hintergrund": F(225, 232, 236),
			"text":        F(1, 88, 122),
			"bande":       F(1, 88, 122),
			"anzeige":     Weiss(),
			"infos":       F(240, 255, 255),
			"billardtuch": F(92, 179, 193),
			"diamanten":   F(180, 230, 255),
			"anz_treffer": F(243, 186, 0),
			"anz_fouls":   F(219, 80, 0),
			"quiz":        F(225, 232, 236),
			"quiz_a0":     F(243, 186, 0),
			"quiz_a1":     F(92, 179, 193),
			"quiz_a2":     F(92, 179, 193),
			"quiz_a3":     F(243, 186, 0)}}
	darkfarbschema = schema{
		farben: map[string]Farbe{
			"rot":         Rot(),
			"weiss":       Weiss(),
			"schwarz":     Schwarz(),
			"hintergrund": F(25, 32, 36),
			"text":        F(125, 170, 195),
			"bande":       F(1, 88, 122),
			"anzeige":     F(25, 32, 36),
			"infos":       F(140, 155, 155),
			"billardtuch": F(92, 179, 193),
			"diamanten":   F(180, 230, 255),
			"anz_treffer": F(143, 86, 0),
			"anz_fouls":   F(119, 40, 0),
			"quiz":        F(25, 32, 36),
			"quiz_a0":     F(80, 64, 19),
			"quiz_a1":     F(40, 61, 65),
			"quiz_a2":     F(40, 61, 65),
			"quiz_a3":     F(80, 64, 19)}}
	farbschema *schema = &standardfarbschema
)

func SetzeStandardFarbSchema() { farbschema = &standardfarbschema }

func SetzeDarkFarbSchema() { farbschema = &darkfarbschema }

func (s *schema) GibFarbe(name string) Farbe {
	c, ok := s.farben[name]
	if !ok {
		return Rot()
	}
	return c
}

func gibFarbe(name string) Farbe {
	if farbschema == nil {
		farbschema = &standardfarbschema
	}
	return farbschema.GibFarbe(name)
}

func FRot() string         { return "rot" }
func FWeiss() string       { return "weiss" }
func FSchwarz() string     { return "schwarz" }
func Fhintergrund() string { return "hintergrund" }
func Ftext() string        { return "text" }
func Fbande() string       { return "bande" }
func Fanzeige() string     { return "anzeige" }
func Fbillardtuch() string { return "billardtuch" }
func Finfos() string       { return "infos" }
func Fdiamanten() string   { return "diamanten" }
func FanzTreffer() string  { return "anz_treffer" }
func FanzFouls() string    { return "anz_fouls" }
func Fquiz() string        { return "quiz" }
func FquizA0() string      { return "quiz_a0" }
func FquizA1() string      { return "quiz_a1" }
func FquizA2() string      { return "quiz_a2" }
func FquizA3() string      { return "quiz_a3" }

// Zusätzliche praktische Farben
func Weiss() Farbe {
	return &rgb{r: 255, g: 255, b: 255}
}

func Schwarz() Farbe {
	return &rgb{}
}

func Rot() Farbe {
	return &rgb{r: 255}
}
