// embed.go — Alle Assets (Fonts, Sounds, Quizfragen) werden zur
// Kompilierzeit ins Binary eingebettet. Kein Dateisystem zur Laufzeit nötig.
// Damit funktioniert das Programm als einzelne .exe und ist WASM-fähig.
package assets

import (
	"embed"
	"io"
)

//go:embed fontfiles/*.ttf soundfiles/*.wav quizfragen/*.csv
var assetFS embed.FS

// lese liefert eine Kopie der eingebetteten Datei als []byte.
// Geeignet für kleine Assets (Fonts, Effekt-Sounds, CSV).
func lese(name string) []byte {
	data, err := assetFS.ReadFile(name)
	if err != nil {
		panic("Asset nicht gefunden: " + name)
	}
	return data
}

// OeffneAsset liefert einen ReadSeeker, der direkt aus dem
// eingebetteten Binary liest — ohne Kopie im Heap.
// Geeignet für große Assets (Musik), die gestreamt werden.
func OeffneAsset(name string) io.ReadSeeker {
	f, err := assetFS.Open(name)
	if err != nil {
		panic("Asset nicht gefunden: " + name)
	}
	return f.(io.ReadSeeker)
}
