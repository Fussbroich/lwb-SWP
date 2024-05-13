package apps

import (
	"time"

	vc "../views_controls"
)

// Das zusätzliche Modell einer Uhr
type appUhr struct {
	uhrzeit time.Time
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
	var g uint16 = b / 16 // Rastermass für dieses App-Design

	// Das Seitenverhältnis des App-Fensters ist B:H = 2:1
	a := appUhr{app: app{titel: "Die Uhrzeit.", breite: 16 * g, hoehe: 8 * g}}

	// Views erzeugen
	var hintergrund vc.Widget = vc.NewFenster()
	var uhrzeitAnzeige vc.Widget = vc.NewDigitalUhrzeitAnzeiger(&a.uhrzeit)
	// Buttonleiste
	a.buttonLeiste = []vc.Widget{
		vc.NewButton("(d)unkel/hell", a.darkmodeAnAus),
		vc.NewButton("(s)chließen", a.quit)}

	//setze Layout
	hintergrund.SetzeKoordinaten(0, 0, a.breite, a.hoehe)
	var xs, ys, xe, ye uint16 = 2 * g, g, 14 * g, 5 * g
	// Anzeige
	uhrzeitAnzeige.SetzeKoordinaten(xs, ys, xe, ye)
	// Buttons unterhalb der Anzeige gleichmäßig verteilt
	zb := a.breite / 5
	hButton := 2 * g / 3
	a.buttonLeiste[0].SetzeKoordinaten(zb, ye+3*hButton/2, 2*zb, ye+5*hButton/2)
	a.buttonLeiste[1].SetzeKoordinaten(3*zb, ye+3*hButton/2, 4*zb, ye+5*hButton/2)

	//setzeFarben
	hintergrund.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	uhrzeitAnzeige.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	for _, b := range a.buttonLeiste {
		b.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	}

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.widgets = append(a.widgets, hintergrund, uhrzeitAnzeige)
	a.widgets = append(a.widgets, a.buttonLeiste...)

	return &a
}

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
func (a *appUhr) Update() {
	a.uhrzeit = time.Now()
}

// Die Tastatursteuerung der App.
//
//	Vor: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: die zur Taste passende Spiel-Aktion ist ausgeführt.
func (a *appUhr) TastaturEingabe(taste uint16, gedrückt uint8, _ uint16) {
	if gedrückt == 1 {
		switch taste {
		case 'd': // Dunkle Umgebung
			a.darkmodeAnAus()
		case 'l':
			a.layoutAnAus()
		case 's', 'q':
			a.quit()
		}
	}
}
