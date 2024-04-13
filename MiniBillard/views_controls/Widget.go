package views_controls

type Widget interface {
	Zeichne()
	ZeichneRand()
	ZeichneLayout()
	SetzeKoordinaten(uint16, uint16, uint16, uint16)
	SetzeFarben(string, string)
	LadeFarben()
	IstAktiv() bool
	AktivAnAus()
	SetzeAktiv()
	SetzeInAktiv()
	SetzeTransparenz(uint8)
	SetzeEckradius(uint16)
	GibStartkoordinaten() (uint16, uint16)
	GibGroesse() (uint16, uint16)
	ImFenster(mausX, mausY uint16) bool
	MausklickBei(mausX, mausY uint16)
}
