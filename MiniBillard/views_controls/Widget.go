package views_controls

type Widget interface {
	Zeichne()
	ZeichneRand()
	ZeichneLayout()
	SetzeKoordinaten(uint16, uint16, uint16, uint16)
	SetzeFarben(Farbe, Farbe)
	SetzeTransparenz(uint8)
	SetzeEckradius(uint16)
	GibStartkoordinaten() (uint16, uint16)
	GibGröße() (uint16, uint16)
	ImFenster(mausX, mausY uint16) bool
	MausklickBei(mausX, mausY uint16)
}
