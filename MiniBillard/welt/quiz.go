package welt

import "math/rand"

type Quiz interface {
	GibAktuelleFrage() QuizFrage
	GibNächsteFrage() QuizFrage
}

type quiz struct {
	fragen   []QuizFrage
	aktuelle QuizFrage
}

func (q *quiz) GibAktuelleFrage() QuizFrage {
	if q.aktuelle == nil {
		q.aktuelle = q.fragen[rand.Intn(len(q.fragen))]
	}
	return q.aktuelle
}

func (q *quiz) GibNächsteFrage() QuizFrage {
	q.aktuelle = q.fragen[rand.Intn(len(q.fragen))]
	return q.aktuelle
}

func NewAlphabetQuiz() *quiz {
	return &quiz{
		fragen: []QuizFrage{
			&quizfrage{frage: "Wie lautet der erste Buchstabe des Alphabets?",
				antworten: [4]string{"A", "B", "C", "D"},
				richtig:   0},
			&quizfrage{frage: "Wie lautet der zweite Buchstabe des Alphabets?",
				antworten: [4]string{"A", "B", "C", "D"},
				richtig:   1},
			&quizfrage{frage: "Wie lautet der dritte Buchstabe des Alphabets?",
				antworten: [4]string{"A", "B", "C", "D"},
				richtig:   2},
			&quizfrage{frage: "Wie lautet der vierte Buchstabe des Alphabets?",
				antworten: [4]string{"A", "B", "C", "D"},
				richtig:   3},
			&quizfrage{frage: "Wie lautet der fünfte Buchstabe des Alphabets?",
				antworten: [4]string{"D", "E", "F", "G"},
				richtig:   0},
			&quizfrage{frage: "Wie lautet der sechste Buchstabe des Alphabets?",
				antworten: [4]string{"D", "E", "F", "G"},
				richtig:   1},
			&quizfrage{frage: "Wie lautet der siebte Buchstabe des Alphabets?",
				antworten: [4]string{"D", "E", "F", "G"},
				richtig:   2}}}
}

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	IstRichtig(uint8) bool
}

type quizfrage struct {
	frage     string
	richtig   uint8
	antworten [4]string
}

func NewFrage(frage, a1, a2, a3, a4 string, richtig uint8) *quizfrage {
	return &quizfrage{frage: frage, antworten: [4]string{a1, a2, a3, a4}, richtig: richtig}
}

func (f *quizfrage) GibFrage() string { return f.frage }

func (f *quizfrage) GibAntworten() [4]string { return f.antworten }

func (f *quizfrage) IstRichtig(i uint8) bool { return i == f.richtig }
