package assets

import (
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
