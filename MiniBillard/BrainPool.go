package main

import (
	"gfx"

	"./hilf"
	"./klaenge"
	"./modelle"
	"./views_controls"
)

var (
	rastermaß    uint16 = 35
	breite, höhe uint16 = 32 * rastermaß, 22 * rastermaß // Größe des gesamten App-Fensters
	// Klänge
	musik, geräusche klaenge.Klang = klaenge.CoolJazz2641SOUND(), klaenge.BillardPubAmbienceSOUND()
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
)

func mausSteuerFunktion(taste uint8, status int8, mausX, mausY uint16) {
	if quitButton != nil && quitButton.IstAktiv() &&
		quitButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		quitButton.MausklickBei(mausX, mausY)
		return
	} else if hilfeButton != nil && hilfeButton.IstAktiv() &&
		hilfeButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		hilfeButton.MausklickBei(mausX, mausY)
	} else if neuesSpielButton != nil && neuesSpielButton.IstAktiv() &&
		neuesSpielButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		neuesSpielButton.MausklickBei(mausX, mausY)
	} else if quizfenster != nil && quizfenster.IstAktiv() &&
		quizfenster.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		quizfenster.MausklickBei(mausX, mausY)
		if quiz.GibAktuelleFrage().RichtigBeantwortet() {
			billard.ReduziereStrafpunkte()
			if billard.GibStrafpunkte() <= billard.GibTreffer() {
				quizfenster.SetzeInAktiv()
				spieltisch.SetzeAktiv() // zurück zum Spielmodus
			}
		} else {
			quiz.NaechsteFrage()
		}
	} else if spieltisch != nil && spieltisch.IstAktiv() &&
		!hilfebox.IstAktiv() && billard.Laeuft() {
		if billard.GibStrafpunkte() > billard.GibTreffer() {
			quiz.NaechsteFrage()
			spieltisch.SetzeInAktiv()
			quizfenster.SetzeAktiv() // zum Quizmodus
		} else if billard.IstStillstand() {
			// zielen und stoßen
			switch taste {
			case 1: // stoßen
				if status == -1 {
					billard.Stosse()
				}
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

func Run() {
	billard.Starte()
	spieltisch.SetzeAktiv()
	quizfenster.SetzeInAktiv()
	hilfebox.SetzeInAktiv()
	renderer.Starte()
	mausSteuerung = views_controls.NewMausProzess(mausSteuerFunktion)
	mausSteuerung.Starte()
	geräusche.StarteLoop()
	musik.StarteLoop()
}

func Quit() {
	geräusche.Stoppe()
	musik.Stoppe()
	renderer.UeberblendeText("Bye!", views_controls.Fanzeige(), views_controls.Ftext(), 30)
	mausSteuerung.Stoppe()
	billard.Stoppe()
	renderer.Stoppe()
}

func initMV() {
	// ######## Modelle und Views zusammenstellen #################################
	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	// Breite, Höhe des Spielfelds
	var bS uint16 = 3 * breite / 4
	var hS uint16 = bS / 2
	// Radius der Kugeln
	var ra uint16 = uint16(0.5 + float64(bS)*57.2/2540)

	// Modelle erzeugen
	billard = modelle.NewMini9BallSpiel(bS, hS, ra)
	quiz = modelle.NewQuizCSV("BeispielQuiz.csv")

	// Views und Zeichner erzeugen
	renderer = views_controls.NewFensterZeichner("BrainPool - Das MiniBillard für Schlaue.")
	hintergrund = views_controls.NewFenster()
	renderer.SetzeFensterHintergrund(hintergrund)
	punktezaehler = views_controls.NewMBPunkteAnzeiger(billard)
	restzeit = views_controls.NewMBRestzeitAnzeiger(billard)
	bande = views_controls.NewFenster()
	spieltisch = views_controls.NewMBSpieltisch(billard)
	quizfenster = views_controls.NewQuizFenster(quiz)
	neuesSpielButton = views_controls.NewButton("neues Spiel",
		func() {
			billard.Reset()
			quizfenster.SetzeInAktiv()
			hilfebox.SetzeInAktiv()
			spieltisch.SetzeAktiv()
		})
	quitButton = views_controls.NewButton("quit",
		func() {
			Quit()
		})
	hilfebox = views_controls.NewTextBox("Hilfe")
	hilfebox.SetzeInAktiv()
	hilfeButton = views_controls.NewButton("?",
		func() {
			hilfebox.AktivAnAus()
		})
	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	renderer.SetzeWidgets(bande, spieltisch, quizfenster, hilfebox, punktezaehler, restzeit,
		neuesSpielButton, quitButton, hilfeButton)
}

func setzeLayout() {
	hintergrund.SetzeKoordinaten(0, 0, breite, höhe)
	var xs, ys, xe, ye uint16 = 4 * rastermaß, 6 * rastermaß, 28 * rastermaß, 18 * rastermaß
	var g3 uint16 = rastermaß + rastermaß/3
	punktezaehler.SetzeKoordinaten(xs-g3, 1*rastermaß, 18*rastermaß, 3*rastermaß)
	restzeit.SetzeKoordinaten(20*rastermaß+g3, rastermaß, xe+g3, 3*rastermaß)
	bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	bande.SetzeEckradius(g3)
	spieltisch.SetzeKoordinaten(xs, ys, xe, ye)
	quizfenster.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	quizfenster.SetzeEckradius(g3)
	neuesSpielButton.SetzeKoordinaten(breite/2-2*rastermaß, ye+g3+rastermaß/2, breite/2+2*rastermaß, ye+g3+g3)
	neuesSpielButton.SetzeEckradius(rastermaß / 3)
	hilfeButton.SetzeKoordinaten(2*rastermaß, ye+g3+rastermaß/2, 4*rastermaß, ye+g3+g3)
	hilfeButton.SetzeEckradius(rastermaß / 3)
	hilfebox.SetzeKoordinaten(breite/4, höhe/4, breite/2, höhe/2)
	quitButton.SetzeKoordinaten(breite-4*rastermaß, ye+g3+rastermaß/2, breite-2*rastermaß, ye+g3+g3)
	quitButton.SetzeEckradius(rastermaß / 3)
}

func setzeFarben() {
	hintergrund.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	spieltisch.SetzeFarben(views_controls.Fbillardtuch(), views_controls.Fdiamanten())
	bande.SetzeFarben(views_controls.Ftext(), views_controls.Fanzeige())
	punktezaehler.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	punktezaehler.SetzeTransparenz(255)
	restzeit.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	quizfenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	neuesSpielButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	hilfeButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	hilfebox.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	quitButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
}

func main() {
	println("Willkommen bei BrainPool")
	// Modelle und Views zusammenstellen
	initMV()
	// Abmessungen und Layout
	setzeLayout()
	// Farben
	setzeFarben()
	// Starte
	Run()
	// Tastatur-Loop
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'p': // Pause
				billard.PauseAnAus()
			case 'c': // dark mode
				renderer.DarkmodeAnAus()
			case 'h': // Hilfe anzeigen
				hilfebox.AktivAnAus()
			case 'd': // Zeitlupe (Testzwecke - zeigt Richtung und Geschwindigkeit der Kugeln)
				billard.ZeitlupeAnAus()
			case 'l': // Layout (Testzwecke zeigt nur Widget-Ränder und Farbschema - Spiel und Steuerung laufen weiter)
				renderer.LayoutAnAus()
			case 'f': // Quizmodus erzwingen (Testzwecke)
				spieltisch.SetzeInAktiv()
				quizfenster.SetzeAktiv()
			case 'q': // quit
				// ######## Stoppe alles #########################################
				Quit()
				return
			}
		}
	}
}
