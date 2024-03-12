package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	wdPath := filepath.Dir(wd)
	fmt.Println(wdPath)
}
