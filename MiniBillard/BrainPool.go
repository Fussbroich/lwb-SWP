package main

import (
	"gfx"

	"./hilf"
	"./klaenge"
	"./modelle"
	"./views_controls"
)

type BPApp interface {
	Run()
	Quit()
}

type bpapp struct {
	läuft        bool
	rastermaß    uint16
	breite, höhe uint16 // Größe des gesamten App-Fensters
	// Klänge
	musik, geräusche klaenge.Klang
	// Modelle
	billard modelle.MiniBillardSpiel
	quiz    modelle.Quiz
	// Views
	spieltisch, quizfenster, hilfebox           views_controls.Widget
	neuesSpielButton, hilfeButton, quitButton   views_controls.Widget
	hintergrund, punktezaehler, restzeit, bande views_controls.Widget
	renderer                                    views_controls.FensterZeichner
	// Controls
	mausSteuerung views_controls.EingabeProzess
}

// ####### baue die App zusammen ##################################################
func NewBPApp(g uint16) *bpapp {
	app := bpapp{rastermaß: g, breite: 32 * g, höhe: 22 * g}
	app.renderer = views_controls.NewFensterZeichner("BrainPool - Das MiniBillard für Schlaue.")

	app.musik = klaenge.CoolJazz2641SOUND()
	app.geräusche = klaenge.BillardPubAmbienceSOUND()

	// ######## Modelle und Views zusammenstellen #################################
	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	// Breite, Höhe des Spielfelds
	var bS uint16 = 3 * app.breite / 4
	var hS uint16 = bS / 2
	// Radius der Kugeln
	var ra uint16 = uint16(0.5 + float64(bS)*57.2/2540)

	// Modelle erzeugen
	app.billard = modelle.NewMini9BallSpiel(bS, hS, ra)
	app.quiz = modelle.NewQuizCSV("BeispielQuiz.csv")

	// Views und Zeichner erzeugen
	app.hintergrund = views_controls.NewFenster()
	app.renderer.SetzeFensterHintergrund(app.hintergrund)
	app.punktezaehler = views_controls.NewMBPunkteAnzeiger(app.billard)
	app.restzeit = views_controls.NewMBRestzeitAnzeiger(app.billard)
	app.bande = views_controls.NewFenster()
	app.spieltisch = views_controls.NewMBSpieltisch(app.billard)
	app.quizfenster = views_controls.NewQuizFenster(app.quiz)
	app.neuesSpielButton = views_controls.NewButton("neues Spiel",
		func() {
			app.billard.Reset()
			app.quizfenster.SetzeInAktiv()
			app.hilfebox.SetzeInAktiv()
			app.spieltisch.SetzeAktiv()
		})
	app.quitButton = views_controls.NewButton("Quit",
		func() {
			app.Quit()
		})
	app.hilfebox = views_controls.NewTextBox("Hilfe")
	app.hilfebox.SetzeInAktiv()
	app.hilfeButton = views_controls.NewButton("?",
		func() {
			app.hilfebox.AktivAnAus()
		})

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	app.renderer.SetzeWidgets(app.bande, app.spieltisch, app.quizfenster, app.hilfebox, app.punktezaehler, app.restzeit,
		app.neuesSpielButton, app.quitButton, app.hilfeButton)

	//setze Layout
	app.hintergrund.SetzeKoordinaten(0, 0, app.breite, app.höhe)
	var xs, ys, xe, ye uint16 = 4 * app.rastermaß, 6 * app.rastermaß, 28 * app.rastermaß, 18 * app.rastermaß
	var g3 uint16 = app.rastermaß + app.rastermaß/3
	app.punktezaehler.SetzeKoordinaten(xs-g3, 1*app.rastermaß, 18*app.rastermaß, 3*app.rastermaß)
	app.restzeit.SetzeKoordinaten(20*app.rastermaß+g3, app.rastermaß, xe+g3, 3*app.rastermaß)
	app.bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	app.bande.SetzeEckradius(g3)
	app.spieltisch.SetzeKoordinaten(xs, ys, xe, ye)
	app.quizfenster.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	app.quizfenster.SetzeEckradius(g3)
	app.neuesSpielButton.SetzeKoordinaten(app.breite/2-2*app.rastermaß, ye+g3+app.rastermaß/2, app.breite/2+2*app.rastermaß, ye+g3+g3)
	app.neuesSpielButton.SetzeEckradius(app.rastermaß / 3)
	app.hilfeButton.SetzeKoordinaten(2*app.rastermaß, ye+g3+app.rastermaß/2, 4*app.rastermaß, ye+g3+g3)
	app.hilfeButton.SetzeEckradius(app.rastermaß / 3)
	app.hilfebox.SetzeKoordinaten(app.breite/4, app.höhe/4, app.breite/2, app.höhe/2)
	app.quitButton.SetzeKoordinaten(app.breite-4*app.rastermaß, ye+g3+app.rastermaß/2, app.breite-2*app.rastermaß, ye+g3+g3)
	app.quitButton.SetzeEckradius(app.rastermaß / 3)

	//setzeFarben
	app.hintergrund.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	app.spieltisch.SetzeFarben(views_controls.Fbillardtuch(), views_controls.Fdiamanten())
	app.bande.SetzeFarben(views_controls.Ftext(), views_controls.Fanzeige())
	app.punktezaehler.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	app.punktezaehler.SetzeTransparenz(255)
	app.restzeit.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	app.quizfenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	app.neuesSpielButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	app.hilfeButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	app.hilfebox.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	app.quitButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())

	return &app
}

