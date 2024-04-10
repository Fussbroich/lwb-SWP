package views_controls

type FensterZeichner interface {
	Starte()
	Stoppe()
	ZeigeLayout()
	Überblende(Widget)
	ÜberblendeText(string, Farbe, Farbe, uint8)
	ÜberblendeAus()
}
