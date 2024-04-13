package views_controls

type FensterZeichner interface {
	Starte()
	Stoppe()
	SetzeFensterHintergrund(Widget)
	SetzeWidgets(...Widget)
	LayoutAnAus()
	DarkmodeAnAus()
	Ueberblende(Widget)
	UeberblendeText(string, string, string, uint8)
	UeberblendeAus()
}
