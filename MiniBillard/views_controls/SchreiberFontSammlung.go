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

var (
	liberationMonoBoldFont       = fontDateipfad("LiberationMono-Bold.ttf")
	liberationMonoRegularFont    = fontDateipfad("LiberationMono-Regular.ttf")
	liberationMonoBoldItalicFont = fontDateipfad("LiberationMono-BoldItalic.ttf")
)

func (f *widget) LiberationMonoBoldSchreiber() *schreiber {
	return &schreiber{fontdatei: liberationMonoBoldFont,
		schriftgroesse: 24}
}

func (f *widget) LiberationMonoRegularSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      liberationMonoRegularFont,
		schriftgroesse: 24}
}

func (f *widget) LiberationMonoBoldItalicSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      liberationMonoBoldItalicFont,
		schriftgroesse: 24}
}
