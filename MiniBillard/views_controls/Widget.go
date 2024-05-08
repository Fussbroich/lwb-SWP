package views_controls

// Ein Widget ist ein darstellbares Objekt mit einem bestimmten Design.
//
// Es hat oft ein Modell zugeordnet, aus dem es Daten abrufen kann.
// Es kann bisweilen auch Mausaktionen enthalten und ausführen.
//
//	Vor.: Das Grafikpaket gfx muss im GOPATH installiert sein und
//	es muss ein Fenster geöffnet sein, bevor Zeichnen aufgerufen wird.
//
//	Die verschiedensten Konstruktoren erzeugen das jeweils passende Widget für eine Aufgabe.
type Widget interface {
	// Konstruktormethode/Setter zum Verändern der Anzeige-Position und -Größe im Grafik-Fenster
	//	Vor.: keine
	//	Eff.: die Einstellung ändern sich
	SetzeKoordinaten(uint16, uint16, uint16, uint16)
	// Konstruktormethode/Setter zum Verändern der Anzeige-Farben im Fenster
	//	Vor.: keine
	//	Eff.: die Einstellung ändern sich
	SetzeFarben(FarbID, FarbID)
	// Konstruktormethode/Setter zum Verändern der Anzeige-Transparenz im Fenster
	//	Vor.: keine
	//	Eff.: die Einstellung ändern sich
	SetzeTransparenz(uint8)
	// Konstruktormethode/Setter zum Verändern des Eckradius
	//	Vor.: keine
	//	Eff.: die Einstellung ändern sich
	SetzeEckradius(uint16)
	// Getter zum Auslesen der Startecke des Widgets im Fenster
	//	Vor.: keine
	//	Erg.: die obere linke Ecke ist geliefert
	GibStartkoordinaten() (uint16, uint16)
	// Getter zum Auslesen der aktuellen Größe des Widgets im Fenster
	//	Vor.: keine
	//	Erg.: die Breite und die Höhe sind geliefert
	GibGroesse() (uint16, uint16)
	// darstellen:
	// Lädt die Anzeigefarben neu aus dem aktiven Farb-Schema.
	//	Vor.: keine
	//	Eff.: Die Farben werden nach Änderung des Schemas ans aktive Schema angepasst
	LadeFarben()
	// darstellen:
	// Manche Widgets sind sehr kompliziert und belasten den gfx-Server -
	// hiermit kann man die Darstellung etwas vereinfachen, falls das Widget es unterstützt.
	//	Vor.: keine
	//	Eff.: Die Darstellung wird ggf. vereinfacht.
	SetzeSchlicht()
	// Die Darstellungsmethode
	//	Vor.: keine
	//	Eff.: Das Widget stellt sich im gfx-Fenster dar, falls es aktiv ist.
	//	Die Oberklasse zeichnet den Hintergrund und ergänzende Klassen zeichnen neue Inhalte.
	//	Aufrufer (ergänzende Klassen) - die den Inhalt ergänzen wollen, müssen
	//	ihren Inhalt erst danach zeichnen, sonst wird
	//	dieser ggf. überdeckt. Transparenz wird beachtet.
	//
	//	Hinweise: Bei Widgets, die immer sichtbar sind, kann das Verhalten geändert werden,
	//	so dass sie sich nur dann neu zeichnen, falls sich Daten geändert haben.
	Zeichne()
	// Die Datenlademethode
	//	Vor.: keine
	//	Eff.: Das Widget lädt seine anzuzeigenden Daten - falls vorhanden - aus seinem zugeordneten Modell neu.
	Update()
	// für besondere Fälle
	// (beispielsweise, wenn der Rand besser zu sehen sein soll)
	ZeichneOffset(uint16)
	// zeichnet einen Rand um das Widget
	ZeichneRand()
	// Testzwecke - zeichnet einen roten Rahmen um das Widget
	ZeichneLayout()
	// aktivieren und deaktivieren:
	// Nur aktive Widgets zeichnen sich auch. Es lassen sich so "Fenster" effektiv ein- und ausblenden.
	//	Vor.: keine
	//	Erg.: Liefert wahr, falls das Widget aktiv ist.
	IstAktiv() bool
	// Schalte Darstellung an bzw. aus.
	DarstellenAnAus()
	// Aktiviert das Widget. Nur aktive Widgets zeichnen sich auch.
	//	Vor.: keine
	//	Eff.: Das Widget zeigt sich an.
	Einblenden()
	// Deaktiviert das Widget. Nur aktive Widgets zeichnen sich auch.
	//	Vor.: keine
	//	Eff.: Das Widget zeigt sich nicht mehr an.
	Ausblenden()
	// Maussteuerung: Stellt fest, ob die Maus innerhalb des Widgets ist.
	//	Vor.: keine
	//	Erg.: Liefert wahr, falls die gegebenen Koordinate im Widget liegt.
	ImFenster(uint16, uint16) bool
	// Maussteuerung: Führt eine hinterlegte Aktion für die gegebene Koordinate aus.
	//	Vor.: keine
	//	Eff.: Falls eine Aktion hinterlegt ist, wird sie ausgeführt, sonst passiert nichts.
	MausklickBei(uint16, uint16)
	// Maussteuerung: Führt eine hinterlegte Aktion für die gegebene Koordinate aus.
	//	Vor.: keine
	//	Eff.: Falls eine Aktion hinterlegt ist, wird sie ausgeführt, sonst passiert nichts.
	MausBei(uint16, uint16)
	// Maussteuerung: Führt eine hinterlegte Aktion aus.
	//	Vor.: keine
	//	Eff.: Falls eine Aktion hinterlegt ist, wird sie ausgeführt, sonst passiert nichts.
	MausScrolltHoch()
	// Maussteuerung: Führt eine hinterlegte Aktion aus.
	//	Vor.: keine
	//	Eff.: Falls eine Aktion hinterlegt ist, wird sie ausgeführt, sonst passiert nichts.
	MausScrolltRunter()
}
