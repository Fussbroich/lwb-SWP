package klaenge

import (
	"errors"
	"gfx"
	"os"
	"path/filepath"
	"time"
)

func assetDateipfad(filename string) (fp string) {
	klaengeDir := "MiniBillard/assets/soundfiles"
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	wdir := filepath.Dir(wd)
	fp = filepath.Join(wdir, klaengeDir, filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return
}

// gfx.SetzeKlangparameter(rate uint32, aufloesung, kanaele, signal uint8, p float64)
//
//	rate      ist die Abtastrate, z.B. 11025, 22050 oder 44100.
//	auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
//	kanaele   ist 1 für mono oder 2 für stereo.
//	signal    gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
//	pulsweite (für Rechtecksignale) gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.

func MassivePulseSound() *klang {
	fp := assetDateipfad("massivePulseLoop.wav")
	return &klang{
		titel: "Massive Pulse",
		dauer: 16 * time.Second,
		autor: "unknown",
		play: func() {
			if gfx.FensterOffen() {
				gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
				gfx.SpieleSound(fp)
			}
		}}
}

func CoolJazz2641SOUND() *klang {
	fp := assetDateipfad("coolJazzLoop2641.wav")
	return &klang{
		titel: "Cool Jazz 2641",
		dauer: 2*time.Minute + 8*time.Second,
		autor: "Julius H. (pixabay)",
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(16000, 2, 2, 1, 0.3)
			gfx.SpieleSound(fp)
		}}
}

func BillardPubAmbienceSOUND() *klang {
	fp := assetDateipfad("billardPubAmbience.wav")
	return &klang{
		titel: "Billard Pub Ambience",
		dauer: time.Minute + 13*time.Second,
		autor: "",
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
			gfx.SpieleSound(fp)
		}}
}

func CueHitsBallSound() *klang {
	fp := assetDateipfad("cueHitsBall.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}

func BallHitsBallSound() *klang {
	fp := assetDateipfad("ballHitsBall.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}

func BallInPocketSound() *klang {
	fp := assetDateipfad("ballIntoPocket.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}

func BallHitsRailSound() *klang {
	fp := assetDateipfad("ballHitsRail.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}
