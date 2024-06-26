package apps

import (
	"time"

	"../klaenge"
	"../modelle"
	vc "../views_controls"
)

// Die Modelle und Views einer Brainpool-App
type bpapp struct {
	// Klänge
	musik      klaenge.Klang
	geraeusche klaenge.Klang
	// Modelle
	billard modelle.MiniBillardSpiel
	quiz    modelle.Quiz
	// Views
	spielFenster    vc.Widget
	quizFenster     vc.Widget
	hilfeFenster    vc.Widget
	gameOverFenster vc.Widget
	app
}

// Zweck: Konstruktor für BrainPool - baut eine App zusammen
//
//	Vor.:  Klänge und Liberation-Fonts sind installiert.
//	Eff.:  Ein App-Objekt steht zum Starten bereit. Die Billard-Simulation läuft bereits.
//	Hinweis: Man startet eine App mit RunApp(App).
func NewBPApp(b uint16) *bpapp {
	if b > 1920 {
		b = 1920 // größtmögliches gfx-Fenster ist 1920 Pixel breit
	}
	if b < 480 {
		b = 480 // kleinere Bildschirme sind zum Spielen ungeeignet
	}

	var g uint16 = b / 32 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 16:11
	a := bpapp{app: app{titel: "BrainPool - Das Mini-Billard für Schlaue.",
		breite: 32 * g, hoehe: 22 * g}}

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

	// Views und Zeichner erzeugen
	var hintergrund vc.Widget = vc.NewFenster()
	var bande, punktezaehler, restzeit vc.Widget = vc.NewFenster(),
		vc.NewMBPunkteAnzeiger(a.billard),
		vc.NewMBRestzeitAnzeiger(a.billard)

	a.spielFenster = vc.NewMBSpieltisch(a.billard)
	a.quizFenster = vc.NewQuizFenster(a.quiz,
		func() { a.billard.ReduziereStrafpunkte(); a.quiz.NaechsteFrage() },
		func() { a.quiz.NaechsteFrage() })

	// Buttons
	buttons := []vc.Widget{
		vc.NewButton("(h)ilfe", a.hilfeAnAus),
		vc.NewButton("(n)eues Spiel", a.neuesSpiel),
		vc.NewButton("(m)usik spielen", a.musikAn),
		vc.NewButton("(d)unkel/hell", a.darkmodeAnAus),
		vc.NewButton("(s)chließen", a.quit)}

	// Hilfe
	var hilfetext string = "Hilfe\n\n" +
		"Im Spielmodus (und nur, wenn alle Kugeln still stehen): " +
		"Maus bewegen ändert die Zielrichtung. Stoß durch klicken mit der linken Maustaste. " +
		"Die Stoßkraft wird durch scrollen der Maus verändert.\n\n" +
		"Du spielst gegen die Zeit. Alle neun Kugeln müssen versenkt werden. " +
		"Es gibt ein Foul, wenn die weiße Kugel reingeht oder wenn bei einem Stoß gar keine Kugel versenkt wird.\n\n" +
		"Im Quizmodus: Klicke die richtigen Antworten an, um Fouls abzuarbeiten.\n\n" +
		"Die übrige Bedienung erfolgt durch Anklicken der Buttons unten " +
		"oder mit der angegebenen Taste auf der Tastatur."

	a.hilfeFenster = vc.NewTextBox(hilfetext, vc.Regular, int(a.breite/56))

	// GAME OVER
	a.gameOverFenster = vc.NewTextBox(" \n  * GAME OVER *", vc.BoldItalic, int(a.breite/12))

	//setze Layout
	hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = 4 * g, 6 * g, 28 * g, 18 * g
	var g3 uint16 = g + g/3
	// oben links ist der Punktezähler
	punktezaehler.SetzeKoordinaten(xs-g3, 1*g, 18*g, 3*g)
	// oben rechts ist der Countdown
	restzeit.SetzeKoordinaten(20*g+g3, g, xe+g3, 3*g)
	// Hintergrund für das Tuch
	bande.SetzeKoordinaten(xs-g3, ys-g3, xe+g3, ye+g3)
	bande.SetzeEckradius(g3)
	// Spielfeld (Tuch)
	a.spielFenster.SetzeKoordinaten(xs, ys, xe, ye)
	// die übrigen Fenster stehen genau vor dem Spielfeld
	a.quizFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.quizFenster.SetzeEckradius(g3 - 2)
	a.hilfeFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.hilfeFenster.SetzeEckradius(g3 - 2)
	a.gameOverFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	a.gameOverFenster.SetzeEckradius(g3 - 2)
	// Buttons unterhalb des Spielfelds gleichmäßig verteilt
	zb := (a.breite - 2*g) / uint16(len(buttons))
	for i, b := range buttons {
		b.SetzeKoordinaten(g+uint16(i)*zb+zb/8, ye+5*g/2, g+uint16(i+1)*zb-zb/8, ye+13*g/4)
		b.SetzeEckradius(g / 3)
	}

	//setzeFarben
	hintergrund.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	a.spielFenster.SetzeFarben(vc.Fbillardtuch, vc.Fdiamanten)
	bande.SetzeFarben(vc.Fbande, vc.Fanzeige)
	punktezaehler.SetzeFarben(vc.Fanzeige, vc.Ftext)
	punktezaehler.SetzeTransparenz(255)
	restzeit.SetzeFarben(vc.Fanzeige, vc.Ftext)
	a.quizFenster.SetzeFarben(vc.Fquiz, vc.Ftext)
	a.hilfeFenster.SetzeFarben(vc.Fquiz, vc.Ftext)
	a.gameOverFenster.SetzeFarben(vc.Fquiz, vc.Ftext)
	for _, b := range buttons {
		b.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	}

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.widgets = append(a.widgets, hintergrund)
	a.widgets = append(a.widgets, bande, a.spielFenster, a.quizFenster, a.gameOverFenster, a.hilfeFenster)
	a.widgets = append(a.widgets, punktezaehler, restzeit)
	a.widgets = append(a.widgets, buttons...)
	a.klickbare = append(a.klickbare, a.spielFenster, a.quizFenster)
	a.klickbare = append(a.klickbare, buttons...)

	// Setze Start-Zustand
	a.quizFenster.Ausblenden()
	a.hilfeFenster.Ausblenden()
	a.gameOverFenster.Ausblenden()
	a.spielFenster.Einblenden()
	// Starte Simulation
	a.billard.Starte()
	a.geraeusche.StarteLoop()
	return &a
}

