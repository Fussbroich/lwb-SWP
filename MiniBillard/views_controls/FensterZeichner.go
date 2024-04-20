package views_controls

// Ein FensterZeichner (Renderer) verwaltet alle Widgets, die zu einer App gehören, ruft
// ihre Darstellungs-Methode auf und steuert die Bildwiederholung. Im "unmittelbaren Modus"
// werden alle Widgets in einem regelmäßigen Takt der Reihe nach in ein einziges Fenster
// gezeichnet.
//
//	Vor.: Das Grafikpaket gfx muss im GOPATH installiert sein.
//
//	NewFensterZeichner() erzeugt ein leeres Objekt.
type FensterZeichner interface {
	// Startet das wiederholte zeichnen der Widgets, sobald ein gfx-Fenster geöffnet wird.
	//
	//	Vor.: keine
	//	Eff.: Die Widgets stellen sich dar.
	Starte()
	// Stoppt das wiederholte zeichnen der Widgets.
	//
	//	Vor.: keine
	//	Eff.: Die Darstellung ist eingefroren.
	Stoppe()
	// Setzt den Hintergrund.
	//
	// Der Hintergrund ist auch nur ein Widget und wird immer zuerst (zuunterst) gezeichnet.
	//	Vor.: keine
	//	Eff.: Der Hintergrund ändert sich.
	SetzeFensterHintergrund(Widget)
	// Setzt die Widgets, die gezeichnet werden sollen.
	//
	// Die Methode ruft man normalerweise nur beim Erzeugen einer App auf.
	// Widgets lassen sich zur Laufzeit ein- und ausblenden, ohne sie auszutauschen.
	SetzeWidgets(...Widget)
	// Ruft die Widgets dazu auf, ihr Layout mit anzuzeigen
	LayoutAnAus()
	// Setzt das dunkle Farbschema und ruft die Widgets zum Laden der Farben auf.
	DarkmodeAnAus()
	// Spontanes ueberblenden eines weiteren Widgets über das Fenster.
	//
	//	Vor.: keine
	//	Eff.: das übergebene Widget wir immer als letztes (zuoberst) angezeigt
	Ueberblende(Widget)
	// Spontanes erzeugen und überblenden eines Text-Widgets über das Fenster.
	// Für Begrüßungen/Abschieds-Texte oder sonstiges, das nicht direkt zur App gehört.
	//
	//	Vor.: keine
	//	Eff.: das neue Widget wir immer als letztes (zuoberst) angezeigt
	UeberblendeText(string, string, string, uint8)
	// Entfernt das ueberblendete Widget wieder.
	//
	//	Vor.: keine
	//	Eff.: Es wird nichts mehr überblendet.
	UeberblendeAus()
}
