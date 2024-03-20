package main

import (
	"gfx"

	"./controls"
	"./klaenge"
	"./views"
	"./welt"
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
	var billard welt.MiniBillardSpiel = welt.NewMini9BallSpiel(bS, hS, ra)
	var quiz welt.Quiz = welt.NewBeispielQuiz()

	// ######## erzeuge App-Fenster ###########################################
	// H: Hallenboden: F(218, 218, 218), Kneipenboden: views.F(104, 76, 65)
	hintergrund := views.NewFenster(0, 0, b, h,
		views.F(225, 232, 236), views.F(1, 88, 122), 0, 0)
	// Anzeige der Punkte
	punktezähler := views.NewMBPunkteAnzeiger(billard, xs-g3, 1*g, 18*g, 3*g,
		views.Weiß(), views.F(1, 88, 122), 255)
	// Anzeige restliche Zeit
	restzeit := views.NewMBRestzeitAnzeiger(billard, 20*g+g3, g, xe+g3, 3*g,
		views.Weiß(), views.F(1, 88, 122), 0)
	// Bande
	bande := views.NewFenster(xs-g3, ys-g3, xe+g3, ye+g3,
		views.F(1, 88, 122), views.Schwarz(), 0, g3)
	// Spielfeld
	tisch := views.NewMBSpieltischFenster(billard, xs, ys, xe, ye,
		views.F(92, 179, 193), views.Schwarz(), 0, 0)
	// neues-Spiel-Button
	neuesSpielButton := views.NewButton(b/2-2*g, ye+g3+g/2, b/2+2*g, ye+g3+g3, "neues Spiel",
		views.Weiß(), views.F(1, 88, 122), 100, g/3)
	//Quizfenster
	quizfenster := views.NewQuizFenster(quiz, xs-g3, ys-g3, xe+g3, ye+g3,
		views.Weiß(), views.F(1, 88, 122), g3)

	//erzeuge App-Control
	var bpapp controls.BPAppControl = controls.NewBPAppControl(
		billard, tisch, punktezähler, restzeit,
		quiz, quizfenster,
		hintergrund, bande, neuesSpielButton)

	// ######## Musik ###########################################################
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// ######## Tastatur-Loop #########################################
	bpapp.Starte()
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'f': // erzwinge fragemodus
				bpapp.QuizmodusAnAus()
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
