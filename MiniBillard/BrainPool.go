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
	var g uint16 = 35 // Rastermaß
	xs, ys, xe, ye := 4*g, 8*g, 28*g, 20*g
	g3 := g + g/3

	println("Öffne Gfx-Fenster")
	gfx.Fenster(32*g, 24*g) //Fenstergröße
	gfx.Fenstertitel("BrainPool - Das MiniBillard für Schlaue.")

	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	var bS, hS uint16 = 24 * g, 12 * g        // Breite, Höhe des "Spielfelds"
	ra := uint16(0.5 + float64(bS)*57.2/2540) // Zeichenradius der Kugeln
	var billard welt.MiniBillardSpiel = welt.NewMiniPoolSpiel(bS, hS, ra)

	// ######## erzeuge App-Fenster ###########################################
	var renderer views.FensterZeichner = views.NewFensterZeichner(
		// Hintergrund: Hallenboden: F(218, 218, 218), Kneipenboden: views.F(104, 76, 65)
		views.NewFenster(0, 0, 32*g, 24*g,
			views.F(218, 218, 218), views.Schwarz(), 0, 0),
		// Anzeige der Punkte
		views.NewMBPunkteAnzeiger(billard, xs-g3, 2*g, 18*g, 5*g,
			views.F(96, 108, 108), views.F(120, 135, 135), 255),
		// Anzeige Countdown
		views.NewFenster(20*g, 2*g, xe, 4*g,
			views.F(96, 108, 108), views.Schwarz(), 200, 0),
		// Bande
		views.NewFenster(xs-g3, ys-g3, xe+g3, ye+g3,
			views.F(96, 108, 108), views.Schwarz(), 0, g3),
		// Spielfeld
		views.NewMBSpieltischFenster(billard, xs, ys, xe, ye,
			views.F(91, 184, 207), views.Schwarz(), 0, 0))

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
				Minus(hilf.V2(float64(xs), float64(ys)))
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
	//renderer.ZeigeLayout()
	mausProzess.StarteRate(15) // gewünschte Abtastrate je Sekunde

	// ######## Tastatur-Loop ###################################################
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
					renderer.ÜberblendeText("Pause", views.F(225, 255, 255), views.F(249, 73, 68), 180)
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
				renderer.ÜberblendeText("Bye!", views.F(225, 255, 255), views.F(249, 73, 68), 30)
				geräusche.Stoppe()
				musik.Stoppe()
				billard.Stoppe()
				renderer.Stoppe()
				mausProzess.Stoppe()
				time.Sleep(100 * time.Millisecond)
				if gfx.FensterOffen() {
					gfx.FensterAus()
				}
				println("BrainPool wird beendet!")
				return
			}
		}
	}
}
