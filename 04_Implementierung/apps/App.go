package apps

// Eine App ist eine grafische Anwendung, die im "unmittelbaren Modus" in einem einzigen
// Fenster läuft. Das bedeutet, dass nach jedem zeitlichen "Tick" (Zeiteinheit) das gesamte
// Fenster mit allen grafischen Elementen (Widgets) neu gezeichnet wird. Die Modelle sind
// Teil der App, und ihr Zustand wird jedesmal neu abgefragt.
//
//	Vor.: Das Grafikpaket gfx muss im GOPATH installiert sein.
type App interface {
	// Setzt den Quitter
	SetzeQuit(func())

	// Die Größe, die das Gfx-Fenster haben muss.
	GibGroesse() (uint16, uint16)

	// Der Titel, den das Gfx-Fenster haben soll.
	GibTitel() string

	// Der App-Loop ruft diese Funktion bei jedem Tick einmal auf.
	Update()

	// Der Zeichen-Loop ruft diese Funktion bei jedem Tick einmal auf.
	// Vor: Gfx Fenster ist offen
	Zeichne()

	// Ein zentraler callback für Maus-Ereignisse
	MausEreignis(uint8, int8, uint16, uint16)

	// Ein zentraler callback für Tastatur-Ereignisse
	TastaturEreignis(uint16, uint8, uint16)
}
