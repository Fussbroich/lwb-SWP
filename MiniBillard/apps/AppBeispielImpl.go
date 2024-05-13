package apps

import (
	"time"

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
	uhrzeit time.Time
	// Views
	uhrzeitAnzeige vc.Widget
	hilfeFenster   vc.Widget
	buttonLeiste   []vc.Widget
	widgets        []vc.Widget
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
	var g uint16 = b / 8 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 2:1
	a := app{titel: "Eine Uhr.", breite: 8 * g, hoehe: 4 * g}

	// Views erzeugen
	var hintergrund vc.Widget = vc.NewFenster()
	a.uhrzeitAnzeige = vc.NewDigitalUhrzeitAnzeiger(&a.uhrzeit)
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
	var xs, ys, xe, ye uint16 = g, g, 7 * g, 2 * g
	// Anzeige
	a.uhrzeitAnzeige.SetzeKoordinaten(xs, ys, xe, ye)
	// die Hilfe steht genau davor
	a.hilfeFenster.SetzeKoordinaten(xs, ys, xe, ye)
	// Buttons unterhalb des Spielfelds gleichmäßig verteilt
	zb := (a.breite - 2*g) / uint16(len(a.buttonLeiste))
	hButton := g / 4
	for i, b := range a.buttonLeiste {
		b.SetzeKoordinaten(
			g+uint16(i)*zb+zb/8, ye+g-hButton/2,
			g+uint16(i+1)*zb-zb/8, ye+g+hButton/2)
	}

	//setzeFarben
	hintergrund.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	a.uhrzeitAnzeige.SetzeFarben(vc.Fanzeige, vc.Ftext)
	a.hilfeFenster.SetzeFarben(vc.Fanzeige, vc.Ftext)
	for _, b := range a.buttonLeiste {
		b.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	}

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.widgets = append(a.widgets, hintergrund, a.uhrzeitAnzeige, a.hilfeFenster)
	a.widgets = append(a.widgets, a.buttonLeiste...)

	// Setze App-Zustand
	a.hilfeFenster.Ausblenden()
	a.uhrzeitAnzeige.Einblenden()
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
			f.Zeichne()
		}
	}
}

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
func (a *app) Update() {
	a.uhrzeit = time.Now()
	for _, f := range a.widgets {
		if f.IstAktiv() {
			f.Update()
		}
	}
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: zeigt das Hilfefenster an oder blendet es wieder aus. Das Spielmodell wird solang angehalten.
func (a *app) hilfeAnAus() {
	if a.hilfeFenster.IstAktiv() {
		a.hilfeFenster.Ausblenden()
		a.uhrzeitAnzeige.Einblenden()
	} else {
		a.uhrzeitAnzeige.Ausblenden()
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
