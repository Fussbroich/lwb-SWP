// Autoren:
//
//		Thomas Schrader
//		Bettina Chang
//
//		Zweck:
//			Das Spielprogramm BrainPool -
//			ein Softwareprojekt im Rahmen der Lehrerweiterbildung Berlin
//
//		Notwendige Software: Linux, Go ab 1.18
//			es läuft auch unter Windows, jedoch in verringerter Komplexität
//	     (Kugeln haben keine Nummern, Diamanten fehlen) um den gfx-Server
//			zu entlasten.
//		verwendete Pakete:
//			gfx, fmt, math, math/rand, strconv, strings, unicode/utf8, time,
//			runtime, os, errors, path/filepath, encoding/csv
//		Notwendige Hardware:
//			PC, Bildschirm, Tastatur, Maus mit Scrollrad
//
//		Datum: 01.05.2024
package main

import (
	"time"

	"./hilf"
	"./views_controls"
)

// Startet die Laufzeit-Elemente einer App.
//
//	Vor.: die App läuft nicht
//	Eff.: Die App wurde gestartet und ein gfx-Fenster geöffnet.
//	Hinweis: Es ist ratsam, die Tastensteuerung
//	lokal loopen zu lassen.
func RunApp(a views_controls.App) {
	// der Zeichen-Loop
	var zeichner views_controls.ZeichenRoutine
	// der Spiel-Loop
	var updater hilf.Routine
	// die Eingabe-Loops
	var mausSteuerung, tastenSteuerung views_controls.EingabeRoutine

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
	zeichner = views_controls.NewZeichenRoutine(a)
	zeichner.StarteRate(60) // go-Routine, öffnet das gfx-Fenster

	// ####### Simulation bringt eigenen Loop ######
	a.Reset() // startet go-Routine

	// ### der eigentliche Spiel-Loop der App läuft auch nebenher ###
	updater = hilf.NewRoutine("Umschalter", a.Update)
	updater.StarteRate(20) // go-Routine mit begrenzter Rate

	// ####### die Maussteuerung läuft ebenfalls nebenher ################
	mausSteuerung = views_controls.NewMausRoutine(a.MausEingabe)
	mausSteuerung.StarteRate(50) // go-Routine mit begrenzter Rate

	// ### Dafür darf der Tastatur-Loop hier existieren ########
	tastenSteuerung = views_controls.NewTastenRoutine(a.TastaturEingabe)
	tastenSteuerung.LoopeHier() // blockiert, bis quit() aufgerufen wird
}

// ####### der Startpunkt ##################################################
func main() {
	// Die gewünschte Fensterbreite in Pixeln wird übergeben.
	// Das Seitenverhältnis des Spiels ist B:H = 16:11
	RunApp(views_controls.NewBPApp(1024))
}
