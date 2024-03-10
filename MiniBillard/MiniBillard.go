package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./views"
	"./welt"
)

var (
	vAnstoß     hilf.Vec2
	updateLäuft bool
)

func starteUpdateProzess(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(12 * time.Millisecond)

	updater := func() {
		defer func() { println("Halte Spiel-Takt an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("Stoppe Spiel-Logik")
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

	var b, h uint16 = 1000, 600
	rand := b / 60
	var breiteSpiel uint16 = 800
	lS, bS := spiel.GibGröße()
	seitenverhältnis := lS / bS // breite/höhe

	maßstab := float64(breiteSpiel) / float64(lS)

	//öffne gfx-Fenster
	gfx.Fenster(b, h)
	gfx.Fenstertitel("unser Programmname")

	//erzeuge Zeichner
	billardTischZeichner :=
		views.NewFensterZeichner(rand, rand, 900+rand, uint16(900/seitenverhältnis)+rand, maßstab)
	hintergrundZeichner :=
		views.NewFensterZeichner(0, 0, b, h, 1.0)

		//definiere Zeichenfunktionen
	fülleHintergrund := func(fenster *views.FensterZeichner, r, g, b uint8) {
		fenster.FülleFläche(r, g, b)
	}

	zeichneBillardSpiel := func(fenster *views.FensterZeichner, spiel welt.MiniBillardSpiel) {
		// warte auf Bewegung der Kugeln
		for updateLäuft {
			time.Sleep(time.Millisecond)
		}
		fenster.ZeichneMiniBillardSpiel(spiel)
		// TODO die Spielelogik und die Skalierung muss hier raus
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			pS := spiel.GibStoßkugel().GibPos()
			fenster.ZeichneBreiteLinie(pS, pS.Plus(vAnstoß.Mal(15)), 5, 250, 175, 50)
		}
	}

	viewer := func() {
		defer func() { println("Halte Zeichen-Takt an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("Stoppe Zeichenprozess")
				takt.Stop()
				return
			case <-takt.C:
				gfx.UpdateAus()
				fülleHintergrund(hintergrundZeichner, 139, 69, 19)
				zeichneBillardSpiel(billardTischZeichner, spiel)
				gfx.UpdateAn()
			}
		}
	}

	// starte Zeichenprozess
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
		defer func() { println("Halte Controller-Takt an"); takt.Stop() }()
		for {
			select {
			case <-stop:
				println("Stoppe Maussteuerung")
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

func starteHintergrundPlayer(stop chan bool) {
	coolJazzTakt := time.NewTicker(2*time.Minute + 8*time.Second)
	ambienceTakt := time.NewTicker(time.Minute + 13*time.Second)
	klaenge.CoolJazzLoop2641SOUND()
	klaenge.BillardPubAmbienceSOUND()
	player := func() {
		defer func() {
			println("Halte Musik-Schleife an")
			coolJazzTakt.Stop()
			ambienceTakt.Stop()
		}()
		for {
			select {
			case <-stop:
				println("Stoppe Musik")
				return
			case <-coolJazzTakt.C:
				klaenge.CoolJazzLoop2641SOUND()
			case <-ambienceTakt.C:
				klaenge.BillardPubAmbienceSOUND()
			}
		}
	}

	// starte Prozess
	println("Starte Musik")
	go player()
}

func main() {
	spiel := welt.New3BallStandardSpiel()

	stopViewer, stopUpdater, stopController := make(chan bool), make(chan bool), make(chan bool)
	stopPlayer := make(chan bool)

	starteZeichenProzess(spiel, stopViewer)
	starteUpdateProzess(spiel, stopUpdater)
	starteMaussteuerung(spiel, stopController)
	starteHintergrundPlayer(stopPlayer)

	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'r': // reset
				spiel.Nochmal() // setze Kugeln wie vor dem letzten Anstoß
			case 'q': // quit
				stopPlayer <- true
				stopViewer <- true
				stopUpdater <- true
				stopController <- true
				time.Sleep(100 * time.Millisecond)
				if gfx.FensterOffen() {
					gfx.FensterAus()
				}
				// goodbye
				return
			}
		}
	}
}
