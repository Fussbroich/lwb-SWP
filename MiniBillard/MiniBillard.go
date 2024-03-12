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

func maussteuerung(spiel welt.MiniBillardSpiel) func() {
	return func() {
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			// TODO: hier hängt der Prozess, solange die Maus nicht im Fenster ist
			taste, _, mausX, mausY := gfx.MausLesen1()
			// TODO: das muss in eine Extra-Klasse
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

}

func view_komponente(spiel welt.MiniBillardSpiel, b, h, rand uint16) func() {
	//erzeuge Zeichner
	bS, hS := spiel.GibGröße()
	xs, ys, xe, ye := rand, rand, rand+uint16(bS+0.5), uint16(hS+0.5)+rand
	billardSpielFenster :=
		views.NewBillardTischZeichner(xs, ys, xe, ye)
	stoßzählerFenster :=
		views.NewSpielinfoZeichner(xs, ye+2, xe, h-rand)
	lernfragenFenster :=
		views.NewHintergrundZeichner(xe+5, rand, b-rand, h-rand)
	hintergrund :=
		views.NewHintergrundZeichner(0, 0, b, h)

	return func() {
		gfx.UpdateAus()
		hintergrund.Zeichne(139, 69, 19)
		billardSpielFenster.Zeichne(spiel)
		stoßzählerFenster.Zeichne(spiel)
		lernfragenFenster.Zeichne(200, 200, 200)
		//definiere Zeichenfunktionen
		// TODO die Skalierung muss hier raus
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			pS := spiel.GibStoßkugel().GibPos()
			billardSpielFenster.ZeichneBreiteLinie(pS, pS.Plus(vAnstoß.Mal(15)), 5, 250, 175, 50)
		}
		gfx.UpdateAn()
	}
}

func main() {
	//öffne gfx-Fenster
	var b, h, rand uint16 = 1280, 720, 10

	gfx.Fenster(b, h)
	gfx.Fenstertitel("unser Programmname")

	var spiel welt.MiniBillardSpiel = welt.New3BallStandardSpiel(900)

	// erzeuge Spiel-Prozesse
	updater := hilf.NewProzess("Spiel-Logik", func() { spiel.Update() })
	zeichner := hilf.NewProzess("View-Komponente", view_komponente(spiel, b, h, rand))
	steuerung := hilf.NewProzess("Maussteuerung", maussteuerung(spiel))
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// starte Spiel-Prozesse
	updater.StarteLoop(12 * time.Millisecond)
	zeichner.StarteLoop(15 * time.Millisecond)
	steuerung.StarteLoop(5 * time.Millisecond)
	musik.StarteLoop()
	geräusche.StarteLoop()
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'r': // reset
				spiel.StoßWiederholen() // setze Kugeln wie vor dem letzten Anstoß
			case 'q': // quit
				geräusche.StoppeLoop()
				musik.StoppeLoop()
				updater.StoppeLoop()
				zeichner.StoppeLoop()
				steuerung.StoppeLoop()
				time.Sleep(100 * time.Millisecond)
				if gfx.FensterOffen() {
					gfx.FensterAus()
				}
				println("MiniBillard wird beendet!")
				return
			}
		}
	}
}
