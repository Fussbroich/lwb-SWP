package modelle

type Quiz interface {
	NaechsteFrage()
	GibAktuelleFrage() QuizFrage
}
