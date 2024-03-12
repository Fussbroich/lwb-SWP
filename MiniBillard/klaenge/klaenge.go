package klaenge

import (
	"errors"
	"gfx" // Fenster muss offen sein
	"os"
	"path/filepath"
	"time"
)

type Klang interface {
	Play()
	StarteLoop()
	StoppeLoop()
}

type klang struct {
	titel    string
	dauer    time.Duration
	filepath string
	play     func()
	stop     chan bool
}

func (s *klang) Play() {
	s.play()
}

func (s *klang) StarteLoop() {
	if s.stop != nil {
		return
	}
	s.stop = make(chan bool)
	takt := time.NewTicker(s.dauer)
	println("Starte Soundloop \"", s.titel, "\"")
	s.play()
	player := func() {
		defer func() { takt.Stop(); println("Stoppe Soundloop \"", s.titel, "\"") }()
		for {
			select {
			case <-s.stop:
				return
			case <-takt.C:
				s.play()
			}
		}
	}
	// starte Prozess
	go player()
}

func (s *klang) StoppeLoop() {
	s.stop <- true
	s.stop = nil
}

// gfx.SetzeKlangparameter(rate uint32, aufloesung, kanaele, signal uint8, p float64)
//
//	rate      ist die Abtastrate, z.B. 11025, 22050 oder 44100.
//	auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
//	kanaele   ist 1 für mono oder 2 für stereo.
//	signal    gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
//	pulsweite (für Rechtecksignale) gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.

func klangDateipfad(filename string) (fp string) {
	klaengeDir := "lwb-SWP/MiniBillard/klaenge"
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	fp = filepath.Join(filepath.Dir(wd), klaengeDir, filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return
}

func MassivePulseSound() *klang {
	fp := klangDateipfad("massivePulseLoop.wav")
	return &klang{
		titel: "Massive Pulse",
		dauer: 16 * time.Second,
		play: func() {
			if gfx.FensterOffen() {
				gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
				gfx.SpieleSound(fp)
			}
		}}
}

func CoolJazz2641SOUND() *klang {
	fp := klangDateipfad("coolJazzLoop2641.wav")
	return &klang{
		titel: "Cool Jazz 2641",
		dauer: 2*time.Minute + 8*time.Second,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(16000, 2, 2, 1, 0.3)
			gfx.SpieleSound(fp)
		}}
}

func BillardPubAmbienceSOUND() *klang {
	fp := klangDateipfad("billardPubAmbience.wav")
	return &klang{
		titel: "Billard Pub Ambience",
		dauer: time.Minute + 13*time.Second,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
			gfx.SpieleSound(fp)
		}}
}

func CueHitsBallSound() *klang {
	fp := klangDateipfad("cueHitsBall.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}

func BallHitsBallSound() *klang {
	fp := klangDateipfad("ballHitsBall.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}

func BallInPocketSound() *klang {
	fp := klangDateipfad("ballIntoPocket.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}

func BallHitsRailSound() *klang {
	fp := klangDateipfad("ballHitsRail.wav")
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound(fp)
		}}
}
