package views_controls

var (
	farbschema         *(map[string]Farbe)
	kugelpalette       *([16]Farbe)
	standardfarbschema map[string]Farbe = map[string]Farbe{
		"hintergrund": F(225, 232, 236),
		"text":        F(1, 88, 122),
		"anzeige":     F(255, 255, 255),
		"infos":       F(240, 255, 255),
		"billardtuch": F(92, 179, 193),
		"diamanten":   F(180, 230, 255),
		"anz_treffer": F(243, 186, 0),
		"anz_fouls":   F(219, 80, 0)}
	darkfarbschema map[string]Farbe = map[string]Farbe{
		"hintergrund": F(25, 32, 36),
		"text":        F(1, 88, 122),
		"anzeige":     F(100, 100, 100),
		"infos":       F(140, 155, 155),
		"billardtuch": F(92, 179, 193),
		"diamanten":   F(180, 230, 255),
		"anz_treffer": F(143, 86, 0),
		"anz_fouls":   F(119, 40, 0)}

	standardPoolPalette [16]Farbe = [16]Farbe{
		F(252, 253, 242), // weiß
		F(255, 201, 78),  // gelb
		F(34, 88, 175),   // blau
		F(249, 73, 68),   // hellrot
		F(84, 73, 149),   // violett
		F(255, 139, 33),  // orange
		F(47, 159, 52),   // grün
		F(155, 53, 30),   // dunkelrot
		F(48, 49, 54),    // schwarz
		F(255, 201, 78),  // gelb
		F(34, 88, 175),   // blau
		F(249, 73, 68),   // hellrot
		F(84, 73, 149),   // violett
		F(255, 139, 33),  // orange
		F(47, 159, 52),   // grün
		F(155, 53, 30)}   // dunkelrot

	englishPoolPalette [16]Farbe = [16]Farbe{
		F(252, 253, 242), // weiß
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(48, 49, 54),    // schwarz
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78)}  // gelb

)

func StandardFarbSchema() { farbschema = &standardfarbschema }

func DarkFarbSchema() { farbschema = &darkfarbschema }

func getFarbe(name string) Farbe {
	if farbschema == nil {
		farbschema = &standardfarbschema
	}
	return (*farbschema)[name]
}

func StandardKugelPalette() { kugelpalette = &standardPoolPalette }

func EnglishKugelPalette() { kugelpalette = &englishPoolPalette }

func KugelFarbe(wert uint8) Farbe {
	if kugelpalette == nil {
		kugelpalette = &standardPoolPalette
	}
	if int(wert) >= len(kugelpalette) {
		return Schwarz()
	}
	return (*kugelpalette)[wert]
}

func Fhintergrund() Farbe { return getFarbe("hintergrund") }

func Ftext() Farbe { return getFarbe("text") }

func Fanzeige() Farbe { return getFarbe("anzeige") }

func Fbillardtuch() Farbe { return getFarbe("billardtuch") }

func Fdiamanten() Farbe { return getFarbe("diamanten") }

func FanzTreffer() Farbe { return getFarbe("anz_treffer") }

func FanzFouls() Farbe { return getFarbe("anz_fouls") }

func Finfos() Farbe {
	return &rgb{r: 240, g: 255, b: 255}
}

func Weiß() Farbe {
	return &rgb{r: 255, g: 255, b: 255}
}

func Schwarz() Farbe {
	return &rgb{}
}
