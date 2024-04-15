/*
Autoren:

	Thomas Schrader
	Bettina Chang

Zweck:

	Softwareprojekt im Rahmen der Lehrerweiterbildung Berlin

Datum: 15.04.2024
*/
package main

import (
	"gfx"
	"time"

	"./hilf"
	"./klaenge"
	"./modelle"
	"./views_controls"
)

// Ein Objekt, das die ganze App initialisiert und steuert
type BPApp interface {
	/*
	   Vor.: Keine
	   Eff.: Die App wurde gestartet und ein gfx-Fenster geöffnet.
	*/
	Run()
	/*
	   Vor.: Keine
	   Eff.: Die App wurde beendet und das gfx-Fenster geschlossen.
	*/
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
	quizFenster      views_controls.Widget
	hilfeFenster     views_controls.Widget
	gameOverFenster  views_controls.Widget
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
	regelWaechter hilf.Routine
}

/*
Zweck: baut die App zusammen
Vor.:  keine
Eff.:  Die App wurde beendet und das gfx-Fenster geschlossen.
*/
func NewBPApp(b uint16) *bpapp {
	if b > 1920 {
		b = 1920 // größtmögliches gfx-Fenster ist 1920 Pixel breit
	}
	if b < 640 {
		b = 640 // kleinere Bildschirme sind zum Spielen ungeeignet
	}
	var g uint16 = b / 32 // Rastermass für dieses App-Design

	a := bpapp{rastermass: g, breite: 32 * g, hoehe: 22 * g}
	a.renderer = views_controls.NewFensterZeichner("BrainPool - Das MiniBillard für Schlaue.")

	a.musik = klaenge.CoolJazz2641SOUND()
	a.geraeusche = klaenge.BillardPubAmbienceSOUND()

	// ######## Modelle und Views zusammenstellen #################################
	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	// Breite, Höhe des Spielfelds
	var bS uint16 = 3 * a.breite / 4
	var hS uint16 = bS / 2
	// Radius der Kugeln
	var ra uint16 = uint16(0.5 + float64(bS)*57.2/2540)

	// Modelle erzeugen
	a.billard = modelle.NewMini9BallSpiel(bS, hS, ra)
	a.billard.SetzeRegeln(func() {
		if a.billard.GibSpielkugel().IstEingelocht() {
			a.billard.ErhoeheStrafpunkte()
			a.billard.StossWiederholen()
			return
		}
	})

	//a.quiz = modelle.NewQuizInformatiksysteme()
	a.quiz = modelle.NewBeispielQuiz()

	// Views und Zeichner erzeugen
	a.hintergrund = views_controls.NewFenster()
	a.renderer.SetzeFensterHintergrund(a.hintergrund)
	a.punktezaehler = views_controls.NewMBPunkteAnzeiger(a.billard)
	a.restzeit = views_controls.NewMBRestzeitAnzeiger(a.billard)
	a.bande = views_controls.NewFenster()
	a.spieltisch = views_controls.NewMBSpieltisch(a.billard)
	a.quizFenster = views_controls.NewQuizFenster(a.quiz)
	a.neuesSpielButton = views_controls.NewButton("neues Spiel", a.neuesSpielStarten)
	a.quitButton = views_controls.NewButton("Quit", a.Quit)
	a.hilfeFenster = views_controls.NewTextBox("Hilfe")
	a.hilfeFenster.Ausblenden() // wäre standardmäßig eingeblendet
	a.hilfeButton = views_controls.NewButton("?", a.hilfeAnAus)
	a.gameOverFenster = views_controls.NewTextBox("GAME OVER!")
	a.gameOverFenster.Ausblenden() // wäre standardmäßig eingeblendet

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.renderer.SetzeWidgets(a.bande, a.spieltisch, a.quizFenster, a.punktezaehler, a.restzeit,
		a.neuesSpielButton, a.quitButton, a.hilfeButton, a.gameOverFenster, a.hilfeFenster)

	//setze Layout
	a.hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = 4 * a.rastermass, 6 * a.rastermass, 28 * a.rastermass, 18 * a.rastermass
	var g3 uint16 = a.rastermass + a.rastermass/3
	a.punktezaehler.SetzeKoordinaten(xs-g3, 1*a.rastermass, 18*a.rastermass, 3*a.rastermass)
	a.restzeit.SetzeKoordinaten(20*a.rastermass+g3, a.rastermass, xe+g3, 3*a.rastermass)
	a.bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	a.bande.SetzeEckradius(g3)
	a.spieltisch.SetzeKoordinaten(xs, ys, xe, ye)
	a.quizFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.quizFenster.SetzeEckradius(g3 - 2)
	a.neuesSpielButton.SetzeKoordinaten(a.breite/2-2*a.rastermass, ye+g3+a.rastermass/2, a.breite/2+2*a.rastermass, ye+g3+g3)
	a.neuesSpielButton.SetzeEckradius(a.rastermass / 3)
	a.hilfeButton.SetzeKoordinaten(2*a.rastermass, ye+g3+a.rastermass/2, 4*a.rastermass, ye+g3+g3)
	a.hilfeButton.SetzeEckradius(a.rastermass / 3)
	a.hilfeFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.hilfeFenster.SetzeEckradius(g3 - 2)
	a.quitButton.SetzeKoordinaten(a.breite-4*a.rastermass, ye+g3+a.rastermass/2, a.breite-2*a.rastermass, ye+g3+g3)
	a.quitButton.SetzeEckradius(a.rastermass / 3)
	a.gameOverFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.gameOverFenster.SetzeEckradius(g3 - 2)

	//setzeFarben
	a.hintergrund.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	a.spieltisch.SetzeFarben(views_controls.Fbillardtuch(), views_controls.Fdiamanten())
	a.bande.SetzeFarben(views_controls.Ftext(), views_controls.Fanzeige())
	a.punktezaehler.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.punktezaehler.SetzeTransparenz(255)
	a.restzeit.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.quizFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	a.neuesSpielButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.hilfeButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.hilfeFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	a.quitButton.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.gameOverFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())

	return &a
}

