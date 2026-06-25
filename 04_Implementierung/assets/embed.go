// embed.go — Alle Assets (Fonts, Sounds, Quizfragen) werden zur
// Kompilierzeit ins Binary eingebettet. Kein Dateisystem zur Laufzeit nötig.
// Damit funktioniert das Programm als einzelne .exe und ist WASM-fähig.
package assets

import "embed"

//go:embed fontfiles/*.ttf soundfiles/*.wav quizfragen/*.csv
var assetFS embed.FS

func lese(name string) []byte {
	data, err := assetFS.ReadFile(name)
	if err != nil {
		panic("Asset nicht gefunden: " + name)
	}
	return data
}
