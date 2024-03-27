package controls

import (
	"gfx"
	"time"

	"../hilf"
	"../views"
	"../welt"
)

type MBAppControl interface {
	Starte()
	Quit()
	ZeitlupeAnAus()
	PauseAnAus()
}

type mbapp struct {
	läuft            bool
	billard          welt.MiniBillardSpiel
	spieltisch       views.Widget
	spielzeit        time.Duration
	pause            bool
	neuesSpielButton views.Widget
	quiz             welt.Quiz
	quizmodus        bool
	quizfenster      views.Widget
	renderer         views.FensterZeichner
	steuerProzess    hilf.Prozess
}

func NewMBAppControl(billard welt.MiniBillardSpiel,
	spieltisch, punktezähler, restzeit views.Widget,
	quiz welt.Quiz, quizfenster,
	hintergrund, bande, neuesSpielButton views.Widget) *mbapp {

	app := mbapp{
		billard: billard, neuesSpielButton: neuesSpielButton, spieltisch: spieltisch,
		quiz: quiz, quizfenster: quizfenster,
		renderer: views.NewFensterZeichner(
			hintergrund, bande, spieltisch, punktezähler, restzeit, neuesSpielButton)}
	app.spielzeit = 4 * time.Minute
	return &app
}

func (app *mbapp) Starte() {
	if app.läuft {
		return
	}
	app.billard.SetzeRestzeit(app.spielzeit)
	app.billard.Starte()
	app.renderer.Starte()
	//renderer.ZeigeLayout()
	app.steuerProzess = hilf.NewProzess("Maussteuerung", app.appSteuerung)
	app.steuerProzess.StarteRate(50) // gewünschte Abtastrate je Sekunde
	app.läuft = true
}

func (app *mbapp) Quit() {
	if !app.läuft {
		return
	}
	app.renderer.ÜberblendeText("Bye!", views.F(225, 255, 255), views.F(249, 73, 68), 30)
	app.steuerProzess.Stoppe()
	app.billard.Stoppe()
	app.renderer.Stoppe()
	app.läuft = false
	time.Sleep(100 * time.Millisecond)
}

func (app *mbapp) appSteuerung() {
	// TODO: hier hängt es, wenn die Maus nicht bewegt wird.
	// Ist der Mauspuffer eine Lösung ?
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
		if app.billard.GibStoßkugel().IstEingelocht() {
			app.billard.ErhöheStrafpunkte()
			app.billard.StoßWiederholen()
		} else if app.billard.GibStrafpunkte() > app.billard.GibTreffer() {
			app.quizmodusAn() // zum Quizmodus
		} else if app.billard.IstStillstand() {
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

func (app *mbapp) ZeitlupeAnAus() {
	if !app.läuft {
		return
	}
	app.billard.ZeitlupeAnAus()
}

func (app *mbapp) PauseAnAus() {
	if !app.läuft {
		return
	}
	if !app.pause {
		app.billard.Stoppe()
		app.renderer.ÜberblendeText("Pause", views.F(225, 255, 255), views.F(249, 73, 68), 180)
	} else {
		app.renderer.ÜberblendeAus()
		app.billard.Starte()
	}
	app.pause = !app.pause
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