// Die Darstell-Funktion - wird vom Zeichen-Loop bei jedem Tick einmal aufgerufen.
//
//	Vor.: Gfx-Fenster ist offen.
func (a *bpapp) Zeichne() { a.app.Zeichne() }

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.:
//	Falls Spielzeit abgelaufen war: Spiel wird beendet.
//	Falls Anzahl Fouls >> Anzahl Treffer: Quiz ist aktiviert.
//	Falls Anzahl Fouls < Anzahl Treffer oder 0: Spiel ist aktiviert
//
//	Hinweis: Die Funktion hier bestimmt lediglich die Umschaltung zwischen Quiz und
//	Spiel-Simulation. Die Simulation bringt einen eigenen Loop.
func (a *bpapp) Update() {
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
	a.app.Update()
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: neues Spiel ist gestartet - alle anderen Fenster sind ausgeblendet
func (a *bpapp) neuesSpiel() {
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
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: zeigt das Hilfefenster an oder blendet es wieder aus. Das Spielmodell wird solang angehalten.
func (a *bpapp) hilfeAnAus() {
	if a.hilfeFenster.IstAktiv() {
		a.hilfeFenster.Ausblenden()
		if a.spielFenster.IstAktiv() {
			a.billard.Starte()
		}
	} else {
		a.hilfeFenster.Einblenden()
		if a.spielFenster.IstAktiv() {
			a.billard.Stoppe()
		}
	}
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: Die Musik ist definiert.
//	Eff.: die hinterlegte Musik startet im Loop
func (a *bpapp) musikAn() {
	a.musik.StarteLoop()
}

// Aktion für einen klickbaren Button.
//
//	Vor.: die App läuft
//	Eff.: Die App und die Loops wurden beendet
func (a *bpapp) quit() {
	a.geraeusche.Stoppe()
	a.musik.Stoppe()
	// flicke einen Abschiedsscreen ein
	a.overlay = vc.NewTextOverlay("Bye!", int(a.hoehe/5))
	a.overlay.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	a.overlay.SetzeFarben(vc.Fanzeige, vc.Ftext)
	a.overlay.SetzeTransparenz(20)
	a.billard.Stoppe()
	// IMMER ZULETZT: rufe quit der app
	a.app.quit()
}

// Die Maussteuerung der App (kann als go-Routine in einem Loop laufen).
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: Gibt einige der möglichen Mausaktionen an passende vc.Widgets weiter.
//	Sonst: keiner
func (a *bpapp) MausEreignis(taste uint8, status int8, x, y uint16) {
	if taste == 1 && status == -1 { // es wurde links geklickt
		// wurde etwas angeklickt?
		for _, b := range a.klickbare {
			if b.IstAktiv() && b.ImFenster(x, y) {
				b.MausklickBei(x, y)
				return
			}
		}
		// sonst gib den Klick ans Spiel
		// (zum Spielen kann man auch außerhalb des Spielfensters klicken)
		if a.spielFenster.IstAktiv() {
			a.spielFenster.MausklickBei(x, y)
			return
		}
	} else { // es wurde nicht links geklickt
		// zielen und Kraft aufbauen
		if a.spielFenster.IstAktiv() {
			switch taste {
			case 4: // vorwärts scrollen
				a.spielFenster.MausScrolltHoch()
			case 5: // rückwärts scrollen
				a.spielFenster.MausScrolltRunter()
			default: // bewegen
				a.spielFenster.MausBei(x, y)
			}
		}
		// Sonst: tue gar nichts
	}
}

// Die Tastatursteuerung der App.
//
//	Vor: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: die zur Taste passende Spiel-Aktion ist ausgeführt.
func (a *bpapp) TastaturEreignis(taste uint16, gedrückt uint8, _ uint16) {
	if gedrückt == 1 {
		switch taste {
		case 'h', 'H': // Hilfe an-aus
			a.hilfeAnAus()
		case 'n', 'N': // neues Spiel
			a.neuesSpiel()
		case 'd', 'D': // Dunkle Umgebung
			a.darkmodeAnAus()
		case 'm', 'M': // Musik spielen, wenn man möchte
			a.musikAn() // go-Routine
		case 's', 'S', 'q', 'Q':
			a.quit()
		// ######  Testzwecke ####################################
		case 't', 'T': // Test-Modus
			a.testAnAus()
		case 'e', 'E': // Spiel testen
			if a.testModus {
				a.billard.ErhoeheStrafpunkte()
			}
		case 'r', 'R': // Spiel testen
			if a.testModus {
				a.billard.ReduziereStrafpunkte()
			}
		case '1': // Spiel testen
			if a.testModus {
				a.billard.SetzeRestzeit(5 * time.Second)
				a.billard.SetzeKugeln1BallTest()
			}
		case '3': // Spiel testen
			if a.testModus {
				a.billard.SetzeSpielzeit(90 * time.Second)
				a.billard.SetzeKugeln3Ball()
			}
		case '9': // Spiel testen
			if a.testModus {
				a.billard.SetzeSpielzeit(4 * time.Minute)
				a.billard.SetzeKugeln9Ball()
			}
		}
	}
}
