package modelle

type quizfrage struct {
	frage     string
	richtig   int
	antworten [4]string
}

func NewQuizFrage(frage, a1, a2, a3, a4 string, richtig int) *quizfrage {
	return &quizfrage{frage: frage, antworten: [4]string{a1, a2, a3, a4}, richtig: richtig}
}

func (f *quizfrage) GibFrage() string { return f.frage }

func (f *quizfrage) GibAntworten() [4]string { return f.antworten }

func (f *quizfrage) IstRichtig(i int) bool { return f.richtig == i }
