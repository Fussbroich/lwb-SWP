package modelle

type Quiz interface {
	NächsteFrage()
	GibAktuelleFrage() QuizFrage
}
