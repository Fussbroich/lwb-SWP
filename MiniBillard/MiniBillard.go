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

func view_komponente(spiel welt.MiniBillardSpiel) func() {
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
		views.NewBillardTischZeichner(xs, ys, xe, ye, float64(breiteSpielFenster)/float64(lS))
	stoßzählerZeichner :=
		views.NewSpielinfoZeichner(xs, ye+2, xe, h-rand)
	lernfragenZeichner :=
		views.NewHintergrundZeichner(xe+5, rand, b-rand, h-rand)
	hintergrundZeichner :=
		views.NewHintergrundZeichner(0, 0, b, h)

	return func() {
		gfx.UpdateAus()
		hintergrundZeichner.Zeichne(139, 69, 19)
		billardSpielZeichner.Zeichne(spiel)
		stoßzählerZeichner.Zeichne(spiel)
		lernfragenZeichner.Zeichne(200, 200, 200)
		//definiere Zeichenfunktionen
		// TODO die Skalierung muss hier raus
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			pS := spiel.GibStoßkugel().GibPos()
			billardSpielZeichner.ZeichneBreiteLinie(pS, pS.Plus(vAnstoß.Mal(15)), 5, 250, 175, 50)
		}
		gfx.UpdateAn()
	}
}

func main() {
	var spiel welt.MiniBillardSpiel = welt.New3BallStandardSpiel()

	// erzeuge Spiel-Prozesse
	updater := hilf.NewProzess("Spiel-Logik", func() { spiel.Update() })
	zeichner := hilf.NewProzess("View-Komponente", view_komponente(spiel))
	steuerung := hilf.NewProzess("Maussteuerung", maussteuerung(spiel))
	music := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	ambience := klaenge.BillardPubAmbienceSOUND()

	// starte Spiel-Prozesse
	updater.StarteLoop(12 * time.Millisecond)
	zeichner.StarteLoop(20 * time.Millisecond)
	steuerung.StarteLoop(5 * time.Millisecond)
	music.StarteLoop()
	ambience.StarteLoop()
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'r': // reset
				spiel.StoßWiederholen() // setze Kugeln wie vor dem letzten Anstoß
			case 'q': // quit
				ambience.StoppeLoop()
				music.StoppeLoop()
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
