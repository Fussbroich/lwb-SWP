package views_controls

import (
	"time"

	"../hilf"
)

// Startet die Laufzeit-Elemente einer App.
//
//	Vor.: die App läuft nicht
//	Eff.: Die App wurde gestartet und ein gfx-Fenster geöffnet.
//	Hinweis: Es ist ratsam, die Tastensteuerung
//	lokal loopen zu lassen.
func RunApp(a App) {
	// der Zeichen-Loop
	var zeichner ZeichenRoutine
	// der Spiel-Loop
	var updater hilf.Routine
	// die Eingabe-Loops
	var mausSteuerung, tastenSteuerung EingabeRoutine

	a.SetzeQuit(func() {
		go mausSteuerung.Stoppe()   // go-Routine, blockiert sonst
		go tastenSteuerung.Stoppe() // go-Routine, blockiert sonst
		updater.Stoppe()
		println("*********************************************")
		println("*** App wird beendet                      ***")
		println("*********************************************")
		time.Sleep(750 * time.Millisecond)
		zeichner.Stoppe()
	})

	println("*********************************************")
	println("*** Starte", a.GibTitel())
	println("*********************************************")

	//  ####### Zeichen-Loop nebenläufig starten ########
	zeichner = NewZeichenRoutine(a)
	zeichner.StarteRate(60) // go-Routine, öffnet das gfx-Fenster

	// ####### Simulation bringt eigenen Loop ######
	a.Reset() // startet go-Routine

	// ### der eigentliche Spiel-Loop der App läuft auch nebenher ###
	updater = hilf.NewRoutine("Umschalter", a.Update)
	updater.StarteRate(20) // go-Routine mit begrenzter Rate

	// ####### die Maussteuerung läuft ebenfalls nebenher ################
	mausSteuerung = NewMausRoutine(a.MausEingabe)
	mausSteuerung.StarteRate(50) // go-Routine mit begrenzter Rate

	// ### Dafür darf der Tastatur-Loop hier existieren ########
	tastenSteuerung = NewTastenRoutine(a.TastaturEingabe)
	tastenSteuerung.LoopeHier() // blockiert, bis quit() aufgerufen wird
}
