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
	spielFenster    views_controls.Widget
	quizFenster     views_controls.Widget
	hilfeFenster    views_controls.Widget
	gameOverFenster views_controls.Widget
	hintergrund     views_controls.Widget
	klickbare       []views_controls.Widget
	renderer        views_controls.FensterZeichner
	// Controls
	mausSteuerung views_controls.EingabeRoutine
	umschalter    hilf.Routine
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

	a := bpapp{rastermass: g, breite: 32 * g, hoehe: 22 * g, klickbare: []views_controls.Widget{}}
	a.renderer = views_controls.NewFensterZeichner("BrainPool - Das MiniBillard für Schlaue.")

	a.musik = klaenge.CoolJazz2641SOUND()
	a.geraeusche = klaenge.BillardPubAmbienceSOUND()

	// ######## Modelle und Views zusammenstellen #################################
	// realer Tisch: 2540 mm x 1270 mm, Kugelradius: 57.2 mm
	// Breite, Höhe des Spielfelds (2:1)
	var bS uint16 = 3 * a.breite / 4
	var hS uint16 = bS / 2 // Anderes Seitenverhältnis geht auch.

	// Modelle erzeugen
	a.billard = modelle.NewMini9BallSpiel(bS, hS)
	a.quiz = modelle.NewQuizInformatiksysteme()
	//a.quiz = modelle.NewBeispielQuiz()

	// Views und Zeichner erzeugen
	a.hintergrund = views_controls.NewFenster()
	a.renderer.SetzeFensterHintergrund(a.hintergrund)
	punktezaehler := views_controls.NewMBPunkteAnzeiger(a.billard)
	restzeit := views_controls.NewMBRestzeitAnzeiger(a.billard)

	bande := views_controls.NewFenster()
	a.spielFenster = views_controls.NewMBSpieltisch(a.billard)
	a.quizFenster = views_controls.NewQuizFenster(a.quiz, func() { a.billard.ReduziereStrafpunkte(); a.quiz.NaechsteFrage() }, func() { a.quiz.NaechsteFrage() })
	a.hilfeFenster = views_controls.NewTextBox("Hilfe")
	a.hilfeFenster.Ausblenden() // wäre standardmäßig eingeblendet
	a.gameOverFenster = views_controls.NewTextBox("GAME OVER!")
	a.gameOverFenster.Ausblenden() // wäre standardmäßig eingeblendet

	//setze Layout
	a.hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = 4 * a.rastermass, 6 * a.rastermass, 28 * a.rastermass, 18 * a.rastermass
	var g3 uint16 = a.rastermass + a.rastermass/3
	punktezaehler.SetzeKoordinaten(xs-g3, 1*a.rastermass, 18*a.rastermass, 3*a.rastermass)
	restzeit.SetzeKoordinaten(20*a.rastermass+g3, a.rastermass, xe+g3, 3*a.rastermass)
	bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	bande.SetzeEckradius(g3)
	a.spielFenster.SetzeKoordinaten(xs, ys, xe, ye)
	a.quizFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.quizFenster.SetzeEckradius(g3 - 2)
	a.hilfeFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.hilfeFenster.SetzeEckradius(g3 - 2)
	a.gameOverFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.gameOverFenster.SetzeEckradius(g3 - 2)

	// Buttonleiste
	hilfeButton := views_controls.NewButton("(h)ilfe an/aus", a.hilfeAnAus)
	neuesSpielButton := views_controls.NewButton("(n)eues Spiel", a.neuesSpielStarten)
	pauseButton := views_controls.NewButton("(m)usik spielen", a.musik.StarteLoop)
	darkButton := views_controls.NewButton("(d)unkel/hell", a.renderer.DarkmodeAnAus)
	quitButton := views_controls.NewButton("(q)uit", a.Quit)

	a.klickbare = []views_controls.Widget{hilfeButton, neuesSpielButton, pauseButton, darkButton, quitButton}
	zb := (a.breite - 2*a.rastermass) / uint16(len(a.klickbare))
	for i, k := range a.klickbare {
		k.SetzeKoordinaten(a.rastermass+uint16(i)*zb+zb/8, ye+5*a.rastermass/2, a.rastermass+uint16(i+1)*zb-zb/8, ye+13*a.rastermass/4)
		k.SetzeEckradius(a.rastermass / 3)
		k.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	}

	//setzeFarben
	a.hintergrund.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	a.spielFenster.SetzeFarben(views_controls.Fbillardtuch(), views_controls.Fdiamanten())
	bande.SetzeFarben(views_controls.Ftext(), views_controls.Fanzeige())
	punktezaehler.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	punktezaehler.SetzeTransparenz(255)
	restzeit.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.quizFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	a.hilfeFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	a.gameOverFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.klickbare = append(a.klickbare, a.quizFenster) // Quizzes sind auch klickbar
	a.renderer.SetzeWidgets(bande, a.spielFenster, a.quizFenster, punktezaehler, restzeit,
		hilfeButton, neuesSpielButton, pauseButton, darkButton, quitButton,
		a.gameOverFenster, a.hilfeFenster)

	return &a
}

