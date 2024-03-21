package welt

import (
	"encoding/csv"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

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

func (q *quiz) Antwort(i int) {
	q.aktuelle.Gewählt(i)
}

type QuizFrage interface {
	GibFrage() string
	GibAntworten() [4]string
	Gewählt(int)
	RichtigBeantwortet() bool
}

type quizfrage struct {
	frage            string
	richtig, gewählt int
	antworten        [4]string
}

func NewQuizFrage(frage, a1, a2, a3, a4 string, richtig int) *quizfrage {
	return &quizfrage{frage: frage, antworten: [4]string{a1, a2, a3, a4}, richtig: richtig, gewählt: -1}
}

func (f *quizfrage) GibFrage() string { return f.frage }

func (f *quizfrage) GibAntworten() [4]string { return f.antworten }

func (f *quizfrage) Gewählt(i int) { f.gewählt = i }

func (f *quizfrage) RichtigBeantwortet() bool { return f.richtig == f.gewählt }

func assetDateipfad(filename string) (fp string) {
	fragenDir := "MiniBillard/welt"
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wdir := filepath.Dir(wd)
	fp = filepath.Join(wdir, fragenDir, filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return
}

func NewQuizCSV(fn string) *quiz {
	q := quiz{}
	q.fragen = []QuizFrage{}
	file, err := os.Open(assetDateipfad(fn))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ';'

	rs, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, r := range rs {
		if len(r) != 6 {
			panic("Falsches CSV-Format")
		}
		i, err := strconv.Atoi(r[5])
		if err != nil || i < 0 || i > 3 {
			panic("Falsches CSV-Format")
		}
		f := NewQuizFrage(
			r[0], r[1], r[2], r[3], r[4],
			i)
		q.fragen = append(q.fragen, f)
	}

	return &q
}
