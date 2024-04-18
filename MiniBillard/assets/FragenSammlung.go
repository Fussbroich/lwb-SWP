package assets

import (
	"encoding/csv"
	"os"
)

func fragenDateipfad(filename string) string {
	dir := "MiniBillard/assets/quizfragen"
	return assetDateipfad(dir, filename)
}

func gibFragenDaten(dateipfad string) [][]string {
	file, err := os.Open(dateipfad)
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
	return rs
}

func BeispielFragenDateipfad() string {
	return fragenDateipfad("BeispielQuiz.csv")
}

func BeispielFragen() [][]string {
	return gibFragenDaten(BeispielFragenDateipfad())
}

func InformatiksystemeFragenDateipfad() string {
	return fragenDateipfad("InformatiksystemQuiz.csv")
}
func InformatiksystemeFragen() [][]string {
	return gibFragenDaten(InformatiksystemeFragenDateipfad())
}
