package main

import (
	"gfx"
	"time"

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
	laeuft     bool
	rastermass uint16
	breite     uint16 // Größe des gesamten App-Fensters
	hoehe      uint16 // Größe des gesamten App-Fensters
	// Klänge
	musik      klaenge.Klang
	geraeusche klaenge.Klang
	// Modelle
	billard modelle.MiniBillardSpiel
	quiz    modelle.Quiz
	// Views
	spieltisch       views_controls.Widget
	quizfenster      views_controls.Widget
	hilfebox         views_controls.Widget
	neuesSpielButton views_controls.Widget
	hilfeButton      views_controls.Widget
	quitButton       views_controls.Widget
	hintergrund      views_controls.Widget
	punktezaehler    views_controls.Widget
	restzeit         views_controls.Widget
	bande            views_controls.Widget
	renderer         views_controls.FensterZeichner
	// Controls
	mausSteuerung views_controls.EingabeRoutine
}

// ####### baue die App zusammen ##################################################
func NewBPApp(g uint16) *bpapp {
	app := bpapp{rastermass: g, breite: 32 * g, hoehe: 22 * g}
	app.renderer = views_controls.NewFensterZeichner("BrainPool - Das MiniBillard für Schlaue.")

	app.musik = klaenge.CoolJazz2641SOUND()
	app.geraeusche = klaenge.BillardPubAmbienceSOUND()

	// ######## Modelle und Views zusammenstellen #################################
	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	// Breite, Höhe des Spielfelds
	var bS uint16 = 3 * app.breite / 4
	var hS uint16 = bS / 2
	// Radius der Kugeln
	var ra uint16 = uint16(0.5 + float64(bS)*57.2/2540)

	// Modelle erzeugen
	app.billard = modelle.NewMini9BallSpiel(bS, hS, ra)
	app.quiz = modelle.NewQuizInformatiksysteme()

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
			app.quizfenster.Ausblenden()
			app.hilfebox.Ausblenden()
			app.spieltisch.Einblenden()
		})
	app.quitButton = views_controls.NewButton("Quit",
		func() {
			app.Quit()
		})
	app.hilfebox = views_controls.NewTextBox("Hilfe")
	app.hilfebox.Ausblenden()
	app.hilfeButton = views_controls.NewButton("?",
		func() {
			app.hilfebox.DarstellenAnAus()
		})

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	app.renderer.SetzeWidgets(app.bande, app.spieltisch, app.quizfenster, app.hilfebox, app.punktezaehler, app.restzeit,
		app.neuesSpielButton, app.quitButton, app.hilfeButton)

	//setze Layout
	app.hintergrund.SetzeKoordinaten(0, 0, app.breite, app.hoehe)
	var xs, ys, xe, ye uint16 = 4 * app.rastermass, 6 * app.rastermass, 28 * app.rastermass, 18 * app.rastermass
	var g3 uint16 = app.rastermass + app.rastermass/3
	app.punktezaehler.SetzeKoordinaten(xs-g3, 1*app.rastermass, 18*app.rastermass, 3*app.rastermass)
	app.restzeit.SetzeKoordinaten(20*app.rastermass+g3, app.rastermass, xe+g3, 3*app.rastermass)
	app.bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	app.bande.SetzeEckradius(g3)
	app.spieltisch.SetzeKoordinaten(xs, ys, xe, ye)
	app.quizfenster.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	app.quizfenster.SetzeEckradius(g3)
	app.neuesSpielButton.SetzeKoordinaten(app.breite/2-2*app.rastermass, ye+g3+app.rastermass/2, app.breite/2+2*app.rastermass, ye+g3+g3)
	app.neuesSpielButton.SetzeEckradius(app.rastermass / 3)
	app.hilfeButton.SetzeKoordinaten(2*app.rastermass, ye+g3+app.rastermass/2, 4*app.rastermass, ye+g3+g3)
	app.hilfeButton.SetzeEckradius(app.rastermass / 3)
	app.hilfebox.SetzeKoordinaten(app.breite/4, app.hoehe/4, app.breite/2, app.hoehe/2)
	app.quitButton.SetzeKoordinaten(app.breite-4*app.rastermass, ye+g3+app.rastermass/2, app.breite-2*app.rastermass, ye+g3+g3)
	app.quitButton.SetzeEckradius(app.rastermass / 3)

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
				// zurück zum Spielmodus
				a.quizfenster.Ausblenden()
				a.billard.Starte()
				a.spieltisch.Einblenden()
			}
		} else {
			a.quiz.NaechsteFrage()
		}
	} else if a.spieltisch != nil && a.spieltisch.IstAktiv() &&
		!a.hilfebox.IstAktiv() && a.billard.Laeuft() {
		if a.billard.GibStrafpunkte() > a.billard.GibTreffer() {
			// stoppe die Zeit und gehe zum Quizmodus
			a.billard.Stoppe()
			a.spieltisch.Ausblenden()
			a.quiz.NaechsteFrage()
			a.quizfenster.Einblenden()
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
	if a.laeuft {
		return
	}
	println("Willkommen bei BrainPool")
	a.billard.Starte()
	a.spieltisch.Einblenden()
	a.quizfenster.Ausblenden()
	a.hilfebox.Ausblenden()
	a.renderer.Starte()
	a.mausSteuerung = views_controls.NewMausRoutine(a.mausSteuerFunktion)
	a.mausSteuerung.Starte()
	a.geraeusche.StarteLoop()
	a.musik.StarteLoop()
	a.laeuft = true
	// ####### der Tastatur-Loop bestimmt das Laufende ##################################################
	for {
		if !gfx.FensterOffen() {
			break
		}
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'p': // Spiel pausieren
				a.billard.PauseAnAus()
			case 'c': // Dunkle Umgebung
				a.renderer.DarkmodeAnAus()
			case 'h': // Hilfe anzeigen
				a.hilfebox.DarstellenAnAus()
			case 'd': // Zeitlupe (Testzwecke)
				a.billard.ZeitlupeAnAus()
			case 'l': // Fenster-Layout anzeigen (Testzwecke)
				a.renderer.LayoutAnAus()
			case 'f': // Quizmodus händisch an-/ausschalten (Testzwecke)
				if a.quizfenster.IstAktiv() {
					a.quizfenster.Ausblenden()
					a.billard.Starte()
					a.spieltisch.Einblenden()
				} else {
					a.billard.Stoppe()
					a.spieltisch.Ausblenden()
					a.quizfenster.Einblenden()
				}
			case 'q':
				a.Quit()
				return
			}
		}
	}
}

// ####### stoppe alles ##################################################
func (a *bpapp) Quit() {
	if !a.laeuft {
		return
	}
	a.geraeusche.Stoppe()
	a.musik.Stoppe()
	a.renderer.UeberblendeText("Bye!", views_controls.Fanzeige(), views_controls.Ftext(), 30)
	go a.mausSteuerung.Stoppe()
	a.billard.Stoppe()
	a.renderer.Stoppe()
	time.Sleep(500 * time.Millisecond)
	println("BrainPool wird beendet")
	if gfx.FensterOffen() {
		println("Schließe Gfx-Fenster")
		gfx.FensterAus()
	}
}

// ####### der Startpunkt ##################################################
func main() {
	// das Rastermaß bestimmt die Größe der App
	NewBPApp(35).Run()
}
