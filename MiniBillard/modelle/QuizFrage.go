package modelle

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	IstRichtig(int) bool
}