// ############### Regele die Umschaltung zwischen den App-Modi #######################
/*
Vor.: das Hilfefenster liegt zuoberst im Renderer
Eff.: zeigt das Hilfefenster an oder blendet es wieder aus. Billard wird angehalten.
*/
func (a *bpapp) hilfeAnAus() {
	if a.hilfeFenster.IstAktiv() {
		a.hilfeFenster.Ausblenden()
		if a.spielFenster.IstAktiv() {
			a.billard.Starte()
		}
	} else {
		a.renderer.UeberblendeAus()
		a.hilfeFenster.Einblenden()
		if a.spielFenster.IstAktiv() {
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
	a.spielFenster.Einblenden()
	if !a.billard.Laeuft() {
		a.billard.Starte()
	}
}

/*
Zweck: die sonstige Umschaltung zwischen Quiz und Spiel(wird als go-Routine ausgelagert)
Vor.: keine
Eff.:

	Falls Spielzeit abgelaufen war: Spiel wird beendet.
	Falls Anzahl Fouls >> Anzahl Treffer: Quiz ist aktiviert.
	Falls Anzahl Fouls < Anzahl Treffer oder 0: Spiel ist aktiviert
*/
func (a *bpapp) quizUmschalterFunktion() {
	if a.spielFenster.IstAktiv() && a.billard.GibRestzeit() == 0 {
		a.billard.Stoppe()
		a.spielFenster.Ausblenden()
		a.gameOverFenster.Einblenden() // Hier ist Ende; man muss ein neues Spiel starten ...
	} else if a.spielFenster.IstAktiv() && a.billard.GibStrafpunkte() > a.billard.GibTreffer()+2 {
		// stoppe die Zeit und gehe zum Quizmodus
		a.billard.Stoppe()
		a.spielFenster.Ausblenden()
		a.quiz.NaechsteFrage()
		a.quizFenster.Einblenden()
	} else if a.quizFenster.IstAktiv() && (a.billard.GibStrafpunkte() == 0 || a.billard.GibStrafpunkte() < a.billard.GibTreffer()) {
		// zurück zum Spielmodus
		a.quizFenster.Ausblenden()
		a.billard.Starte()
		a.spielFenster.Einblenden()
	}
}

/*
Zweck: Die App wird mit der Maus bedient. Die Maussteuerung wird als go-Routine ausgelagert.
Vor.: keine
Eff.:

	Falls Quit, Hilfe oder Neues Spiel angeklickt wurde, ist die enstprechende Aktion ausgeführt.
	Falls das Quiz aktiv ist und das Quizfenster angeklickt wurde, wird die Antwort ausgewertet.
	  -- Falls dadurch genügend Strafpunkte abgebaut wurden, wird das Spiel aktiviert.
	Falls das Spiel läuft und alle Kugeln still stehen: Es ist der Queue bewegt bzw.
	ist der Stoss erfolgt.
	Sonst: keiner
*/
func (a *bpapp) mausSteuerFunktion(taste uint8, status int8, mausX, mausY uint16) {
	if taste == 1 && status == -1 { // es wurde links geklickt
		// anklickbare Widgets abfragen
		var benutzt bool
		for _, b := range a.klickbare {
			if b.IstAktiv() && b.ImFenster(mausX, mausY) {
				b.MausklickBei(mausX, mausY)
				benutzt = true
				break
			}
		}
		// falls niemand den Klick wollte, gib ihn ans Spiel
		if !benutzt && a.spielFenster.IstAktiv() {
			// kann auch außerhalb des Tuchs klicken
			a.spielFenster.MausklickBei(mausX, mausY)
		}
	} else { // es wurde nicht geklickt
		if a.spielFenster.IstAktiv() {
			// sonst: zielen und Kraft aufbauen
			switch taste {
			case 4: // vorwärts scrollen
				a.spielFenster.MausScrolltHoch()
			case 5: // rückwärts scrollen
				a.spielFenster.MausScrolltRunter()
			default: // bewegen
				a.spielFenster.MausBei(mausX, mausY)
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

	Taste 'h' gedrückt: Hilfe angezeigt/ausgeblendet
	Taste 'r' gedrückt: neues Spiel begonnen
	Taste 'p' gedrückt: Spiel pausiert
	Taste 'd' gedrückt: Umgebung abgedunkelt/aufgehellt
	Taste 'm' gedrückt: Musik startet
	Taste 'q' gedrückt: App ist beendet.

	(für Testzwecke)

	Taste 'z' gedrückt: Bewegungen erfolgen in Zeitlupe
	Taste 'l' gedrückt: Layout der App ist angezeigt
	Taste '1' gedrückt: Teste mit 1 Kugel
	Taste '3' gedrückt: Teste mit 3 Kugeln
	Taste '9' gedrückt: Spiele mit 9 Kugeln
*/
func (a *bpapp) Run() {
	if a.laeuft {
		return
	}
	println("Willkommen bei BrainPool")
	a.billard.Starte()
	a.spielFenster.Einblenden()
	a.quizFenster.Ausblenden()
	a.hilfeFenster.Ausblenden()
	a.renderer.Starte() // go-Routine
	a.mausSteuerung = views_controls.NewMausRoutine(a.mausSteuerFunktion)
	a.mausSteuerung.StarteRate(20)                                         // go-Routine
	a.umschalter = hilf.NewRoutine("Umschalter", a.quizUmschalterFunktion) // go-Routine
	a.umschalter.StarteRate(5)                                             // go-Routine
	a.geraeusche.StarteLoop()                                              // go-Routine
	a.laeuft = true

	// ####### der Tastatur-Loop darf hier existieren ####################
	for {
		if !gfx.FensterOffen() {
			break
		}
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 'h': // Hilfe an-aus
				a.hilfeAnAus()
			case 'n': // neues Spiel
				a.neuesSpielStarten()
			case 'p': // Spiel pausieren
				a.billard.PauseAnAus()
			case 'd': // Dunkle Umgebung
				a.renderer.DarkmodeAnAus()
			case 'm': // Musik spielen, wenn man möchte
				a.musik.StarteLoop() // go-Routine
			case 'q':
				a.Quit()
				return
				// ######  Testzwecke ####################################
			case 'z': // Zeitlupe
				a.billard.ZeitlupeAnAus()
			case 'l': // Fenster-Layout anzeigen
				a.renderer.LayoutAnAus()
			case '1': // Spiel testen
				a.billard.SetzeRestzeit(10 * time.Second)
				a.billard.SetzeKugeln1BallTest()
			case '3': // Spiel testen
				a.billard.SetzeSpielzeit(90 * time.Second)
				a.billard.SetzeKugeln3Ball()
			case '9': // Spiel testen
				a.billard.SetzeSpielzeit(4 * time.Minute)
				a.billard.SetzeKugeln9Ball()
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
	a.umschalter.Stoppe()
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
	// Die gewünschte Fensterbreite in Pixeln wird übergeben.
	// Das Seitenverhältnis des Spiels ist B:H = 16:11
	NewBPApp(1024).Run() // läuft bis Spiel beendet wird
}
