package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./modelle"
	"./views_controls"
)

var (
	// Modelle
	billard   modelle.MiniBillardSpiel
	quiz      modelle.Quiz
	quizmodus bool
	// Views
	hintergrund, bande      views_controls.Widget
	punktezaehler, restzeit views_controls.Widget
	spieltisch              views_controls.Widget
	neuesSpielButton        views_controls.Widget
	quizfenster             views_controls.Widget
	renderer                views_controls.FensterZeichner
	// Controls
	mausSteuerung views_controls.EingabeProzess
)

func mausSteuerFunktion(taste uint8, status int8, mausX, mausY uint16) {
	// im Quizmodus
	if quizmodus && quizfenster.ImFenster(mausX, mausY) {
		if taste == 1 && status == -1 {
			quizfenster.MausklickBei(mausX, mausY)
			if quiz.GibAktuelleFrage().RichtigBeantwortet() {
				billard.ReduziereStrafpunkte()
				if billard.GibStrafpunkte() <= billard.GibTreffer() {
					renderer.UeberblendeAus()
					billard.Starte()
					quizmodus = false // zurück zum Spielmodus
				}
			} else {
				quiz.NaechsteFrage()
			}
		}
		// neues Spiel starten geht immer
	} else if neuesSpielButton.ImFenster(mausX, mausY) && taste == 1 {
		renderer.UeberblendeAus()
		quizmodus = false
		billard.Reset()
		billard.Starte()
		// im Spielmodus
	} else if billard.Laeuft() {
		if billard.GibStrafpunkte() > billard.GibTreffer() {
			billard.Stoppe()
			quiz.NaechsteFrage()
			renderer.Ueberblende(quizfenster)
			quizmodus = true // zum Quizmodus
		} else if billard.IstStillstand() {
			// zielen und stoßen
			switch taste {
			case 1: // stoßen
				billard.Stosse()
			case 4: // Stoßkraft erhöhen
				billard.SetzeStosskraft(billard.GibVStoss().Betrag() + 1)
			case 5: // Stoßkraft verringern
				billard.SetzeStosskraft(billard.GibVStoss().Betrag() - 1)
			default: // zielen
				xs, ys := spieltisch.GibStartkoordinaten()
				billard.SetzeStossRichtung((hilf.V2(float64(mausX), float64(mausY))).
					Minus(billard.GibSpielkugel().GibPos()).
					Minus(hilf.V2(float64(xs), float64(ys))))
			}
		}
	}
}

func main() {
	println("Willkommen bei BrainPool")
	var rastermaß uint16 = 35 // Rastermaß bestimmt Größe der gesamten App
	var b, h uint16 = 32 * rastermaß, 22 * rastermaß

	// ######## Modelle und Views zusammenstellen #################################
	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	// Breite, Höhe des Spielfelds
	var bS, hS uint16 = 24 * rastermaß, 12 * rastermaß
	// Radius der Kugeln
	var ra uint16 = uint16(0.5 + float64(bS)*57.2/2540)

	// Modelle erzeugen
	billard = modelle.NewMini9BallSpiel(bS, hS, ra)
	quiz = modelle.NewQuizCSV("BeispielQuiz.csv")

	// Views erzeugen
	hintergrund = views_controls.NewFenster(0, 0, b, h, views_controls.F(225, 232, 236), views_controls.F(1, 88, 122), 0, 0)

	var xs, ys, xe, ye uint16 = 4 * rastermaß, 6 * rastermaß, 28 * rastermaß, 18 * rastermaß
	var g3 uint16 = rastermaß + rastermaß/3

	punktezaehler = views_controls.NewMBPunkteAnzeiger(billard, xs-g3, 1*rastermaß, 18*rastermaß, 3*rastermaß, views_controls.Weiß(), views_controls.F(1, 88, 122), 255)
	restzeit = views_controls.NewMBRestzeitAnzeiger(billard, 20*rastermaß+g3, rastermaß, xe+g3, 3*rastermaß, views_controls.Weiß(), views_controls.F(1, 88, 122), 0)
	bande = views_controls.NewFenster(xs-g3, ys-g3, xe+g3, ye+g3, views_controls.F(1, 88, 122), views_controls.Schwarz(), 0, g3)
	spieltisch = views_controls.NewMBSpieltisch(billard, xs, ys, xe, ye, views_controls.F(92, 179, 193), views_controls.F(180, 230, 255), 0, 0)
	neuesSpielButton = views_controls.NewButton(b/2-2*rastermaß, ye+g3+rastermaß/2, b/2+2*rastermaß, ye+g3+g3, "neues Spiel", views_controls.Weiß(), views_controls.F(1, 88, 122), 100, rastermaß/3)
	quizfenster = views_controls.NewQuizFenster(quiz, xs-g3, ys-g3, xe+g3, ye+g3, views_controls.Weiß(), views_controls.F(1, 88, 122), g3)

	// ######## Starte alles #########################################
	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	renderer = views_controls.NewFensterZeichner(hintergrund, bande, spieltisch, punktezaehler, restzeit, neuesSpielButton)
	mausSteuerung = views_controls.NewMausProzess(mausSteuerFunktion)

	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h) //Fenstergröße
	gfx.Fenstertitel("BrainPool - Das MiniBillard für Schlaue.")
	billard.Starte()
	renderer.Starte()
	mausSteuerung.Starte()

	// ######## Tastatur-Loop #########################################
	var musik, geräusche klaenge.Klang = klaenge.CoolJazz2641SOUND(), klaenge.BillardPubAmbienceSOUND()
	geräusche.StarteLoop()
	musik.StarteLoop()

	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'd': // Debug
				billard.ZeitlupeAnAus()
			case 'p': // Pause
				billard.PauseAnAus()
			case 'q': // quit
				// ######## Stoppe alles #########################################
				geräusche.Stoppe()
				musik.Stoppe()
				renderer.UeberblendeText("Bye!", views_controls.F(225, 255, 255), views_controls.F(249, 73, 68), 30)
				mausSteuerung.Stoppe()
				billard.Stoppe()
				renderer.Stoppe()
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
