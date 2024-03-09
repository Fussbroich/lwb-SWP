package main

import (
	"gfx"
	"time"

	"./hilf"
	"./welt"
)

var (
	vAnstoß     hilf.Vec2
	updateLäuft bool
)

func starteUpdateProzess(spiel welt.MiniBillardSpiel, stop chan bool) {
	takt := time.NewTicker(8 * time.Millisecond)

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
		gfx.Stiftfarbe(225, 255, 255)
		l, b := spiel.GibGröße()
		gfx.Vollrechteck(0, 0, uint16(l), uint16(b))
		// zeichne den Bahnbelag
		gfx.Stiftfarbe(60, 179, 113)
		for _, d := range spiel.GibBahnDreiecke() {
			hilf.ZeichneVollDreieck(d[0], d[1], d[2])
		}
		// zeichne die Taschen
		for _, t := range spiel.GibTaschen() {
			pos := t.GibPos()
			gfx.Stiftfarbe(0, 0, 0)
			gfx.Vollkreis(uint16(pos.X()), uint16(pos.Y()), uint16(t.GibRadius()))
		}
		gfx.Stiftfarbe(210, 105, 30)
		// zeichne die Banden
		for _, b := range spiel.GibBanden() {
			hilf.ZeichneBreiteLinieRechts(b.GibVon(), b.GibNach(), 15)
		}
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
		if spiel.IstStillstand() {
			kS := spiel.GibStoßkugel()
			// TODO: hier hängt der Prozess, solange die Maus nicht im Fenster ist
			taste, _, mausX, mausY := gfx.MausLesen1()
			vAnstoß = (hilf.V2(float64(mausX), float64(mausY))).Minus(kS.GibPos()).Mal(1.0 / 15)
			vabs := vAnstoß.Betrag()
			if vabs > 12 {
				vAnstoß = vAnstoß.Mal(12 / vabs)
			}
			if taste == 1 {
				spiel.Anstoß(vAnstoß)
			}
		}
	}
	mousecontroller := func() {
		defer func() { println("MiniBillard: Halte Controller-Takt an"); takt.Stop() }()
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

func main() {
	//spiel := welt.NewHexaBahnSpiel(600)
	//spiel := welt.NewLBahnSpiel(600)
	spiel := welt.New3BallLinieStandardSpiel(800)
	//spiel := welt.New3BallStandardSpiel(800)

	//spiel := welt.NewNewtonRauteSpiel(600, 350)

	stopViewer, stopUpdater, stopController := make(chan bool), make(chan bool), make(chan bool)
	starteZeichenProzess(spiel, stopViewer)
	starteUpdateProzess(spiel, stopUpdater)
	starteMaussteuerung(spiel, stopController)

	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'r': // reset
				spiel.Nochmal() // setze Kugeln wie vor dem letzten Anstoß
			case 'q': // quit
				stopViewer <- true
				stopUpdater <- true
				stopController <- true
				return
			}
		}
	}
}
