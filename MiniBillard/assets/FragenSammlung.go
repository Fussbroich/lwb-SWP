package assets

import (
	"encoding/csv"
	"os"
)

func fragenDateipfad(filename string) string {
	dir := "MiniBillard/assets/quizfragen"
	return assetDateipfad(dir, filename)
}

func gibFragenDaten(filename string) [][]string {
	file, err := os.Open(fragenDateipfad(filename))
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

func BeispielFragen() [][]string {
	return gibFragenDaten("BeispielQuiz.csv")
}

func InformatiksystemeFragen() [][]string {
	return gibFragenDaten("InformatiksystemQuiz.csv")
}
