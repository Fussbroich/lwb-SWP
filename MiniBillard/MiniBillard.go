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

	println("Starte MiniBillard")

	// ######## lege App-Größe fest ###########################################
	var b, h, rand uint16 = 960, 720, 30
	var spiel welt.MiniBillardSpiel = welt.NewMiniPoolSpiel(900) // Spieltisch mit Breite
	bS, hS := spiel.GibGröße()
	xs, ys, xe, ye := rand, rand, rand+uint16(bS+0.5), uint16(hS+0.5)+rand

	// ######## erzeuge App-Fenster ###########################################
	var hintergrund views.Fenster = views.NewFenster(0, 0, b, h, views.F(139, 69, 19))
	var billardSpielFenster views.Fenster = views.NewMBSpielfeldFenster(spiel,
		xs, ys, xe, ye)
	var eingelochteAzeiger views.Fenster = views.NewMBEingelochteFenster(spiel,
		xs, ye+2, xe, ye+(h-ye-2-rand)/2, views.F(80, 80, 80))
	var spielinfoFenster views.Fenster = views.NewMBSpielinfoFenster(spiel,
		xs, ye+(h-ye-2-rand)/2, xe, h-rand, views.F(80, 80, 80), views.F(200, 200, 200))

	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h)
	gfx.Fenstertitel("BrainPool — Das MiniBillard für Schlaumeier.")

	// ######## erzeuge Spiel-Prozesse #########################################
	zeichenProzess := hilf.NewProzess("View-Komponente",
		func() {
			gfx.UpdateAus()
			gfx.Cls()
			hintergrund.Zeichne()
			billardSpielFenster.Zeichne()
			spielinfoFenster.Zeichne()
			eingelochteAzeiger.Zeichne()
			gfx.UpdateAn()
		})

	mausProzess := hilf.NewProzess("Maussteuerung",
		func() {
			// Die Maussteuerung ist nur aktiv, wenn alle Kugeln stehen.
			if spiel.IstStillstand() && !spiel.GibStoßkugel().IstEingelocht() {
				taste, _, mausX, mausY := gfx.MausLesen1()
				vStoß := (hilf.V2(float64(mausX), float64(mausY))).
					Minus(spiel.GibStoßkugel().GibPos()).
					Minus(hilf.V2(float64(rand), float64(rand)))
				// die Stoßstärke wird in "Kugelradien" gemessen
				spiel.SetzeVStoß(vStoß.Mal(1 / spiel.GibStoßkugel().GibRadius()))
				if taste == 1 {
					spiel.Stoße()
				}
			}
		})
	// ######## Musik ###########################################################
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// ######## starte Spiel-Prozesse ###########################################
	spiel.Starte()
	zeichenProzess.StarteRate(120)
	mausProzess.StarteRate(120)
	musik.StarteLoop()
	geräusche.StarteLoop()

	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'd': // Debug
				spiel.DebugAnAus()
			case 'p': // Pause
				spiel.PauseAnAus()
			case 'n': // nochmal
				spiel.StoßWiederholen() // setze Kugeln wie vor dem letzten Stoß
			case 'r': // reset
				spiel.Reset() // setze Kugeln wie vor dem Anstoß
			case 'q': // quit
				geräusche.Stoppe()
				musik.Stoppe()
				spiel.Stoppe()
				zeichenProzess.Stoppe()
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
