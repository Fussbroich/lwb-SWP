package apps

import (
	"time"

	vc "../views_controls"
)

// Die zusätzlichen Modelle und Views einer Digitaluhr
type appUhr struct {
	// Modelle
	uhrzeit time.Time
	// Views
	uhrzeitAnzeige vc.Widget
	hilfeFenster   vc.Widget
	app
}

// Zweck: Konstruktor baut eine App zusammen
//
//	Eff.:  Ein App-Objekt steht zum Starten bereit.
//	Hinweis: Man startet eine App mit RunApp(App).
func NewBeispielApp(b uint16) *appUhr {
	if b > 1920 {
		b = 1920 // größtmögliches gfx-Fenster ist 1920 Pixel breit
	}
	var g uint16 = b / 8 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 2:1
	a := appUhr{app: app{titel: "Eine Uhr.", breite: 8 * g, hoehe: 4 * g}}

	// Views erzeugen
	var hintergrund vc.Widget = vc.NewFenster()
	a.uhrzeitAnzeige = vc.NewDigitalUhrzeitAnzeiger(&a.uhrzeit)
	// Buttonleiste
	a.buttonLeiste = []vc.Widget{
		vc.NewButton("(h)ilfe", a.hilfeAnAus),
		vc.NewButton("(d)unkel/hell", a.darkmodeAnAus),
		vc.NewButton("(s)chließen", a.quit)}

	// Hilfe
	var hilfetext string = "Hilfe\n\nKlicke die Buttons unten an oder drücke die angegebene Taste."
	a.hilfeFenster = vc.NewTextBox(hilfetext, vc.Regular, int(a.breite/30))

	//setze Layout
	hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = g, g / 2, 7 * g, 5 * g / 2
	// Anzeige
	a.uhrzeitAnzeige.SetzeKoordinaten(xs, ys, xe, ye)
	// die Hilfe steht genau davor
	a.hilfeFenster.SetzeKoordinaten(xs, ys, xe, ye)
	// Buttons unterhalb des Spielfelds gleichmäßig verteilt
	zb := (a.breite - 2*g) / uint16(len(a.buttonLeiste))
	hButton := g / 3
	for i, b := range a.buttonLeiste {
		b.SetzeKoordinaten(
			g+uint16(i)*zb+zb/8, ye+2*g/3-hButton/2,
			g+uint16(i+1)*zb-zb/8, ye+2*g/3+hButton/2)
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

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
func (a *appUhr) Update() {
	a.uhrzeit = time.Now()
}

// Aktion für einen klickbaren Button oder eine Taste
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: zeigt das Hilfefenster an oder blendet es wieder aus.
func (a *appUhr) hilfeAnAus() {
	if a.hilfeFenster.IstAktiv() {
		a.hilfeFenster.Ausblenden()
		a.uhrzeitAnzeige.Einblenden()
	} else {
		a.uhrzeitAnzeige.Ausblenden()
		a.hilfeFenster.Einblenden()
	}
}

// Die Tastatursteuerung der App.
//
//	Vor: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: die zur Taste passende Spiel-Aktion ist ausgeführt.
func (a *appUhr) TastaturEingabe(taste uint16, gedrückt uint8, _ uint16) {
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
