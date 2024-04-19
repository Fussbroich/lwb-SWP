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
// Datum: 19.04.2024
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

	HilfeAnAus()
	NeuesSpielStarten()
	PauseAnAus()
	DarkmodeAnAus()
	MusikAn()
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

// Zweck: Konstruktor - baut eine App zusammen
//
//	Vor.:  keine
//	Eff.:  Die App steht zum Starten bereit.
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
		"Die übrige Bedienung erfolgt mit den Buttons unten."

	var g uint16 = b / 32 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 16:11
	a := bpapp{rastermass: g, breite: 32 * g, hoehe: 22 * g, klickbare: []views_controls.Widget{}}

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
	a.hilfeFenster = views_controls.NewTextBox(hilfetext)
	a.hilfeFenster.Ausblenden() // wäre standardmäßig eingeblendet
	a.gameOverFenster = views_controls.NewTextBox("GAME OVER!\n\n\nStarte ein neues Spiel.")
	a.gameOverFenster.Ausblenden() // wäre standardmäßig eingeblendet

	a.renderer = views_controls.NewFensterZeichner()
	a.renderer.SetzeFensterHintergrund(a.hintergrund)

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

	// Buttonleiste
	hilfeButton := views_controls.NewButton("(h)ilfe an/aus", a.HilfeAnAus)
	neuesSpielButton := views_controls.NewButton("(n)eues Spiel", a.NeuesSpielStarten)
	pauseButton := views_controls.NewButton("(m)usik spielen", a.MusikAn)
	darkButton := views_controls.NewButton("(d)unkel/hell", a.DarkmodeAnAus)
	quitButton := views_controls.NewButton("(q)uit", a.Quit)

	a.klickbare = []views_controls.Widget{hilfeButton, neuesSpielButton, pauseButton, darkButton, quitButton}
	zb := (a.breite - 2*a.rastermass) / uint16(len(a.klickbare))
	for i, k := range a.klickbare {
		k.SetzeKoordinaten(a.rastermass+uint16(i)*zb+zb/8, ye+5*a.rastermass/2, a.rastermass+uint16(i+1)*zb-zb/8, ye+13*a.rastermass/4)
		k.SetzeEckradius(a.rastermass / 3)
		k.SetzeFarben(views_controls.Fhintergrund(), views_controls.Ftext())
	}

	// Quizzes sind auch klickbar, aber das Spielfenster wird besonders behandelt
	a.klickbare = append(a.klickbare, a.quizFenster)
	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.renderer.SetzeWidgets(bande, a.spielFenster, a.quizFenster, punktezaehler, restzeit,
		hilfeButton, a.hilfeFenster, neuesSpielButton, pauseButton, darkButton, quitButton,
		a.gameOverFenster)

	return &a
}

// #### Regele die Steuerung der App und die Umschaltung zwischen den App-Modi ##################

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: zeigt das Hilfefenster an oder blendet es wieder aus. Das Spiel wird solang angehalten.
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

// Umschalter zwischen den App-Zuständen (wird als go-Routine ausgelagert)
//
// Zweck: die Umschaltung zwischen Quiz und Spiel gemäß der Regeln.
//
//	Vor.: keine
//	Eff.:
//	Falls Spielzeit abgelaufen war: Spiel wird beendet.
//	Falls Anzahl Fouls >> Anzahl Treffer: Quiz ist aktiviert.
//	Falls Anzahl Fouls < Anzahl Treffer oder 0: Spiel ist aktiviert
func (a *bpapp) quizUmschalterFunktion() func() {
	return func() {
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
}

// Die Maussteuerung (wird als go-Routine ausgelagert)
//
//	Vor.: keine
//	Eff.: Gibt einige der möglichen Mausaktionen an passende Widgets weiter.
//	Sonst: keiner
func (a *bpapp) mausSteuerFunktion() func(uint8, int8, uint16, uint16) {
	return func(taste uint8, status int8, mausX, mausY uint16) {
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
}

// Die Tastatursteuerung.
//
//	Vor: keine
//	Eff.: die zur Taste passende Spiel-Aktion ist ausgeführt
func (a *bpapp) tastenSteuerFunktion() func(uint16, uint8, uint16) bool {
	return func(taste uint16, gedrückt uint8, tiefe uint16) bool {
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
				return true
				// ######  Testzwecke ####################################
			case 's': // Zeitlupe
				a.billard.ZeitlupeAnAus()
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
		return false
	}
}

// Zweck: startet die Laufzeit-Elemente der App
//
//	Vor.: die App läuft nicht
//	Eff.: Die App wurde gestartet und ein gfx-Fenster geöffnet.
func (a *bpapp) Run() {
	if a.laeuft {
		return
	}
	println("Willkommen bei BrainPool")
	println("Öffne Gfx-Fenster")
	b, h := a.hintergrund.GibGroesse()
	gfx.Fenster(b, h) //Fenstergröße
	gfx.Fenstertitel("BrainPool - Das MiniBillard für Schlaue.")

	a.billard.Starte()
	a.spielFenster.Einblenden()
	a.quizFenster.Ausblenden()
	a.hilfeFenster.Ausblenden()
	a.gameOverFenster.Ausblenden()
	a.mausSteuerung = views_controls.NewMausRoutine(a.mausSteuerFunktion())
	a.mausSteuerung.StarteRate(20) // go-Routine
	// der eigentliche Event-Loop der App läuft nebenher
	a.umschalter = hilf.NewRoutine("Umschalter", a.quizUmschalterFunktion())
	a.umschalter.StarteRate(20) // go-Routine
	a.geraeusche.StarteLoop()   // go-Routine
	a.renderer.Starte()         // go-Routine
	a.laeuft = true

	// ####### der Tastatur-Loop darf dafür hier existieren ####################
	var aktion func(uint16, uint8, uint16) bool = a.tastenSteuerFunktion()
	for {
		taste, gedrückt, tiefe := gfx.TastaturLesen1() // blockiert, bis Taste gedrückt
		if aktion(taste, gedrückt, tiefe) {
			return
		}
	}
}

// Zweck: stoppt die Laufzeit-Elemente der App
//
//	Vor.: die App läuft
//	Eff.: Die App wurde beendet und das gfx-Fenster geschlossen.
func (a *bpapp) Quit() {
	if !a.laeuft {
		return
	}
	a.geraeusche.Stoppe()
	a.musik.Stoppe()
	a.renderer.UeberblendeText("Bye!", views_controls.Fanzeige(), views_controls.Ftext(), 30)
	go a.mausSteuerung.Stoppe() // go-Routine
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
