package welt

import "math/rand"

type Quiz interface {
	NächsteFrage()
	GibAktuelleFrage() QuizFrage
}

type quiz struct {
	fragen   []QuizFrage
	aktuelle QuizFrage
}

func (q *quiz) NächsteFrage() {
	q.aktuelle = q.fragen[rand.Intn(len(q.fragen))]
}

func (q *quiz) GibAktuelleFrage() QuizFrage {
	if q.aktuelle == nil {
		q.aktuelle = q.fragen[rand.Intn(len(q.fragen))]
	}
	return q.aktuelle
}

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	IstRichtig(uint8) bool
}

type quizfrage struct {
	frage            string
	richtig, gewählt uint8
	antworten        [4]string
}

func NewFrage(frage, a1, a2, a3, a4 string, richtig uint8) *quizfrage {
	return &quizfrage{frage: frage, antworten: [4]string{a1, a2, a3, a4}, richtig: richtig}
}

func (f *quizfrage) GibFrage() string { return f.frage }

func (f *quizfrage) GibAntworten() [4]string { return f.antworten }

func (f *quizfrage) IstRichtig(i uint8) bool { return i == f.richtig }
