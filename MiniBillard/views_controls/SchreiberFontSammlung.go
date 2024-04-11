package views_controls

import (
	"errors"
	"os"
	"path/filepath"
)

func fontDateipfad(filename string) string {
	fontsDir := "MiniBillard/assets/fontfiles"
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	workDir := filepath.Dir(wd)
	fp := filepath.Join(workDir, fontsDir, filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return fp
}

func LiberationMonoBold(größe int) *schreiber {
	return &schreiber{
		fontdatei:      fontDateipfad("LiberationMono-Bold.ttf"),
		schriftgroesse: größe}
}

func LiberationMonoRegular(größe int) *schreiber {
	return &schreiber{
		fontdatei:      fontDateipfad("LiberationMono-Regular.ttf"),
		schriftgroesse: größe}
}

func LiberationMonoBoldItalic(größe int) *schreiber {
	return &schreiber{
		fontdatei:      fontDateipfad("LiberationMono-BoldItalic.ttf"),
		schriftgroesse: größe}
}
