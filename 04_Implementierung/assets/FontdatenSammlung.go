package assets

func fontDateipfad(filename string) string {
	dir := "04_Implementierung/assets/fontfiles"
	return assetDateipfad(dir, filename)
}

func MonoBoldFontDateipfad() string { return fontDateipfad("LiberationMono-Bold.ttf") }

func MonoRegularFontDateipfad() string { return fontDateipfad("LiberationMono-Regular.ttf") }

func MonoBoldItalicFontDateipfad() string { return fontDateipfad("LiberationMono-BoldItalic.ttf") }
