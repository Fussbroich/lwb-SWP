package views_controls

type KugelPalette interface {
	GibFarbe(uint8) Farbe
}

// Farben der Kugeln
type palette struct {
	farben [16]Farbe
}

var (
	standardPoolPalette palette = palette{
		farben: [16]Farbe{
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
			F(155, 53, 30)}}  // dunkelrot
	englishPoolPalette palette = palette{
		farben: [16]Farbe{
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
			F(255, 201, 78)}} // gelb
	kugelpalette *palette = &standardPoolPalette
)

func StandardKugelPalette() { kugelpalette = &standardPoolPalette }

func EnglishKugelPalette() { kugelpalette = &englishPoolPalette }

func (p *palette) GibFarbe(wert uint8) Farbe {
	if int(wert) >= len(p.farben) {
		return Schwarz()
	}
	return p.farben[wert]
}

func gibKugelFarbe(wert uint8) Farbe {
	if kugelpalette == nil {
		kugelpalette = &standardPoolPalette
	}
	return kugelpalette.GibFarbe(wert)
}
