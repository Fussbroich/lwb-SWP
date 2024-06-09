package modelle

// Ein Quiz ist eine Sammlung von Fragen, die man aufrufen kann.
// Dabei bestimmt die Implementierung, in welcher Reihenfolge die Fragen präsentiert werden.
//
//	NewBeispielQuiz() - erzeugt ein Quiz mit sinnlosen Testfragen
//	NewQuizInformatiksysteme() - erzeugt ein Quiz mit Fragen zur Computertechnik.
type Quiz interface {
	// Schaltet das Quiz auf die nächste Frage, die "dran" ist.
	//
	//	Vor.: keine
	//	Eff.: Die nächste Frage ist gesetzt.
	NaechsteFrage()

	// Liefert die aktuell gesetzte Frage.
	//
	//	Vor.: keine
	//	Erg.: Die aktuelle Frage.
	GibAktuelleFrage() QuizFrage
}
