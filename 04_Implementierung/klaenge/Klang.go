package klaenge

// Ein Wrapper-Objekt für einen Klang - üblicherweise
// eine wav-Datei mit praktischen Methoden zum Abspielen.
type Klang interface {
	// Spielt den Klang einmal ab.
	Play()
	// Spielt den Klang in Endlosschleife.
	StarteLoop()
	// Stoppt die Endlosschleife.
	// Hinweis: Einen laufenden Klang kann man so nicht abbrechen.
	// Lediglich die Wiederholungs-Schleife wird abgebrochen.
	Stoppe()
}
