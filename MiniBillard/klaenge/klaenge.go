package klaenge

import (
	"gfx" // Fenster muss offen sein
	"time"
)

type Klang interface {
	Play()
	StarteLoop()
	StoppeLoop()
}

type klang struct {
	titel string
	dauer time.Duration
	play  func()
	stop  chan bool
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
func MassivePulseSound() *klang {
	return &klang{
		titel: "Massive Pulse",
		dauer: 16 * time.Second,
		play: func() {
			if gfx.FensterOffen() {
				gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
				gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\massivePulseLoop.wav")
			}
		}}
}

func CoolJazz2641SOUND() *klang {
	return &klang{
		titel: "Cool Jazz 2641",
		dauer: 2*time.Minute + 8*time.Second,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(16000, 2, 2, 1, 0.3)
			gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\coolJazzLoop2641.wav")
		}}
}

func BillardPubAmbienceSOUND() *klang {
	return &klang{
		titel: "Billard Pub Ambience",
		dauer: time.Minute + 13*time.Second,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
			gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\billardPubAmbience.wav")
		}}
}

func CueHitsBallSound() *klang {
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\cueHitsBall.wav")
		}}
}

func BallHitsBallSound() *klang {
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\ballHitsBall.wav")
		}}
}

func BallInPocketSound() *klang {
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\ballIntoPocket.wav")
		}}
}

func BallHitsRailSound() *klang {
	return &klang{
		dauer: 300 * time.Millisecond,
		play: func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
			gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\ballHitsRail.wav")
		}}
}
