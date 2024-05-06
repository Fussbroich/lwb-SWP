package views_controls

type schema struct {
	farben map[FarbID]Farbe
}

type FarbID uint8

const (
	Frot FarbID = iota
	Fweiss
	Fschwarz
	Fhintergrund
	Ftext
	Fbande
	Fanzeige
	Fbillardtuch
	Finfos
	Fdiamanten
	FanzTreffer
	FanzFouls
	Fquiz
	FquizA0
	FquizA1
	FquizA2
	FquizA3
)

var (
	standardfarbschema schema = schema{
		farben: map[FarbID]Farbe{
			Frot:         Rot(),
			Fweiss:       Weiss(),
			Fschwarz:     Schwarz(),
			Fhintergrund: F(225, 232, 236),
			Ftext:        F(1, 88, 122),
			Fbande:       F(1, 88, 122),
			Fanzeige:     Weiss(),
			Finfos:       F(240, 255, 255),
			Fbillardtuch: F(92, 179, 193),
			Fdiamanten:   F(180, 230, 255),
			FanzTreffer:  F(243, 186, 0),
			FanzFouls:    F(219, 80, 0),
			Fquiz:        F(225, 232, 236),
			FquizA0:      F(243, 186, 0),
			FquizA1:      F(92, 179, 193),
			FquizA2:      F(92, 179, 193),
			FquizA3:      F(243, 186, 0)}}
	darkfarbschema = schema{
		farben: map[FarbID]Farbe{
			Frot:         Rot(),
			Fweiss:       Weiss(),
			Fschwarz:     Schwarz(),
			Fhintergrund: F(25, 32, 36),
			Ftext:        F(125, 170, 195),
			Fbande:       F(1, 88, 122),
			Fanzeige:     F(25, 32, 36),
			Finfos:       F(140, 155, 155),
			Fbillardtuch: F(92, 179, 193),
			Fdiamanten:   F(180, 230, 255),
			FanzTreffer:  F(143, 86, 0),
			FanzFouls:    F(119, 40, 0),
			Fquiz:        F(25, 32, 36),
			FquizA0:      F(80, 64, 19),
			FquizA1:      F(40, 61, 65),
			FquizA2:      F(40, 61, 65),
			FquizA3:      F(80, 64, 19)}}
	farbschema *schema = &standardfarbschema
)

func SetzeStandardFarbSchema() { farbschema = &standardfarbschema }

func SetzeDarkFarbSchema() { farbschema = &darkfarbschema }

func (s *schema) GibFarbe(id FarbID) Farbe {
	c, ok := s.farben[id]
	if !ok {
		return Rot()
	}
	return c
}

func gibFarbe(id FarbID) Farbe {
	if farbschema == nil {
		farbschema = &standardfarbschema
	}
	return farbschema.GibFarbe(id)
}
