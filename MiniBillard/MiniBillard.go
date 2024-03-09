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
		for {
			select {
			case <-stop:
				println("Stoppe Zeichenprozess")
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
		if spiel.IstStillstand() && gfx.FensterOffen() {
			kS := spiel.GibStoßkugel()
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
		for {
			select {
			case <-stop:
				println("Stoppe Maussteuerung")
				takt.Stop()
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
	//spiel := welt.NewStandardSpielNewtonLinie(800, 400)
	spiel := welt.NewStandardSpiel(600, 300)
	//spiel := welt.NewNewtonRauteSpiel(600, 350)

	stopViewer, stopUpdater, stopController := make(chan bool), make(chan bool), make(chan bool)
	starteZeichenProzess(spiel, stopViewer)
	starteUpdateProzess(spiel, stopUpdater)
	starteMaussteuerung(spiel, stopController)

	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			println("Taste", taste, "gedrückt")
			switch taste {
			case 114: // Taste R
				spiel.Nochmal()
			case 113: // Taste Q
				stopController <- true
				stopUpdater <- true
				stopViewer <- true
				if gfx.FensterOffen() {
					gfx.FensterAus()
				}
				return
			}
		}
	}
}
