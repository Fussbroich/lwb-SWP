package modelle

// Eine QuizFrage ist ein Container für eine Frage mit 4 Antwortmöglichkeiten.
//
//	NewQuizFrage(frage, a1, a2, a3, a4 string, richtig int) erzeugt eine Frage.
type QuizFrage interface {
	// Getter für den Fragetext
	//	Vor.: keine
	//	Erg.: der Fragetext
	GibFrage() string
	// Getter für die 4 Antwortmöglichkeiten.
	//	Vor.: keine
	//	Erg.: die Antworten
	GibAntworten() [4]string
	// Prüfe, ob eine bestimmte Antwort die richtige ist.
	//	Vor.: keine
	//	Eff.: wahr, falls die Antwort stimmt, sonst falsch.
	IstRichtig(int) bool
}