// ############### Regele die Umschaltung zwischen den App-Modi #######################
// Todo: Bündele Steuerung in einem Zustandsautomaten

/*
Vor.: das Hilfefenster liegt zuoberst im Renderer
Eff.: zeigt das Hilfefenster an oder blendet es wieder aus
*/
func (a *bpapp) hilfeAnAus() {
	if a.hilfeFenster.IstAktiv() {
		a.hilfeFenster.Ausblenden()
		if a.spieltisch.IstAktiv() {
			a.billard.Starte()
		}
	} else {
		a.renderer.UeberblendeAus()
		a.hilfeFenster.Einblenden()
		if a.spieltisch.IstAktiv() {
			a.billard.Stoppe()
		}
	}
}

/*
Vor.: keine
Eff.: neues Spiel ist gestartet - Quiz ist ausgeblendet
*/
func (a *bpapp) neuesSpielStarten() {
	a.renderer.UeberblendeAus()
	a.quizFenster.Ausblenden()
	a.hilfeFenster.Ausblenden()
	a.gameOverFenster.Ausblenden()
	a.billard.Reset()
	a.spieltisch.Einblenden()
	if !a.billard.Laeuft() {
		a.billard.Starte()
	}
}

/*
Zweck: die restliche Zustands-Umschaltung (wird als go-Routine ausgelagert)
Vor.: keine
Eff.:

	Falls Spielzeit abgelaufen war: Spiel ist beendet.
	Falls Anzahl Fouls > Anzahl Treffer: Quiz ist aktiviert.
	Falls Anzahl Fouls <= Anzahl Treffer: Spiel ist aktiviert
*/
func (a *bpapp) quizUmschalterFunktion() {
	if a.spieltisch.IstAktiv() && a.billard.GibRestzeit() == 0 {
		a.billard.Stoppe()
		a.spieltisch.Ausblenden()
		a.gameOverFenster.Einblenden() // Hier ist Ende; man muss ein neues Spiel starten ...
	} else if a.spieltisch.IstAktiv() && a.billard.GibStrafpunkte() > a.billard.GibTreffer() {
		// stoppe die Zeit und gehe zum Quizmodus
		a.billard.Stoppe()
		a.spieltisch.Ausblenden()
		a.quiz.NaechsteFrage()
		a.quizFenster.Einblenden()
	} else if a.quizFenster.IstAktiv() && a.billard.GibStrafpunkte() <= a.billard.GibTreffer() {
		// zurück zum Spielmodus
		a.quizFenster.Ausblenden()
		a.billard.Starte()
		a.spieltisch.Einblenden()
	}
}

