package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./welt"
)

var (
	vAnstoß     hilf.Vec2
	updateLäuft bool
)

func starteUpdateProzess(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(12 * time.Millisecond)

	updater := func() {
		defer func() { println("MiniBillard: Halte Spiel-Takt an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("MiniBillard: Stoppe Spiel-Logik")
				takt.Stop()
				return
			case <-takt.C:
				updateLäuft = true
				spiel.BewegeKugeln()
				updateLäuft = false
			}
		}
	}

	// starte Prozess
	println("Starte Spiel-Logik")
	go updater()
}

func starteZeichenProzess(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(20 * time.Millisecond)

	l, b := spiel.GibGröße()

	gfx.Fenster(uint16(l), uint16(b))
	gfx.Fenstertitel("Mini Billard")

	zeichne := func() {
		gfx.UpdateAus()
		gfx.Cls()
		l, b := spiel.GibGröße()
		// zeichne den Belag
		gfx.Stiftfarbe(60, 179, 113)
		gfx.Vollrechteck(0, 0, uint16(l), uint16(b))
		// zeichne die Taschen
		for _, t := range spiel.GibTaschen() {
			pos := t.GibPos()
			gfx.Stiftfarbe(0, 0, 0)
			gfx.Vollkreis(uint16(pos.X()), uint16(pos.Y()), uint16(t.GibRadius()))
		}
		// warte auf Bewegung der Kugeln
		for updateLäuft {
			time.Sleep(time.Millisecond)
		}
		// zeichne die Kugeln
		for _, k := range spiel.GibKugeln() {
			if k.IstEingelocht() {
				continue
			}
			pos := k.GibPos()
			ra := k.GibRadius()
			r, g, b := k.GibFarbe()
			gfx.Stiftfarbe(0, 0, 0)
			gfx.Vollkreis(uint16(pos.X()), uint16(pos.Y()), uint16(ra))
			gfx.Stiftfarbe(r, g, b)
			gfx.Vollkreis(uint16(pos.X()), uint16(pos.Y()), uint16(ra-1))
		}
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			pS := spiel.GibStoßkugel().GibPos()
			gfx.Stiftfarbe(250, 175, 50)
			hilf.ZeichneBreiteLinie(pS, pS.Plus(vAnstoß.Mal(15)), 5)
		}
		gfx.UpdateAn()
	}

	viewer := func() {
		defer func() { println("MiniBillard: Halte Zeichen-Takt an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("MiniBillard: Stoppe Zeichenprozess")
				takt.Stop()
				if gfx.FensterOffen() {
					gfx.FensterAus()
				}
				return
			case <-takt.C:
				zeichne()
			}
		}
	}

	// starte Prozess
	println("Starte Zeichenprozess")
	go viewer()
}

func starteMaussteuerung(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(5 * time.Millisecond)
	maustest := func() {
		if !gfx.FensterOffen() {
			return
		}
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			// TODO: hier hängt der Prozess, solange die Maus nicht im Fenster ist
			taste, _, mausX, mausY := gfx.MausLesen1()
			vAnstoß = (hilf.V2(float64(mausX), float64(mausY))).Minus(spiel.GibStoßkugel().GibPos()).Mal(1.0 / 15)
			vabs := vAnstoß.Betrag()
			if vabs > 10 {
				vAnstoß = vAnstoß.Mal(10 / vabs)
			}
			if taste == 1 {
				spiel.Anstoß(vAnstoß)
				klaenge.CueHitsBallSound()
			}
		}
	}
	mousecontroller := func() {
		defer func() { println("MiniBillard: Halte Controller-Takt an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("MiniBillard: Stoppe Maussteuerung")
				return
			case <-takt.C:
				maustest()
			}
		}
	}

	// starte Prozess
	println("Starte Maussteuerung")
	go mousecontroller()
}

func starteSound(stop chan bool) {
	takt := time.NewTicker(2*time.Minute + 8*time.Second)
	klaenge.CoolJazzLoop2641SOUND()
	music := func() {
		defer func() { println("MiniBillard: Halte Musik-Schleife an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("MiniBillard: Stoppe Musik")
				return
			case <-takt.C:
				klaenge.CoolJazzLoop2641SOUND()
			}
		}
	}

	// starte Prozess
	println("Starte Musik")
	go music()
}

func main() {
	spiel := welt.New3BallStandardSpiel(800)

	stopViewer, stopUpdater, stopController := make(chan bool), make(chan bool), make(chan bool)
	stopSound := make(chan bool)

	starteZeichenProzess(spiel, stopViewer)
	starteUpdateProzess(spiel, stopUpdater)
	starteMaussteuerung(spiel, stopController)
	starteSound(stopSound)

	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'r': // reset
				spiel.Nochmal() // setze Kugeln wie vor dem letzten Anstoß
			case 'q': // quit
				stopSound <- true
				stopViewer <- true
				stopUpdater <- true
				stopController <- true
				return
			}
		}
	}
}
