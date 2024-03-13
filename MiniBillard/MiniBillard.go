package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./views"
	"./welt"
)

func gibMaussteuerung(spiel welt.MiniBillardSpiel, xs, ys uint16) func() {
	return func() {
		if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
			// TODO: hier hängt der Prozess, solange die Maus nicht im Fenster ist
			taste, _, mausX, mausY := gfx.MausLesen1()
			vStoß := (hilf.V2(float64(mausX), float64(mausY))).
				Minus(spiel.GibStoßkugel().GibPos()).
				Minus(hilf.V2(float64(xs), float64(ys)))
			spiel.SetzeVStoß(vStoß.Mal(1 / spiel.GibStoßkugel().GibRadius()))
			if taste == 1 {
				spiel.Stoße()
			}
		}
	}

}

func gibAppZeichner(spiel welt.MiniBillardSpiel, b, h, rand uint16) func() {
	//erzeuge Zeichner
	bS, hS := spiel.GibGröße()
	xs, ys, xe, ye := rand, rand, rand+uint16(bS+0.5), uint16(hS+0.5)+rand
	billardSpielFenster :=
		views.NewMBSpielfeldZeichner(xs, ys, xe, ye)
	eingelochteAzeiger :=
		views.NewMBEingelochteZeichner(xs, ye+2, xe, ye+(h-ye-2-rand)/2, views.F(80, 80, 80))
	spielinfoFenster :=
		views.NewMBSpielinfoZeichner(xs, ye+(h-ye-2-rand)/2, xe, h-rand, views.F(80, 80, 80), views.F(200, 200, 200))
	lernfragenFenster :=
		views.NewHintergrundZeichner(xe+5, rand, b-rand, h-rand, views.F(200, 200, 200))
	hintergrund :=
		views.NewHintergrundZeichner(0, 0, b, h, views.F(139, 69, 19))

	return func() {
		gfx.UpdateAus()
		hintergrund.Zeichne()
		billardSpielFenster.Zeichne(spiel)
		spielinfoFenster.Zeichne(spiel)
		eingelochteAzeiger.Zeichne(spiel)
		lernfragenFenster.Zeichne()
		gfx.UpdateAn()
	}
}

func main() {
	//öffne gfx-Fenster
	var b, h, rand uint16 = 1280, 720, 30
	var spieltischBreite uint16 = 900
	println("Starte MiniBillard")
	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h)
	gfx.Fenstertitel("Das MiniBillard für Schlaumeier.")

	var spiel welt.MiniBillardSpiel = welt.NewMiniPoolSpiel(spieltischBreite)

	// erzeuge Spiel-Prozesse
	updater := hilf.NewProzess("Spiel-Logik", func() { spiel.Update() })
	zeichner := hilf.NewProzess("View-Komponente", gibAppZeichner(spiel, b, h, rand))
	steuerung := hilf.NewProzess("Maussteuerung", gibMaussteuerung(spiel, rand, rand))
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
			case 'n': // nochmal
				spiel.StoßWiederholen() // setze Kugeln wie vor dem letzten Stoß
			case 'r': // reset
				spiel.Reset() // setze Kugeln wie vor dem Anstoß
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
