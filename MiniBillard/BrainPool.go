package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./views"
	"./welt"
)

type BrainPoolApp interface {
	Starte()
	Quit()
	ZeitlupeAnAus()
	PauseAnAus()
	QuizmodusAnAus()
}

type bpapp struct {
	billard       welt.MiniBillardSpiel
	quiz          welt.Quiz
	quizfenster   views.Fenster
	standardzeit  time.Duration
	renderer      views.FensterZeichner
	maussteuerung hilf.Prozess
	quizmodus     bool
	pause         bool
}

func NewBrainPoolApp(billard welt.MiniBillardSpiel, spieltisch views.Fenster,
	quiz welt.Quiz, quizfenster views.Fenster,
	hintergrund, bande, punktezähler, restzeit, neuesSpielButton views.Fenster) *bpapp {

	xs, ys := spieltisch.GibStartkoordinaten()
	app := bpapp{billard: billard, quiz: quiz, quizfenster: quizfenster,
		renderer: views.NewFensterZeichner(hintergrund, bande, spieltisch, punktezähler, restzeit, neuesSpielButton)}
	app.standardzeit = 4 * time.Minute

	app.maussteuerung = hilf.NewProzess("Maussteuerung",
		func() {
			// TODO: hier hängt es, wenn die Maus nicht bewegt wird.
			// Der Mauspuffer ist aber keine Lösung ...
			taste, _, mausX, mausY := gfx.MausLesen1()

			// Prüfe, wo die Maus gerade ist
			if app.quizmodus && quizfenster.ImFenster(mausX, mausY) {
				if taste == 1 {
					// Antwort prüfen
				}
			} else if neuesSpielButton.ImFenster(mausX, mausY) && taste == 1 {
				app.renderer.ÜberblendeAus()
				billard.Reset()
				billard.SetzeRestzeit(app.standardzeit)
				billard.Starte()
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
			}
		})
	return &app
}
func (app *bpapp) ZeitlupeAnAus() { app.billard.ZeitlupeAnAus() }

func (app *bpapp) PauseAnAus() {
	if !app.pause {
		app.billard.Stoppe()
		app.renderer.ÜberblendeText("Pause", views.F(225, 255, 255), views.F(249, 73, 68), 180)
	} else {
		app.renderer.ÜberblendeAus()
		app.billard.Starte()
	}
	app.pause = !app.pause
}

func (app *bpapp) QuizmodusAnAus() {
	if !app.quizmodus {
		app.billard.Stoppe()
		app.quiz.NächsteFrage()
		app.renderer.Überblende(app.quizfenster)
	} else {
		app.renderer.ÜberblendeAus()
		app.billard.Starte()
	}
	app.quizmodus = !app.quizmodus
}

func (app *bpapp) Starte() {
	// ######## starte Spiel-Prozesse ###########################################
	app.billard.SetzeRestzeit(app.standardzeit)
	app.billard.Starte()
	app.renderer.Starte()
	//renderer.ZeigeLayout()
	app.maussteuerung.StarteRate(50) // gewünschte Abtastrate je Sekunde
}

func (app *bpapp) Quit() {
	app.renderer.ÜberblendeText("Bye!", views.F(225, 255, 255), views.F(249, 73, 68), 30)
	app.maussteuerung.Stoppe()
	app.billard.Stoppe()
	app.renderer.Stoppe()
}

// ######## Hier kommt die gesamte App #########################################
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
	var bpapp BrainPoolApp = NewBrainPoolApp(billard, tisch, quiz, quizfenster,
		hintergrund, bande, punktezähler, restzeit, neuesSpielButton)

	// ######## Musik ###########################################################
	musik := klaenge.CoolJazz2641SOUND()
	//pulse := klaenge.MassivePulseSound()
	geräusche := klaenge.BillardPubAmbienceSOUND()

	// ######## Tastatur-Loop ###################################################
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
			case 'n': // nochmal
				//billard.StoßWiederholen() // setze Kugeln wie vor dem letzten Stoß
			case 'q': // quit
				geräusche.Stoppe()
				musik.Stoppe()
				bpapp.Quit()
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
