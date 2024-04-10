package modelle

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	Gewählt(int)
	RichtigBeantwortet() bool
}
