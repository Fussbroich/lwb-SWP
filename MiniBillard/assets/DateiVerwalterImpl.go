package assets

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
)

func assetDateipfad(dir, filename string) string {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	wdir := filepath.Dir(wd)
	fp := filepath.Join(wdir, dir, filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return fp
}

func KlangDateipfad(filename string) string {
	dir := "MiniBillard/assets/soundfiles"
	return assetDateipfad(dir, filename)
}

func FragenDateipfad(filename string) string {
	dir := "MiniBillard/assets/quizfragen"
	return assetDateipfad(dir, filename)
}

func GibFragenCSV(filename string) [][]string {
	file, err := os.Open(FragenDateipfad(filename))
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

func FontDateipfad(filename string) string {
	dir := "MiniBillard/assets/fontfiles"
	return assetDateipfad(dir, filename)
}
