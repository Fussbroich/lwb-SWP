package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./views"
	"./welt"
)

// ######## Hier kommt die gesamte App #########################################
func main() {

	println("Starte BrainPool")

	// ######## lege App-Größe fest ###########################################
	var b, h, rand uint16 = 960, 720, 30
	var billard welt.MiniBillardSpiel = welt.NewMiniPoolSpiel(900) // Spieltisch mit Breite
	bS, hS := billard.GibGröße()
	xs, ys, xe, ye := rand, rand, rand+uint16(bS+0.5), uint16(hS+0.5)+rand

	// ######## erzeuge App-Fenster ###########################################
	var renderer views.FensterZeichner = views.NewFensterZeichner(
		// Hintergrund
		views.NewFenster(0, 0, b, h, views.F(104, 76, 65)),
		// Spieltisch
		views.NewMBSpielfeldFenster(billard, xs, ys, xe, ye),
		// Anzeige der eingelochten
		views.NewMBEingelochteFenster(billard, xs, ye+rand, xe, ye-rand+(h-ye)/2, views.F(96, 108, 108), views.F(120, 135, 135)),
		// Infos zum Spielverlauf
		views.NewMBSpielinfoFenster(billard, xs, ye-rand+(h-ye-5)/2+5, xe, h-rand, views.F(96, 108, 108), views.F(120, 135, 135)))

	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h)
	gfx.Fenstertitel("BrainPool - Das MiniBillard für Schlaue.")

	// ######## definiere Maussteuerung #########################################
	mausProzess := hilf.NewProzess("Maussteuerung",
		func() {
			if !billard.Läuft() {
				return
			}
			// Die Maussteuerung ist nur aktiv, wenn alle Kugeln stehen.
			if !billard.IstStillstand() {
				return
			}
			if billard.GibStoßkugel().IstEingelocht() {
				return
			}
			taste, _, mausX, mausY := gfx.MausLesen1()
			vStoß := (hilf.V2(float64(mausX), float64(mausY))).
				Minus(billard.GibStoßkugel().GibPos()).
				Minus(hilf.V2(float64(rand), float64(rand)))
			// die Stoßstärke wird in "Kugelradien" gemessen
			billard.SetzeVStoß(vStoß.Mal(1 / billard.GibStoßkugel().GibRadius()))

			// der Stoß wird ausgeführt
			if taste == 1 {
				billard.Stoße()
			}
		})

	// ######## Musik ###########################################################
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// ######## starte Spiel-Prozesse ###########################################
	billard.Starte()
	renderer.Starte()
	mausProzess.StarteRate(15) // gewünschte Abtastrate je Sekunde

	// ######## Tastatur-Loop######## ###########################################
	var pause bool
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'd': // Debug
				billard.ZeitlupeAnAus()
			case 'm': // Musik an
				// einmal an bleibt an; stoppen geht mit gfx nicht.
				musik.StarteLoop()
				geräusche.StarteLoop()
			case 'p': // Pause
				if !pause {
					billard.Stoppe()
					renderer.ÜberblendeText("Pause", views.F(249, 73, 68))
				} else {
					renderer.ÜberblendeAus()
					billard.Starte()
				}
				pause = !pause
			case 'n': // nochmal
				billard.StoßWiederholen() // setze Kugeln wie vor dem letzten Stoß
			case 'r': // reset
				billard.Reset() // setze Kugeln wie vor dem Anstoß
			case 'q': // quit
				geräusche.Stoppe()
				musik.Stoppe()
				billard.Stoppe()
				renderer.Stoppe()
				mausProzess.Stoppe()
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
