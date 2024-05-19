package apps

import (
	"time"

	vc "../views_controls"
)

// Das Modell einer Uhr
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
	// Buttons
	buttons := []vc.Widget{
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
	buttons[0].SetzeKoordinaten(zb, ye+3*hButton/2, 2*zb, ye+5*hButton/2)
	buttons[1].SetzeKoordinaten(3*zb, ye+3*hButton/2, 4*zb, ye+5*hButton/2)

	//setzeFarben
	hintergrund.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	uhrzeitAnzeige.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	for _, b := range buttons {
		b.SetzeFarben(vc.Fhintergrund, vc.Ftext)
	}

	// Reihenfolge der Views ist teilweise wichtig (obere decken untere ab)
	a.widgets = append(a.widgets, hintergrund, uhrzeitAnzeige)
	a.widgets = append(a.widgets, buttons...)
	a.klickbare = append(a.klickbare, buttons...)

	return &a
}

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
func (a *appUhr) Update() {
	a.uhrzeit = time.Now()
	a.app.Update()
}

// Die Darstell-Funktion - wird vom Zeichen-Loop bei jedem Tick einmal aufgerufen.
//
//	Vor.: Gfx-Fenster ist offen.
func (a *appUhr) Zeichne() { a.app.Zeichne() }

// Die Tastatursteuerung der App (wird standardmäßig lokal geloopt).
//
//	Vor: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: die zur Taste passende Aktion ist ausgeführt.
func (a *appUhr) TastaturEreignis(taste uint16, gedrückt uint8, tiefe uint16) {
	a.app.TastaturEreignis(taste, gedrückt, tiefe)
}

// Die Maussteuerung der App (wird als go-Routine in einem Loop gestartet).
//
//	Vor.: Alle Modelle und die Fenster der App sind definiert.
//	Eff.: Gibt einige der möglichen Mausaktionen an passende Widgets weiter.
//	Sonst: keiner
func (a *appUhr) MausEreignis(taste uint8, status int8, x, y uint16) {
	a.app.MausEreignis(taste, status, x, y)
}
