package welt

import "math/rand"

type Quiz interface {
	NächsteFrage()
	GibAktuelleFrage() QuizFrage
	Antwort(int)
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

func (q *quiz) Antwort(i int) {
	q.aktuelle.Gewählt(i)
}

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	Gewählt(int)
}

type quizfrage struct {
	frage            string
	richtig, gewählt int
	antworten        [4]string
}

func NewQuizFrage(frage, a1, a2, a3, a4 string, richtig int) *quizfrage {
	return &quizfrage{frage: frage, antworten: [4]string{a1, a2, a3, a4}, richtig: richtig}
}

func (f *quizfrage) GibFrage() string { return f.frage }

func (f *quizfrage) GibAntworten() [4]string { return f.antworten }

func (f *quizfrage) Gewählt(i int) { f.gewählt = i }
