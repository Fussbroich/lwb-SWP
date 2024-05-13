package apps

import (
	vc "../views_controls"
)

// Die Modelle und Views einer Beispiel-App
type app struct {
	quitter func()
	// Größe des gesamten App-Fensters
	breite uint16
	hoehe  uint16
	titel  string
	// Modelle
	// Views
	spielFenster vc.Widget
	hilfeFenster vc.Widget
	buttonLeiste []vc.Widget
	widgets      []vc.Widget
	//Darstellungsvariable
	darkmode bool
}

// Zweck: Konstruktor baut eine App zusammen
//
//	Eff.:  Ein App-Objekt steht zum Starten bereit.
//	Hinweis: Man startet eine App mit RunApp(App).
func NewBeispielApp(b uint16) *app {
	if b > 1920 {
		b = 1920 // größtmögliches gfx-Fenster ist 1920 Pixel breit
	}
	if b < 480 {
		b = 480 // kleinere Bildschirme sind zum Spielen ungeeignet
	}

	var g uint16 = b / 32 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 16:11
	a := app{titel: "Eine App.",
		breite: 32 * g, hoehe: 22 * g}

	// ######## Modelle und Views zusammenstellen #################################
	// Views erzeugen
	var hintergrund vc.Widget = vc.NewFenster()
	a.spielFenster = vc.NewFenster()
	// Buttonleiste
	a.buttonLeiste = []vc.Widget{
		vc.NewButton("(h)ilfe", a.hilfeAnAus),
		vc.NewButton("(d)unkel/hell", a.darkmodeAnAus),
		vc.NewButton("(s)chließen", a.quit)}

	// Hilfe
	var hilfetext string = "Hilfe\n\nHier steht der Hilfetext"
	a.hilfeFenster = vc.NewTextBox(hilfetext, vc.Regular, int(a.breite/56))

	//setze Layout
	hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = 4 * g, 6 * g, 28 * g, 18 * g
	var g3 uint16 = g + g/3
	// Spielfeld
	a.spielFenster.SetzeKoordinaten(xs, ys, xe, ye)
	// die übrigen Fenster stehen genau vor dem Spielfeld
	a.hilfeFenster.SetzeKoordinaten(xs-g3+2, ys-g3+2, xe+g3-2, ye+g3-2)
	// Buttons unterhalb des Spielfelds gleichmäßig verteilt
	zb := (a.breite - 2*g) / uint16(len(a.buttonLeiste))
	for i, b := range a.buttonLeiste {
		b.SetzeKoordinaten(g+uint16(i)*zb+zb/8, ye+5*g/2, g+uint16(i+1)*zb-zb/8, ye+13*g/4)
	}

	//setzeFarben
	hintergrund.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	a.spielFenster.SetzeFarben(vc.Fbillardtuch, vc.Fdiamanten)
	a.hilfeFenster.SetzeFarben(vc.Fquiz, vc.Ftext)
	for _, b := range a.buttonLeiste {
		b.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	}

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.widgets = append(a.widgets, hintergrund)
	a.widgets = append(a.widgets, a.spielFenster, a.hilfeFenster)
	a.widgets = append(a.widgets, a.buttonLeiste...)

	// Setze App-Zustand
	a.hilfeFenster.Ausblenden()
	a.spielFenster.Einblenden()
	return &a
}

func (a *app) SetzeQuit(f func()) { a.quitter = f }

func (a *app) GibGroesse() (uint16, uint16) { return a.breite, a.hoehe }

func (a *app) GibTitel() string { return a.titel }

// Die Darstell-Funktion - wird vom Zeichen-Loop bei jedem Tick einmal aufgerufen.
//
//	Vor.: Gfx-Fenster ist offen.
func (a *app) Zeichne() {
	for _, f := range a.widgets {
		if f.IstAktiv() {
			f.Update()
			f.Zeichne()
		}
	}
}

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
func (a *app) Update() {
	// hier kommt die Logik rein
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: zeigt das Hilfefenster an oder blendet es wieder aus. Das Spielmodell wird solang angehalten.
func (a *app) hilfeAnAus() {
	if a.hilfeFenster.IstAktiv() {
		a.hilfeFenster.Ausblenden()
		a.spielFenster.Einblenden()
	} else {
		a.spielFenster.Ausblenden()
		a.hilfeFenster.Einblenden()
	}
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: keine
//	Eff.: die GUI wird zwischen hell und dunkel umgeschaltet
func (a *app) darkmodeAnAus() {
	if !a.darkmode {
		vc.SetzeDarkFarbSchema()
	} else {
		vc.SetzeStandardFarbSchema()
	}
	for _, w := range a.widgets {
		w.LadeFarben()
	}
	a.darkmode = !a.darkmode
}

// Aktion für einen klickbaren Button.
//
//	Vor.: die App läuft
//	Eff.: Die App und die Loops wurden beendet
func (a *app) quit() {
	//tue etwas

	// IMMER ZULETZT: rufe Quitter der AppRunner-Funktion
	if a.quitter != nil {
		a.quitter()
	}
}

// Die Maussteuerung der App (kann als go-Routine in einem Loop laufen).
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: Gibt einige der möglichen Mausaktionen an passende vc.Widgets weiter.
//	Sonst: keiner
func (a *app) MausEingabe(taste uint8, status int8, mausX, mausY uint16) {
	if taste == 1 && status == -1 { // es wurde links geklickt
		// wurde ein Button angeklickt?
		for _, b := range a.buttonLeiste {
			if b.IstAktiv() && b.ImFenster(mausX, mausY) {
				b.MausklickBei(mausX, mausY)
				return
			}
		}
		// sonst gib den Klick ans Spiel
		// (zum Spielen kann man auch außerhalb des Spielfensters klicken)
		if a.spielFenster.IstAktiv() {
			// kann auch außerhalb des Tuchs klicken
			a.spielFenster.MausklickBei(mausX, mausY)
		}
	}
}

// Die Tastatursteuerung der App.
//
//	Vor: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: die zur Taste passende Spiel-Aktion ist ausgeführt.
func (a *app) TastaturEingabe(taste uint16, gedrückt uint8, _ uint16) {
	if gedrückt == 1 {
		switch taste {
		case 'h': // Hilfe an-aus
			a.hilfeAnAus()
		case 'd': // Dunkle Umgebung
			a.darkmodeAnAus()
		case 's', 'q':
			a.quit()
		}
	}
}
