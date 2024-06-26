package apps

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
var laeuft bool

func RunApp(a App) {
	if laeuft {
		return
	}
	laeuft = true

	// der Zeichen-Loop
	var zeichner hilf.Routine
	// der Spiel-Loop
	var updater hilf.Routine
	// die Eingabe-Loops
	var mausSteuerung, tastenSteuerung hilf.Routine

	a.SetzeQuit(func() {
		if !laeuft {
			return
		}
		laeuft = false
		println("*********************************************")
		println("*** App wird beendet                      ***")
		println("*********************************************")
		go mausSteuerung.Stoppe()   // go-Routine, blockiert sonst
		go tastenSteuerung.Stoppe() // go-Routine, blockiert sonst
		updater.Stoppe()
		time.Sleep(750 * time.Millisecond)
		zeichner.Stoppe()
	})

	println("*********************************************")
	println("*** Starte", a.GibTitel())
	println("*********************************************")

	//  ####### Zeichen-Loop nebenläufig starten ########
	zeichner = NewZeichenRoutine(a)
	zeichner.StarteMitRate(60) // go-Routine, öffnet das gfx-Fenster

	// ### der eigentliche Spiel-Loop der App läuft auch nebenher ###
	updater = hilf.NewRoutine("Updater", a.Update)
	updater.StarteMitRate(20) // go-Routine mit begrenzter Rate

	// ####### die Maussteuerung läuft ebenfalls nebenher ################
	mausSteuerung = NewMausRoutine(a.MausEreignis)
	mausSteuerung.StarteMitRate(50) // go-Routine mit begrenzter Rate

	// ### Dafür darf der Tastatur-Loop hier existieren ########
	tastenSteuerung = NewTastenRoutine(a.TastaturEreignis)
	tastenSteuerung.LoopeHier() // blockiert, bis quit() aufgerufen wird
}
