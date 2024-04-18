package modelle

import (
	"strconv"

	"../assets"
)

func quizFragenCSV(fn string) (fragen []QuizFrage) {
	rs := assets.GibFragenCSV(fn)
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
		fragen = append(fragen, f)
	}
	return
}

func newQuizCSV(fn string) *quiz {
	return &quiz{fragen: quizFragenCSV(fn)}
}

func NewBeispielQuiz() *quiz {
	return newQuizCSV("BeispielQuiz.csv")
}

func NewQuizInformatiksysteme() *quiz {
	return newQuizCSV("InformatiksystemQuiz.csv")
}
