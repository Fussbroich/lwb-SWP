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

	var b, h uint16 = 1000, 600
	rand := b / 60

	//öffne gfx-Fenster
	gfx.Fenster(b, h)
	gfx.Fenstertitel("unser Programmname")

	//erzeuge ein View
	myView := views.NewGfxView(rand, 900+rand, rand, gfx.Grafikspalten())

	//definiere Zeichenfunktionen
	zeichneHintergrund := func(br, hö uint16, r, g, b uint8) {
		gfx.Stiftfarbe(r, g, b)
		gfx.Vollrechteck(0, 0, br, hö)
	}

	zeichne := func(v *views.GfxView, spiel welt.MiniBillardSpiel) {
		gfx.Cls()
		l, b := spiel.GibGröße()
		// zeichne den Belag
		v.ZeichneVollRechteck(hilf.V2(0, 0), l, b, 60, 179, 113)
		// zeichne die Taschen
		for _, t := range spiel.GibTaschen() {
			v.ZeichneVollKreis(t.GibPos(), t.GibRadius(), 0, 0, 0)
		}
		// warte auf Bewegung der Kugeln
		for updateLäuft {
			time.Sleep(time.Millisecond)
		}
		// zeichne die Kugeln
		// TODO: das ist Logik und muss hier weg
		for _, k := range spiel.GibKugeln() {
			if k.IstEingelocht() {
				continue
			}
			r, g, b := k.GibFarbe()
			v.ZeichneVollKreis(k.GibPos(), k.GibRadius(), 0, 0, 0)
			v.ZeichneVollKreis(k.GibPos(), k.GibRadius()-1, r, g, b)
		}
		// TODO: das ist Logik und muss hier weg
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			pS := spiel.GibStoßkugel().GibPos()
			v.ZeichneBreiteLinie(pS, pS.Plus(vAnstoß.Mal(15)), 5, 250, 175, 50)
		}
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
				gfx.UpdateAus()
				zeichneHintergrund(b, h, 139, 69, 19)
				zeichne(myView, spiel)
				//zeichneRahmen(uint16(l*maßstab)+2*rand, uint16(b*maßstab)+2*rand,
				//	rand, rand, uint16(l*maßstab)+rand, uint16(b*maßstab)+rand,
				//	139, 69, 19)
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

func starteHintergrundPlayer(stop chan bool) {
	coolJazzTakt := time.NewTicker(2*time.Minute + 8*time.Second)
	ambienceTakt := time.NewTicker(time.Minute + 13*time.Second)
	klaenge.CoolJazzLoop2641SOUND()
	klaenge.BillardPubAmbienceSOUND()
	player := func() {
		defer func() {
			println("MiniBillard: Halte Musik-Schleife an")
			coolJazzTakt.Stop()
			ambienceTakt.Stop()
		}()
		for {
			select {
			case <-stop:
				println("MiniBillard: Stoppe Musik")
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
	spiel := welt.New3BallStandardSpiel(800)

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
				return
			}
		}
	}
}
