package views_controls

type FarbSchema interface {
	GetFarbe(string) Farbe
}

type schema struct {
	farben map[string]Farbe
}

var (
	standardfarbschema schema = schema{
		farben: map[string]Farbe{
			"hintergrund": F(225, 232, 236),
			"text":        F(1, 88, 122),
			"anzeige":     Weiss(),
			"infos":       F(240, 255, 255),
			"billardtuch": F(92, 179, 193),
			"diamanten":   F(180, 230, 255),
			"anz_treffer": F(243, 186, 0),
			"anz_fouls":   F(219, 80, 0),
			"quiz":        Weiss(),
			"quiz_a0":     F(155, 155, 0),
			"quiz_a1":     F(255, 255, 0),
			"quiz_a2":     F(0, 255, 255),
			"quiz_a3":     F(255, 0, 255)}}
	darkfarbschema = schema{
		farben: map[string]Farbe{
			"hintergrund": F(25, 32, 36),
			"text":        F(1, 88, 122),
			"anzeige":     F(100, 100, 100),
			"infos":       F(140, 155, 155),
			"billardtuch": F(92, 179, 193),
			"diamanten":   F(180, 230, 255),
			"anz_treffer": F(143, 86, 0),
			"anz_fouls":   F(119, 40, 0),
			"quiz":        Weiss(),
			"quiz_a0":     F(155, 155, 0),
			"quiz_a1":     F(255, 255, 0),
			"quiz_a2":     F(0, 255, 255),
			"quiz_a3":     F(255, 0, 255)}}
	farbschema *schema = &standardfarbschema
)

func StandardFarbSchema() { farbschema = &standardfarbschema }

func DarkFarbSchema() { farbschema = &darkfarbschema }

func (s *schema) GibFarbe(name string) Farbe {
	c, ok := s.farben[name]
	if !ok {
		return Weiss()
	}
	return c
}

func gibFarbe(name string) Farbe {
	if farbschema == nil {
		farbschema = &standardfarbschema
	}
	return farbschema.GibFarbe(name)
}

func Fhintergrund() string { return "hintergrund" }
func Ftext() string        { return "text" }
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
