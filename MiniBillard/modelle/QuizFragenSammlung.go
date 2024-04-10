package modelle

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

func assetDateipfad(filename string) (fp string) {
	fragenDir := "MiniBillard/assets/quizfragen"
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

func QuizFragenCSV(fn string) (fragen []QuizFrage) {
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
		fragen = append(fragen, f)
	}
	return
}

func NewQuizCSV(fn string) *quiz {
	return &quiz{fragen: QuizFragenCSV(fn)}
}
