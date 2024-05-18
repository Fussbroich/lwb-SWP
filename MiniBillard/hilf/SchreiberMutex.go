package hilf

// Zur Absicherung der Veränderung von Spielmodellen - falsche Anzeigen werden kurzfristig akzeptiert
type SchreiberMutex interface {
	// Viele Leser dürfen jederzeit "rein"
	LeserEin()
	LeserAus()
	GibAnzLeser() uint // Testzwecke

	// nur *Ein* Schreiber darf rein, selbst wenn Leser drin sind
	SchreiberEin()
	SchreiberAus()
}
