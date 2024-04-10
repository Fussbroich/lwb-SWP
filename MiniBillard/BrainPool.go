package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./modelle"
	"./views_controls"
)

type BrainPoolApp interface {
	Run()
	Quit()
}

type mbapp struct {
	läuft            bool
	billard          modelle.MiniBillardSpiel
	spieltisch       views_controls.Widget
	spielzeit        time.Duration
	neuesSpielButton views_controls.Widget
	quiz             modelle.Quiz
	quizmodus        bool
	quizfenster      views_controls.Widget
	renderer         views_controls.FensterZeichner
	steuerProzess    hilf.Prozess
}

func NewBPApp(billard modelle.MiniBillardSpiel,
	spieltisch, punktezähler, restzeit views_controls.Widget,
	quiz modelle.Quiz, quizfenster,
	hintergrund, bande, neuesSpielButton views_controls.Widget) *mbapp {
	app := mbapp{
		billard: billard, neuesSpielButton: neuesSpielButton, spieltisch: spieltisch,
		quiz: quiz, quizfenster: quizfenster,
		renderer: views_controls.NewFensterZeichner(
			hintergrund, bande, spieltisch, punktezähler, restzeit, neuesSpielButton)}
	app.spielzeit = 4 * time.Minute
	return &app
}

func (app *mbapp) Run() {
	if app.läuft {
		return
	}
	app.billard.SetzeRestzeit(app.spielzeit)
	app.billard.Starte()
	app.renderer.Starte()
	//renderer.ZeigeLayout()
	app.steuerProzess = hilf.NewProzess(
		"Maussteuerung",
		app.mausSteuerung)
	app.steuerProzess.StarteRate(50) // gewünschte Abtastrate je Sekunde
	app.läuft = true
}

func (app *mbapp) Quit() {
	if !app.läuft {
		return
	}
	app.renderer.ÜberblendeText("Bye!", views_controls.F(225, 255, 255), views_controls.F(249, 73, 68), 30)
	app.steuerProzess.Stoppe()
	app.billard.Stoppe()
	app.renderer.Stoppe()
	app.läuft = false
	time.Sleep(100 * time.Millisecond)
}

func (app *mbapp) mausSteuerung() {
	taste, status, mausX, mausY := gfx.MausLesen1()

	// im Quizmodus
	if app.quizmodus && app.quizfenster.ImFenster(mausX, mausY) {
		if taste == 1 && status == -1 {
			app.quizfenster.MausklickBei(mausX, mausY)
			if app.quiz.GibAktuelleFrage().RichtigBeantwortet() {
				app.billard.ReduziereStrafpunkte()
				if app.billard.GibStrafpunkte() <= app.billard.GibTreffer() {
					app.quizmodusAus() // zurück zum Spielmodus
				}
			} else {
				app.quiz.NächsteFrage()
			}
		}
		// neues Spiel starten geht immer
	} else if app.neuesSpielButton.ImFenster(mausX, mausY) && taste == 1 {
		app.renderer.ÜberblendeAus()
		app.quizmodus = false
		app.billard.Reset()
		app.billard.SetzeRestzeit(app.spielzeit)
		app.billard.Starte()
		// im Spielmodus
	} else if app.billard.Läuft() {
		// TODO Maussteuerung und Regelprüfung trennen!
		if app.billard.GibStrafpunkte() > app.billard.GibTreffer() {
			app.quizmodusAn() // zum Quizmodus
		} else if app.billard.IstStillstand() {
			// zielen und stoßen
			switch taste {
			case 1: // stoßen
				app.billard.Stoße()
			case 4: // Stoßkraft erhöhen
				app.billard.SetzeStoßStärke(app.billard.GibVStoß().Betrag() + 1)
			case 5: // Stoßkraft verringern
				app.billard.SetzeStoßStärke(app.billard.GibVStoß().Betrag() - 1)
			default: // zielen
				xs, ys := app.spieltisch.GibStartkoordinaten()
				app.billard.SetzeStoßRichtung((hilf.V2(float64(mausX), float64(mausY))).
					Minus(app.billard.GibStoßkugel().GibPos()).
					Minus(hilf.V2(float64(xs), float64(ys))))
			}
		}
	}
}

func (app *mbapp) quizmodusAn() {
	app.billard.Stoppe()
	app.quiz.NächsteFrage()
	app.renderer.Überblende(app.quizfenster)
	app.quizmodus = true
}

func (app *mbapp) quizmodusAus() {
	app.renderer.ÜberblendeAus()
	app.billard.Starte()
	app.quizmodus = false
}

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
	// Hallenboden
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
	var bpapp BrainPoolApp = NewBPApp(
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
				billard.ZeitlupeAnAus()
			case 'm': // Musik an
				// einmal an bleibt an; stoppen geht mit gfx nicht.
				musik.StarteLoop()
				geräusche.StarteLoop()
			case 'p': // Pause
				billard.PauseAnAus()
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
