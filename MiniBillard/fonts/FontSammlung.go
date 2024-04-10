package fonts

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

func LiberationMonoBold(größe int) *font {
	return &font{
		name:         "LiberationMono bold",
		dateipfad:    fontDateipfad("LiberationMono-Bold.ttf"),
		schriftgröße: größe}
}

func LiberationMonoRegular(größe int) *font {
	return &font{
		name:         "LiberationMono regular",
		dateipfad:    fontDateipfad("LiberationMono-Regular.ttf"),
		schriftgröße: größe}
}

func LiberationMonoBoldItalic(größe int) *font {
	return &font{
		name:         "LiberationMono bold, italic",
		dateipfad:    fontDateipfad("LiberationMono-BoldItalic.ttf"),
		schriftgröße: größe}
}
