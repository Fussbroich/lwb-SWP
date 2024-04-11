package modelle

type quizfrage struct {
	frage             string
	richtig, gewaehlt int
	antworten         [4]string
}

func NewQuizFrage(frage, a1, a2, a3, a4 string, richtig int) *quizfrage {
	return &quizfrage{frage: frage, antworten: [4]string{a1, a2, a3, a4}, richtig: richtig, gewaehlt: -1}
}

func (f *quizfrage) GibFrage() string { return f.frage }

func (f *quizfrage) GibAntworten() [4]string { return f.antworten }

func (f *quizfrage) Gewaehlt(i int) { f.gewaehlt = i }

func (f *quizfrage) RichtigBeantwortet() bool { return f.richtig == f.gewaehlt }
