package gfx

// klang.go — WAV-Sound-Wiedergabe über Ebitengine-Audio.
// Geladene WAV-Dateien werden gecacht, damit wiederholtes Abspielen
// desselben Sounds ohne Disk-I/O auskommt.

import (
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const audioSampleRate = 44100

var (
	audioCtx *audio.Context

	// Cache: Pfad → dekodierte PCM-Daten ([]byte)
	soundCache = make(map[string][]byte)
)

func initAudio() {
	if audioCtx == nil {
		audioCtx = audio.NewContext(audioSampleRate)
	}
}

// setzeKlangparameter wird für Kompatibilität akzeptiert.
// Ebitengine verwaltet Audioparameter intern — die Werte werden ignoriert.
func setzeKlangparameter(_ uint32, _, _, _ uint8, _ float64) {
	initAudio()
}

// spieleSound spielt eine WAV-Datei im Hintergrund ab.
// Beim ersten Aufruf wird die Datei geladen und gecacht.
func spieleSound(pfad string) {
	initAudio()

	daten, ok := soundCache[pfad]
	if !ok {
		// Erstmaliges Laden und Dekodieren
		f, err := os.Open(pfad)
		if err != nil {
			return
		}
		stream, err := wav.DecodeWithSampleRate(audioSampleRate, f)
		if err != nil {
			f.Close()
			return
		}
		daten, err = io.ReadAll(stream)
		f.Close()
		if err != nil {
			return
		}
		soundCache[pfad] = daten
	}

	player := audioCtx.NewPlayerFromBytes(daten)
	player.Play()
}
