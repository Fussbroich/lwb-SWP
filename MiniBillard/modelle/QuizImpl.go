package modelle

import (
	"math/rand"
)

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
