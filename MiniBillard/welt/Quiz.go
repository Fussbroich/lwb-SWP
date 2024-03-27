package welt

type Quiz interface {
	NächsteFrage()
	GibAktuelleFrage() QuizFrage
}

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	Gewählt(int)
	RichtigBeantwortet() bool
}
