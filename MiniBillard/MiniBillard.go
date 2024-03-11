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
	vAnstoß hilf.Vec2
)

func starteUpdateProzess(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(12 * time.Millisecond)

	updater := func() {
		defer func() { takt.Stop() }()
		for {
			select {
			case <-stop:
				println("Stoppe Spiel-Logik")
				takt.Stop()
				return
			case <-takt.C:
				spiel.Update()
			}
		}
	}

	// starte Prozess
	println("Starte Spiel-Logik")
	go updater()
}

func starteHintergrundPlayer(stop chan bool) {
	coolJazzTakt := time.NewTicker(2*time.Minute + 8*time.Second)
	ambienceTakt := time.NewTicker(time.Minute + 13*time.Second)
	klaenge.CoolJazzLoop2641SOUND()
	klaenge.BillardPubAmbienceSOUND()
	player := func() {
		defer func() { coolJazzTakt.Stop(); ambienceTakt.Stop() }()
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
			if vabs > 12 {
				vAnstoß = vAnstoß.Mal(12 / vabs)
			}
			if taste == 1 {
				spiel.Anstoß(vAnstoß)
				klaenge.CueHitsBallSound()
			}
		}
	}
	mousecontroller := func() {
		defer func() { takt.Stop() }()
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

func starteZeichenProzess(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(20 * time.Millisecond)

	var b, h uint16 = 1200, 800
	rand := b / 60

	//öffne gfx-Fenster
	gfx.Fenster(b, h)
	gfx.Fenstertitel("unser Programmname")

	//erzeuge Zeichner
	lS, bS := spiel.GibGröße()
	seitenverhältnis := lS / bS // breite/höhe
	var breiteSpielFenster uint16 = 800 - 2*rand
	xs, ys, xe, ye := rand, rand, 900+rand, uint16(900/seitenverhältnis)+rand
	billardSpielZeichner :=
		views.NewFensterZeichner(xs, ys, xe, ye, float64(breiteSpielFenster)/float64(lS))

	stoßzählerZeichner :=
		views.NewFensterZeichner(xs, ye+2, xe, h-rand, 1.0)
	lernfragenZeichner :=
		views.NewFensterZeichner(xe+5, rand, b-rand, h-rand, 1.0)
	hintergrundZeichner :=
		views.NewFensterZeichner(0, 0, b, h, 1.0)

		//definiere Zeichenfunktionen
	fülleHintergrund := func(fenster *views.FensterZeichner, r, g, b uint8) {
		fenster.FülleFläche(r, g, b)
	}

	zeigeSpielinfo := func(fenster *views.FensterZeichner, spiel welt.MiniBillardSpiel) {
		fenster.ZeigeSpielinfo(spiel)
	}

	zeichneBillardSpiel := func(fenster *views.FensterZeichner, spiel welt.MiniBillardSpiel) {
		// warte auf Bewegung der Kugeln
		fenster.ZeichneMiniBillardSpiel(spiel)
		// TODO die Skalierung muss hier raus
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			pS := spiel.GibStoßkugel().GibPos()
			fenster.ZeichneBreiteLinie(pS, pS.Plus(vAnstoß.Mal(15)), 5, 250, 175, 50)
		}
	}

	viewer := func() {
		defer func() { takt.Stop() }()
		for {
			select {
			case <-stop:
				println("Stoppe Zeichenprozess")
				takt.Stop()
				return
			case <-takt.C:
				gfx.UpdateAus()
				fülleHintergrund(hintergrundZeichner, 139, 69, 19)
				zeichneBillardSpiel(billardSpielZeichner, spiel)
				zeigeSpielinfo(stoßzählerZeichner, spiel)
				fülleHintergrund(lernfragenZeichner, 200, 200, 200)
				gfx.UpdateAn()
			}
		}
	}

	// starte Zeichenprozess
	println("Starte Zeichenprozess")
	go viewer()
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
