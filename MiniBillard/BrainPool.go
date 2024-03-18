package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./views"
	"./welt"
)

/*
type BrainPoolApp interface {
	MusikAn()
	PauseAnAus()
	DebugAnAus()
	StoßWiederholen()
	NeuesSpiel()
	Quit()
}

type mbapp struct {
	billard       welt.MiniBillardSpiel
	renderer      views.FensterZeichner
	maussteuerung hilf.Prozess
}

func NewMBApp(billard welt.MiniBillardSpiel, fenster ...views.Fenster) *mbapp {}

func (app *mbapp) MusikAn() {}
func (app *mbapp) PauseAnAus() {}
func (app *mbapp) DebugAnAus() {}
func (app *mbapp) StoßWiederholen() {}
func (app *mbapp) NeuesSpiel() {}
func (app *mbapp) Quit() {}
*/

// ######## Hier kommt die gesamte App #########################################
func main() {

	// ######## lege App-Größe fest ###########################################
	var g uint16 = 35 // Rastermaß
	xs, ys, xe, ye := 4*g, 6*g, 28*g, 18*g
	b, h := 32*g, 22*g
	g3 := g + g/3
	standardzeit := 4 * time.Minute

	println("Willkommen bei BrainPool")
	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h) //Fenstergröße
	gfx.Fenstertitel("BrainPool - Das MiniBillard für Schlaue.")

	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	var bS, hS uint16 = 24 * g, 12 * g        // Breite, Höhe des "Spielfelds"
	ra := uint16(0.5 + float64(bS)*57.2/2540) // Zeichenradius der Kugeln
	var billard welt.MiniBillardSpiel = welt.NewMiniPoolSpiel(bS, hS, ra)
	var quiz welt.Quiz = welt.NewAlphabetQuiz()
	var quizfenster views.Fenster
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

	var renderer views.FensterZeichner = views.NewFensterZeichner(
		hintergrund,
		punktezähler,
		restzeit,
		bande, tisch,
		neuesSpielButton)

	// ######## definiere Maussteuerung #########################################^
	inFenster := func(x, y uint16, f views.Fenster) bool {
		xs, ys := f.GibStartkoordinaten()
		b, h := f.GibGröße()
		return x > xs && x < xs+b && y > ys && y < ys+h
	}

	mausProzess := hilf.NewProzess("Maussteuerung",
		func() {
			// TODO: hier hängt das Spiel, wenn die Maus nicht bewegt wird
			// Der Mauspuffer ist aber zu langsam
			taste, _, mausX, mausY := gfx.MausLesen1()

			// Prüfe, wo die Maus gerade ist
			if quizfenster != nil && inFenster(mausX, mausY, quizfenster) {

			} else if inFenster(mausX, mausY, neuesSpielButton) && taste == 1 {
				billard.Reset()
				billard.SetzeRestzeit(standardzeit)
			} else if billard.Läuft() && billard.IstStillstand() && !billard.GibStoßkugel().IstEingelocht() {
				switch taste {
				case 1:
					billard.Stoße()
				case 4:
					billard.SetzeStoßStärke(billard.GibVStoß().Betrag() + 1)
				case 5:
					billard.SetzeStoßStärke(billard.GibVStoß().Betrag() - 1)
				default:
					billard.SetzeStoßRichtung((hilf.V2(float64(mausX), float64(mausY))).
						Minus(billard.GibStoßkugel().GibPos()).
						Minus(hilf.V2(float64(xs), float64(ys))))
				}
				// die Stoßstärke wird in "Kugelradien" gemessen
			}
		})

	// ######## Musik ###########################################################
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// ######## starte Spiel-Prozesse ###########################################
	billard.SetzeRestzeit(standardzeit)
	billard.Starte()
	renderer.Starte()
	//renderer.ZeigeLayout()
	mausProzess.StarteRate(50) // gewünschte Abtastrate je Sekunde

	// ######## Tastatur-Loop ###################################################
	var pause, fragemodus bool
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'f': // fragemodus
				if !fragemodus {
					billard.Stoppe()
					quizfenster = views.NewQuizFenster(quiz.GibNächsteFrage(), xs, ys, xe, ye,
						views.Weiß(), views.F(1, 88, 122))
					renderer.Überblende(quizfenster)
				} else {
					renderer.ÜberblendeAus()
					quizfenster = nil
					billard.Starte()
				}
				fragemodus = !fragemodus
			//
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
			case 'q': // quit
				renderer.ÜberblendeText("Bye!", views.F(225, 255, 255), views.F(249, 73, 68), 30)
				geräusche.Stoppe()
				musik.Stoppe()
				mausProzess.Stoppe()
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
