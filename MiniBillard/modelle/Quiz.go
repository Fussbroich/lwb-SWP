package modelle

// TODO Ein Quiz ist ...
//
//	New...()
type Quiz interface {
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	NaechsteFrage()
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	GibAktuelleFrage() QuizFrage
}
