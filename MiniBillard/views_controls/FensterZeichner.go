package views_controls

type FensterZeichner interface {
	Starte()
	Stoppe()
	ZeigeLayout()
	Ueberblende(Widget)
	UeberblendeText(string, Farbe, Farbe, uint8)
	UeberblendeAus()
}
