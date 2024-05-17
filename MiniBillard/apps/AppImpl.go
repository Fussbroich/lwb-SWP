package apps

import (
	vc "../views_controls"
)

// Die Variablen und Views einer App
type app struct {
	quitter func()
	// Größe des gesamten App-Fensters
	breite uint16
	hoehe  uint16
	titel  string
	// Views
	buttonLeiste []vc.Widget
	widgets      []vc.Widget
	overlay      vc.Widget
	testInfo     vc.Widget
	//Darstellungsvariable
	darkmode  bool
	testModus bool
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
	if a.testModus {
		for _, f := range a.widgets {
			if f.IstAktiv() {
				f.ZeichneLayout()
			}
		}
		a.testInfo.Zeichne()
	}
	if a.overlay != nil {
		a.overlay.Zeichne()
	}

}

// Die Update-Funktion - wird vom Spiel-Loop bei jedem Tick einmal aufgerufen
func (a *app) Update() {
	// tue etwas
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

// Testzwecke: zeige Widget-Layout und Infos
func (a *app) testAnAus() {
	if a.testInfo == nil {
		b, h := a.GibGroesse()
		a.testInfo = vc.NewInfoText(func() string { return "Test" })
		a.testInfo.SetzeKoordinaten(2*b/5, 2*h/5, 3*b/5, 3*h/5)
		a.testInfo.SetzeFarben(vc.Fhintergrund, vc.Frot)
	}
	a.testInfo.LadeFarben()
	a.testModus = !a.testModus
}

// Aktion für einen klickbaren Button.
//
//	Vor.: die App läuft
//	Eff.: Die App und die Loops wurden beendet
func (a *app) quit() {
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
		case 'd': // Dunkle Umgebung
			a.darkmodeAnAus()
		case 's', 'q':
			a.quit()
			// ################ Testzwecke #####################
		case 't': // schalte in den Testmodus
			a.testAnAus()
		}
	}
}
