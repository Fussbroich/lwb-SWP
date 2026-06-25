package assets

import (
	"encoding/csv"
	"strings"
)

func leseCSV(name string) [][]string {
	reader := csv.NewReader(strings.NewReader(string(lese(name))))
	reader.Comma = ';'
	rs, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return rs
}

func BeispielFragen() [][]string         { return leseCSV("quizfragen/BeispielQuiz.csv") }
func InformatiksystemeFragen() [][]string { return leseCSV("quizfragen/InformatiksystemQuiz.csv") }