// ############### die ganze App wird mit der Maus gesteuert ########################
func (a *bpapp) mausSteuerFunktion(taste uint8, status int8, mausX, mausY uint16) {
	if a.quitButton != nil && a.quitButton.IstAktiv() &&
		a.quitButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.quitButton.MausklickBei(mausX, mausY)
		return
	} else if a.hilfeButton != nil && a.hilfeButton.IstAktiv() &&
		a.hilfeButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.hilfeButton.MausklickBei(mausX, mausY)
	} else if a.neuesSpielButton != nil && a.neuesSpielButton.IstAktiv() &&
		a.neuesSpielButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.neuesSpielButton.MausklickBei(mausX, mausY)
	} else if a.quizfenster != nil && a.quizfenster.IstAktiv() &&
		a.quizfenster.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.quizfenster.MausklickBei(mausX, mausY)
		if a.quiz.GibAktuelleFrage().RichtigBeantwortet() {
			a.billard.ReduziereStrafpunkte()
			if a.billard.GibStrafpunkte() <= a.billard.GibTreffer() {
				a.quizfenster.SetzeInAktiv()
				a.spieltisch.SetzeAktiv() // zurück zum Spielmodus
			}
		} else {
			a.quiz.NaechsteFrage()
		}
	} else if a.spieltisch != nil && a.spieltisch.IstAktiv() &&
		!a.hilfebox.IstAktiv() && a.billard.Laeuft() {
		if a.billard.GibStrafpunkte() > a.billard.GibTreffer() {
			a.quiz.NaechsteFrage()
			a.spieltisch.SetzeInAktiv()
			a.quizfenster.SetzeAktiv() // zum Quizmodus
		} else if a.billard.IstStillstand() {
			// zielen und stoßen
			switch taste {
			case 1: // stoßen
				if status == -1 {
					a.billard.Stosse()
				}
			case 4: // Stoßkraft erhöhen
				a.billard.SetzeStosskraft(a.billard.GibVStoss().Betrag() + 1)
			case 5: // Stoßkraft verringern
				a.billard.SetzeStosskraft(a.billard.GibVStoss().Betrag() - 1)
			default: // zielen
				xs, ys := a.spieltisch.GibStartkoordinaten()
				a.billard.SetzeStossRichtung((hilf.V2(float64(mausX), float64(mausY))).
					Minus(a.billard.GibSpielkugel().GibPos()).
					Minus(hilf.V2(float64(xs), float64(ys))))
			}
		}
	}
}

// ####### starte alles ##################################################
func (a *bpapp) Run() {
	if a.läuft {
		return
	}
	println("Willkommen bei BrainPool")
	a.billard.Starte()
	a.spieltisch.SetzeAktiv()
	a.quizfenster.SetzeInAktiv()
	a.hilfebox.SetzeInAktiv()
	a.renderer.Starte()
	a.mausSteuerung = views_controls.NewMausProzess(a.mausSteuerFunktion)
	a.mausSteuerung.Starte()
	a.geräusche.StarteLoop()
	a.musik.StarteLoop()
	a.läuft = true
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'p':
				a.billard.PauseAnAus()
			case 'c':
				a.renderer.DarkmodeAnAus()
			case 'h':
				a.hilfebox.AktivAnAus()
			case 'd':
				a.billard.ZeitlupeAnAus()
			case 'l':
				a.renderer.LayoutAnAus()
			case 'f':
				a.spieltisch.SetzeInAktiv()
				a.quizfenster.SetzeAktiv()
			case 'q':
				a.Quit()
				return
			}
		}
	}
}

// ####### stoppe alles ##################################################
func (a *bpapp) Quit() {
	if !a.läuft {
		return
	}
	a.geräusche.Stoppe()
	a.musik.Stoppe()
	a.renderer.UeberblendeText("Bye!", views_controls.Fanzeige(), views_controls.Ftext(), 30)
	a.mausSteuerung.Stoppe()
	a.billard.Stoppe()
	a.renderer.Stoppe()
	println("BrainPool wird beendet")
}

// ####### der Startpunkt ##################################################
func main() {
	// das Rastermaß bestimmt die Größe der App
	NewBPApp(35).Run()
}
