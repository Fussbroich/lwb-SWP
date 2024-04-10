package main

import (
	"gfx"

	"./klaenge"
	"./modelle"
	"./views_controls"
)

// ######## Hier wird die App zusammengestellt und gestartet ############################
func main() {

	// ######## lege App-Größe fest ###########################################
	var g uint16 = 35 // Rastermaß
	xs, ys, xe, ye := 4*g, 6*g, 28*g, 18*g
	b, h := 32*g, 22*g
	g3 := g + g/3

	println("Willkommen bei BrainPool")
	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h) //Fenstergröße
	gfx.Fenstertitel("BrainPool - Das MiniBillard für Schlaue.")

	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	var bS, hS uint16 = 24 * g, 12 * g        // Breite, Höhe des "Spielfelds"
	ra := uint16(0.5 + float64(bS)*57.2/2540) // Zeichenradius der Kugeln
	var billard modelle.MiniBillardSpiel = modelle.NewMini9BallSpiel(bS, hS, ra)
	var quiz modelle.Quiz = modelle.NewQuizCSV("BeispielQuiz.csv")

	// ######## erzeuge App-Fenster ###########################################
	// H: Hallenboden: F(218, 218, 218), Kneipenboden: views_controls.F(104, 76, 65)
	hintergrund := views_controls.NewFenster(0, 0, b, h,
		views_controls.F(225, 232, 236), views_controls.F(1, 88, 122), 0, 0)
	// Anzeige der Punkte
	punktezähler := views_controls.NewMBPunkteAnzeiger(billard, xs-g3, 1*g, 18*g, 3*g,
		views_controls.Weiß(), views_controls.F(1, 88, 122), 255)
	// Anzeige restliche Zeit
	restzeit := views_controls.NewMBRestzeitAnzeiger(billard, 20*g+g3, g, xe+g3, 3*g,
		views_controls.Weiß(), views_controls.F(1, 88, 122), 0)
	// Bande
	bande := views_controls.NewFenster(xs-g3, ys-g3, xe+g3, ye+g3,
		views_controls.F(1, 88, 122), views_controls.Schwarz(), 0, g3)
	// Spielfeld
	tisch := views_controls.NewMBSpieltisch(billard, xs, ys, xe, ye,
		views_controls.F(92, 179, 193), views_controls.F(180, 230, 255), 0, 0)
	// neues-Spiel-Button
	neuesSpielButton := views_controls.NewButton(b/2-2*g, ye+g3+g/2, b/2+2*g, ye+g3+g3,
		"neues Spiel",
		views_controls.Weiß(), views_controls.F(1, 88, 122), 100, g/3)
	//Quizfenster
	quizfenster := views_controls.NewQuizFenster(quiz, xs-g3, ys-g3, xe+g3, ye+g3,
		views_controls.Weiß(), views_controls.F(1, 88, 122), g3)

	//erzeuge App-Control
	var bpapp views_controls.MBAppControl = views_controls.NewMBAppControl(
		billard, tisch, punktezähler, restzeit,
		quiz, quizfenster,
		hintergrund, bande, neuesSpielButton)

	// ######## Musik ###########################################################
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// ######## Tastatur-Loop #########################################
	bpapp.Run()
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'd': // Debug
				bpapp.ZeitlupeAnAus()
			case 'm': // Musik an
				// einmal an bleibt an; stoppen geht mit gfx nicht.
				musik.StarteLoop()
				geräusche.StarteLoop()
			case 'p': // Pause
				bpapp.PauseAnAus()
			case 'q': // quit
				geräusche.Stoppe()
				musik.Stoppe()
				bpapp.Quit()
				if gfx.FensterOffen() {
					gfx.FensterAus()
				}
				println("BrainPool wird beendet!")
				return
			}
		}
	}
}
