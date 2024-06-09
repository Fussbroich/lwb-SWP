package assets

import (
	"encoding/csv"
	"os"
)

func fragenDateipfad(filename string) string {
	dir := "04_Implementierung/assets/quizfragen"
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

func BeispielFragen() [][]string {
	return gibFragenDaten(fragenDateipfad("BeispielQuiz.csv"))
}

func InformatiksystemeFragen() [][]string {
	return gibFragenDaten(fragenDateipfad("InformatiksystemQuiz.csv"))
}
