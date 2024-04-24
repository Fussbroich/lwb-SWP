// Autoren:
//
//	Thomas Schrader
//	Bettina Chang
//
// Zweck:
//
//	Das Spielprogramm BrainPool -
//	ein Softwareprojekt im Rahmen der Lehrerweiterbildung Berlin
//
// Datum: 22.04.2024
package main

import (
	"time"

	"./hilf"
	"./klaenge"
	"./modelle"
	"./views_controls"
)

// Eine App ist eine grafische Anwendung, die im "unmittelbaren Modus" in einem einzigen
// Fenster läuft. Das bedeutet, dass nach jedem zeitlichen "Tick" (Zeiteinheit) das gesamte
// Fenster mit allen grafischen Elementen (Widgets) neu gezeichnet wird. Die Modelle sind
// Teil der App, und ihr Zustand wird jedesmal neu abgefragt.
//
//	Vor.: Das Grafikpaket gfx muss im GOPATH installiert sein.
//
//	Erzeugung eines BrainPool-Spiels (Bislang einzige Verwendung ;-) mit
//	NewBPApp(uint16) - erhält die gewünschte Fensterbreite in Pixeln
type App interface {
	// Startet die Laufzeit-Elemente der App.
	//
	//	Vor.: die App läuft nicht
	//	Eff.: Die App wurde gestartet und ein gfx-Fenster geöffnet.
	Run()

	// Stoppt die Laufzeit-Elemente der App auf geregelte Art und Weise.
	//
	//	Vor.: die App läuft
	//	Eff.: Die App wurde beendet und das gfx-Fenster geschlossen.
	Quit()
}

type bpapp struct {
	laeuft bool
	quit   bool
	// Größe des gesamten App-Fensters
	breite uint16
	hoehe  uint16
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
	quitButton      views_controls.Widget
	buttonLeiste    []views_controls.Widget
	// Verwalter für den Zeichen-Loop
	renderer views_controls.FensterZeichner
	// separates Eingabe-Control
	mausSteuerung   views_controls.EingabeRoutine
	tastenSteuerung views_controls.EingabeRoutine
	// Spielmodus-Umschalter
	umschalter hilf.Routine
}

// Zweck: Konstruktor für BrainPool - baut eine App zusammen
//
//	Vor.:  keine
//	Eff.:  Ein App-Objekt steht zum Starten bereit.
func NewBPApp(b uint16) *bpapp {
	if b > 1920 {
		b = 1920 // größtmögliches gfx-Fenster ist 1920 Pixel breit
	}
	if b < 640 {
		b = 640 // kleinere Bildschirme sind zum Spielen ungeeignet
	}

	var hilfetext string = "Hilfe\n\n" +
		"Im Spielmodus (und nur, wenn alle Kugeln still stehen): " +
		"Maus bewegen ändert die Zielrichtung. Stoß durch klicken mit der linken Maustaste. " +
		"Die Stoßkraft wird durch scrollen der Maus verändert.\n\n" +
		"Du spielst gegen die Zeit. Alle neun Kugel müssen versenkt werden. " +
		"Es gibt ein Foul, wenn die weiße Kugel reingeht oder wenn bei einem Stoß gar keine Kugel versenkt wird.\n\n" +
		"Im Quizmodus: Klicke die richtigen Antworten an, um Fouls abzuarbeiten.\n\n" +
		"Die übrige Bedienung erfolgt durch Anklicken der Buttons unten " +
		"oder mit der angegebenen Taste auf der Tastatur."

	var g uint16 = b / 32 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 16:11
	a := bpapp{breite: 32 * g, hoehe: 22 * g, buttonLeiste: []views_controls.Widget{}}

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
	punktezaehler := views_controls.NewMBPunkteAnzeiger(a.billard)
	restzeit := views_controls.NewMBRestzeitAnzeiger(a.billard)

	bande := views_controls.NewFenster()
	a.spielFenster = views_controls.NewMBSpieltisch(a.billard)
	a.quizFenster = views_controls.NewQuizFenster(a.quiz, func() { a.billard.ReduziereStrafpunkte(); a.quiz.NaechsteFrage() }, func() { a.quiz.NaechsteFrage() })
	a.hilfeFenster = views_controls.NewTextBox(hilfetext, 18)
	a.hilfeFenster.Ausblenden() // wäre standardmäßig eingeblendet
	a.gameOverFenster = views_controls.NewTextBox("GAME OVER!\n\n\nStarte ein neues Spiel.", 24)
	a.gameOverFenster.Ausblenden() // wäre standardmäßig eingeblendet

	a.renderer = views_controls.NewFensterZeichner()
	a.renderer.SetzeFensterHintergrund(a.hintergrund)
	a.renderer.SetzeFensterTitel("BrainPool - Das MiniBillard für Schlaue.")

	//setze Layout
	a.hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = 4 * g, 6 * g, 28 * g, 18 * g
	var g3 uint16 = g + g/3
	punktezaehler.SetzeKoordinaten(xs-g3, 1*g, 18*g, 3*g)
	restzeit.SetzeKoordinaten(20*g+g3, g, xe+g3, 3*g)
	bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	bande.SetzeEckradius(g3)
	a.spielFenster.SetzeKoordinaten(xs, ys, xe, ye)
	a.quizFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.quizFenster.SetzeEckradius(g3 - 2)
	a.hilfeFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.hilfeFenster.SetzeEckradius(g3 - 2)
	a.gameOverFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.gameOverFenster.SetzeEckradius(g3 - 2)

	//setzeFarben
	a.hintergrund.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	a.spielFenster.SetzeFarben(views_controls.Fbillardtuch(), views_controls.Fdiamanten())
	bande.SetzeFarben(views_controls.Fbande(), views_controls.Fanzeige())
	punktezaehler.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	punktezaehler.SetzeTransparenz(255)
	restzeit.SetzeFarben(views_controls.Fanzeige(), views_controls.Ftext())
	a.quizFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	a.hilfeFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())
	a.gameOverFenster.SetzeFarben(views_controls.Fquiz(), views_controls.Ftext())

	// Buttonleiste
	hilfeButton := views_controls.NewButton("(h)ilfe", a.HilfeAnAus)
	neuesSpielButton := views_controls.NewButton("(n)eues Spiel", a.NeuesSpielStarten)
	pauseButton := views_controls.NewButton("(m)usik spielen", a.MusikAn)
	darkButton := views_controls.NewButton("(d)unkel/hell", a.DarkmodeAnAus)
	a.quitButton = views_controls.NewButton("(q)uit", a.Quit)
	a.buttonLeiste = []views_controls.Widget{hilfeButton, neuesSpielButton, pauseButton, darkButton, a.quitButton}

	zb := (a.breite - 2*g) / uint16(len(a.buttonLeiste))
	for i, b := range a.buttonLeiste {
		b.SetzeKoordinaten(g+uint16(i)*zb+zb/8, ye+5*g/2, g+uint16(i+1)*zb-zb/8, ye+13*g/4)
		b.SetzeEckradius(g / 3)
		b.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	}

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.renderer.SetzeWidgets(bande, a.spielFenster, a.quizFenster, punktezaehler, restzeit,
		hilfeButton, a.hilfeFenster, neuesSpielButton, pauseButton, darkButton, a.quitButton,
		a.gameOverFenster)

	return &a
}

