package views_controls

type Widget interface {
	// Konstruktormethoden
	SetzeKoordinaten(uint16, uint16, uint16, uint16)
	SetzeFarben(string, string)
	LadeFarben()
	SetzeTransparenz(uint8)
	SetzeEckradius(uint16)
	GibStartkoordinaten() (uint16, uint16)
	GibGroesse() (uint16, uint16)
	// darstellen
	Zeichne()
	ZeichneRand()
	ZeichneLayout()
	// aktivieren und deaktivieren
	IstAktiv() bool
	DarstellenAnAus()
	Einblenden()
	Ausblenden()
	// Maussteuerung
	ImFenster(uint16, uint16) bool
	MausklickBei(uint16, uint16)
	MausBei(uint16, uint16)
	MausScrolltHoch()
	MausScrolltRunter()
}
