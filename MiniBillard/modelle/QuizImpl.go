package modelle

import (
	"math/rand"
	"strconv"

	"../assets"
)

type quiz struct {
	fragen   []QuizFrage
	aktuelle QuizFrage
}

func newQuiz(daten [][]string) *quiz {
	var fragen []QuizFrage
	for _, r := range daten {
		if len(r) != 6 {
			panic("Falsches Fragen-Format")
		}
		i, err := strconv.Atoi(r[5])
		if err != nil || i < 0 || i > 3 {
			panic("Falsches Fragen-Format")
		}
		f := NewQuizFrage(
			r[0], r[1], r[2], r[3], r[4],
			i)
		fragen = append(fragen, f)
	}
	return &quiz{fragen: fragen}
}

func NewBeispielQuiz() *quiz {
	return newQuiz(assets.BeispielFragen())
}

func NewQuizInformatiksysteme() *quiz {
	return newQuiz(assets.InformatiksystemeFragen())
}

func (q *quiz) NaechsteFrage() {
	q.aktuelle = q.fragen[rand.Intn(len(q.fragen))]
}

func (q *quiz) GibAktuelleFrage() QuizFrage {
	if q.aktuelle == nil {
		q.aktuelle = q.fragen[rand.Intn(len(q.fragen))]
	}
	return q.aktuelle
}