// #### Regele die Steuerung der App und die Umschaltung zwischen den App-Modi ##################

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: zeigt das Hilfefenster an oder blendet es wieder aus. Das Spielmodell wird solang angehalten.
func (a *bpapp) HilfeAnAus() {
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

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: neues Spiel ist gestartet - alle anderen Fenster sind ausgeblendet
func (a *bpapp) NeuesSpielStarten() {
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

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: die hinterlegte Musik startet im Loop
func (a *bpapp) MusikAn() {
	a.musik.StarteLoop()
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: die GUI wird zwischen hell und dunkel umgeschaltet
func (a *bpapp) DarkmodeAnAus() {
	a.renderer.DarkmodeAnAus()
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: das Spiel (und der Countdown) hält an, bzw. läuft weiter
func (a *bpapp) PauseAnAus() {
	a.billard.PauseAnAus()
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: das Spiel (und der Countdown) laufen in Zeitlupe.
func (a *bpapp) ZeitlupeAnAus() {
	a.billard.ZeitlupeAnAus()
}

// Umschalter zwischen den App-Zuständen (wird als go-Routine ausgelagert)
// Regeln die Umschaltung zwischen Quiz und Spiel gemäß der Regeln.
//
//	Vor.: keine
//	Eff.:
//	Falls Spielzeit abgelaufen war: Spiel wird beendet.
//	Falls Anzahl Fouls >> Anzahl Treffer: Quiz ist aktiviert.
//	Falls Anzahl Fouls < Anzahl Treffer oder 0: Spiel ist aktiviert
func (a *bpapp) quizUmschalterFunktion() {
	tr, st, rz := a.billard.GibTreffer(), a.billard.GibStrafpunkte(), a.billard.GibRestzeit()
	if a.spielFenster.IstAktiv() &&
		rz == 0 {
		//Spielzeit abgelaufen
		a.billard.Stoppe()
		a.quizFenster.Ausblenden()
		a.spielFenster.Ausblenden()
		a.gameOverFenster.Einblenden()
	} else if a.spielFenster.IstAktiv() &&
		st > tr+2 { // Es sind zu viele Strafpunkte
		// zum Quizmodus
		a.billard.Stoppe()
		a.spielFenster.Ausblenden()
		a.quiz.NaechsteFrage()
		a.quizFenster.Einblenden()
	} else if a.quizFenster.IstAktiv() &&
		(st == 0 || st < tr) { // Strafpunkte sind abgebaut
		// zurück zum Spielmodus
		a.quizFenster.Ausblenden()
		a.billard.Starte()
		a.spielFenster.Einblenden()
	}
}

// Die Maussteuerung der App (kann als go-Routine in einem Loop laufen).
//
//	Vor.: keine
//	Eff.: Gibt einige der möglichen Mausaktionen an passende Widgets weiter.
//	Sonst: keiner
func (a *bpapp) mausSteuerFunktion(taste uint8, status int8, mausX, mausY uint16) {
	if taste == 1 && status == -1 { // es wurde links geklickt
		// Buttonleiste abfragen
		for _, b := range a.buttonLeiste {
			if b.IstAktiv() && b.ImFenster(mausX, mausY) {
				b.MausklickBei(mausX, mausY)
				return
			}
		}
		// schauen, ob das Quiz angeklickt wurde
		if a.quizFenster.IstAktiv() && a.quizFenster.ImFenster(mausX, mausY) {
			a.quizFenster.MausklickBei(mausX, mausY)
			return
		}
		// Falls bislang niemand den Klick wollte, gib ihn ans Spiel.
		// (Zum Spielen kann man auch außerhalb des Spieltisches klicken, daher die Sonderbehandlung ...)
		if a.spielFenster.IstAktiv() {
			// kann auch außerhalb des Tuchs klicken
			a.spielFenster.MausklickBei(mausX, mausY)
		}
	} else { // es wurde gar nicht geklickt
		if a.spielFenster.IstAktiv() {
			// zielen und Kraft aufbauen
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

// Die Tastatursteuerung der App.
//
//	Vor: keine
//	Eff.: die zur Taste passende Spiel-Aktion ist ausgeführt.
//	Erg.: Soll der Afrufer die Abfrage beenden (Quit) true, sonst false.
func (a *bpapp) tastenSteuerFunktion(taste uint16, gedrückt uint8, _ uint16) {
	if gedrückt == 1 {
		switch taste {
		case 'h': // Hilfe an-aus
			a.HilfeAnAus()
		case 'n': // neues Spiel
			a.NeuesSpielStarten()
		case 'p': // Spiel pausieren
			a.PauseAnAus()
		case 'd': // Dunkle Umgebung
			a.DarkmodeAnAus()
		case 'm': // Musik spielen, wenn man möchte
			a.MusikAn() // go-Routine
		case 'q':
			a.Quit()
			// ######  Testzwecke ####################################
		case 's': // Zeitlupe
			a.ZeitlupeAnAus()
		case 'l': // Fenster-Layout anzeigen
			a.renderer.LayoutAnAus()
		case 'e': // Spiel testen
			a.billard.ErhoeheStrafpunkte()
		case 'r': // Spiel testen
			a.billard.ReduziereStrafpunkte()
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

// Startet die Laufzeit-Elemente der BrainPool App.
//
//	Vor.: die App läuft nicht
//	Eff.: Die App wurde gestartet und ein gfx-Fenster geöffnet.
func (a *bpapp) Run() {
	if a.laeuft {
		return
	}
	println("Willkommen bei BrainPool")
	a.billard.Starte() // Modell bereit zum Spielen
	a.spielFenster.Einblenden()
	a.quizFenster.Ausblenden()
	a.hilfeFenster.Ausblenden()
	a.gameOverFenster.Ausblenden()
	a.geraeusche.StarteLoop() // go-Routine
	a.renderer.Starte()       // go-Routine
	a.laeuft = true

	a.mausSteuerung = views_controls.NewMausRoutine(a.mausSteuerFunktion)
	a.mausSteuerung.StarteRate(50) // go-Routine mit begrenzter Rate

	//  ####### der eigentliche Event-Loop der App läuft nebenher #############
	a.umschalter = hilf.NewRoutine("Umschalter", a.quizUmschalterFunktion)
	a.umschalter.StarteRate(20) // go-Routine mit begrenzter Rate

	// Dafür darf der Tastatur-Loop hier existieren
	a.tastenSteuerung = views_controls.NewTastenRoutine(a.tastenSteuerFunktion)
	a.tastenSteuerung.StarteHier()
}

// Stoppt die Laufzeit-Elemente der BrainPool App auf geregelte Art und Weise.
//
//	Vor.: die App läuft
//	Eff.: Die App wurde beendet und das gfx-Fenster geschlossen.
func (a *bpapp) Quit() {
	if !a.laeuft {
		return
	}
	a.quit = true
	a.geraeusche.Stoppe()
	a.musik.Stoppe()
	a.renderer.UeberblendeText("Bye!", views_controls.Fanzeige(), views_controls.Ftext(), 30)
	go a.mausSteuerung.Stoppe()   // go-Routine
	go a.tastenSteuerung.Stoppe() // go-Routine
	a.umschalter.Stoppe()
	a.billard.Stoppe()
	time.Sleep(500 * time.Millisecond)
	println("BrainPool wird beendet")
	a.renderer.Stoppe()
}

// ####### der Startpunkt ##################################################
func main() {
	// Die gewünschte Fensterbreite in Pixeln wird übergeben.
	// Das Seitenverhältnis des Spiels ist B:H = 16:11
	NewBPApp(1024).Run() // läuft bis Spiel beendet wird
}