/*
Zweck: Die App wird mit der Maus bedient. Die Maussteuerung wird als go-Routine ausgelagert.
Vor.: keine
Eff.:

	Falls Quit, Hilfe oder Neues Spiel angeklickt wurde, ist die enstprechende Aktion ausgeführt.
	Falls Quiz aktiv ist und Quizfenster angeklickt wurde, wird Antwort ausgewertet.
	  -- Falls dadurch genügend Strafpunkte abgebaut wurden, wird Spiel aktiviert.
	Falls Spiel läuft und alle Kugeln still stehen: Es ist gezielt bzw. gestoßen.
	Sonst: keiner
*/
func (a *bpapp) mausSteuerFunktion(taste uint8, status int8, mausX, mausY uint16) {
	if a.quitButton != nil && a.quitButton.IstAktiv() &&
		a.quitButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.quitButton.MausklickBei(mausX, mausY)
	} else if a.hilfeButton != nil && a.hilfeButton.IstAktiv() &&
		a.hilfeButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.hilfeButton.MausklickBei(mausX, mausY)
	} else if a.neuesSpielButton != nil && a.neuesSpielButton.IstAktiv() &&
		a.neuesSpielButton.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.neuesSpielButton.MausklickBei(mausX, mausY)
	} else if a.quizFenster.IstAktiv() &&
		a.quizFenster.ImFenster(mausX, mausY) && taste == 1 && status == -1 {
		a.quizFenster.MausklickBei(mausX, mausY)
		// Todo: Hier werden Regeln und Maussteuerung vermischt ...
		if a.quiz.GibAktuelleFrage().RichtigBeantwortet() {
			a.billard.ReduziereStrafpunkte()
		} else {
			a.quiz.NaechsteFrage()
		}
	} else if a.spieltisch.IstAktiv() && a.billard.Laeuft() {
		if a.billard.IstStillstand() {
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

/*
Zweck: starte die Laufzeit-Elemente der App
Vor.: die App läuft nicht
Eff.: die App läuft

Ferner:
Zweck: Einige Aspekte der App können mit der Tastatur bedient werden.
Eff.:

	Taste 'p' gedrückt: Spiel pausiert
	Taste 'c' gedrückt: Umgebung abgedunkelt/aufgehellt
	Taste 'h' gedrückt: Hilfe angezeigt/ausgeblendet
	Taste 'r' gedrückt: neues Spiel begonnen
	Taste 'd' gedrückt: Bewegungen erfolgen in Zeitlupe (Testzwecke)
	Taste 'l' gedrückt: Layout der App ist angezeigt (Testzwecke)
	Taste 'q' gedrückt: App ist beendet.
	Sonst: keiner
*/
func (a *bpapp) Run() {
	if a.laeuft {
		return
	}
	println("Willkommen bei BrainPool")
	a.billard.Starte()
	a.spieltisch.Einblenden()
	a.quizFenster.Ausblenden()
	a.hilfeFenster.Ausblenden()
	a.renderer.Starte() // go-Routine
	a.mausSteuerung = views_controls.NewMausRoutine(a.mausSteuerFunktion)
	a.mausSteuerung.StarteRate(20)                                            // go-Routine
	a.regelWaechter = hilf.NewRoutine("Umschalter", a.quizUmschalterFunktion) // go-Routine
	a.regelWaechter.StarteRate(5)                                             // go-Routine
	a.geraeusche.StarteLoop()                                                 // go-Routine
	a.musik.StarteLoop()                                                      // go-Routine
	a.laeuft = true

	// ####### der Tastatur-Loop darf hier existieren ####################
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
			case 'h': // Hilfe an-aus
				a.hilfeAnAus()
			case 'r': // neues Spiel
				a.neuesSpielStarten()
			case 'd': // Zeitlupe (Testzwecke)
				a.billard.ZeitlupeAnAus()
			case 'l': // Fenster-Layout anzeigen (Testzwecke)
				a.renderer.LayoutAnAus()
			case 't': // Spiel testen
				a.billard.Stoppe()
				a.billard.SetzeRestzeit(10 * time.Second)
				a.billard.SetzeKugelnTest()
				a.billard.Starte()
			case 'q':
				a.Quit()
				return
			}
		}
	}
}

/*
Zweck: stoppt die Laufzeit-Elemente der App
Vor.: die App läuft
Eff.: die App läuft nicht
*/
func (a *bpapp) Quit() {
	if !a.laeuft {
		return
	}
	a.geraeusche.Stoppe()
	a.musik.Stoppe()
	a.renderer.UeberblendeText("Bye!", views_controls.Fanzeige(), views_controls.Ftext(), 30)
	go a.mausSteuerung.Stoppe()
	a.regelWaechter.Stoppe()
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
	// Die gewünschte Bildbreite in Pixeln wird übergeben
	NewBPApp(1440).Run() // läuft bis Spiel beendet wird
}
