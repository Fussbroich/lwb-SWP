package controls

import (
	"gfx"
	"time"

	"../hilf"
	"../views"
	"../welt"
)

type BPAppControl interface {
	Starte()
	Quit()
	ZeitlupeAnAus()
	PauseAnAus()
	NochmalVersuchen()
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

func NewBPAppControl(billard welt.MiniBillardSpiel, spieltisch, punktezähler, restzeit views.MBSpielView,
	quiz welt.Quiz, quizfenster views.QuizView, hintergrund, bande, neuesSpielButton views.Fenster) *bpapp {

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
				app.quizmodus = false
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

func (app *bpapp) NochmalVersuchen() { app.billard.StoßWiederholen() }

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
	time.Sleep(100 * time.Millisecond)
}
