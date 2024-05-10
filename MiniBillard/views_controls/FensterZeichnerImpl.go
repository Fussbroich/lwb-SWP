package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type fzeichner struct {
	breite        uint16
	hoehe         uint16
	fenstertitel  string
	widgets       []Widget
	overlay       Widget
	updater       hilf.Routine
	updaterLaeuft bool
	layoutModus   bool
	darkmode      bool
	schlicht      bool
	rate          uint64
}

func NewFensterZeichner() *fzeichner {
	hintergrund := NewFenster()
	hintergrund.SetzeKoordinaten(0, 0, 640, 480)
	return &fzeichner{rate: 80}
}

func (r *fzeichner) SetzeFensterGroesse(b, h uint16) { r.breite, r.hoehe = b, h }

func (r *fzeichner) SetzeFensterTitel(t string) { r.fenstertitel = t }

func (r *fzeichner) AddWidgets(w ...Widget) {
	r.widgets = append(r.widgets, w...)
}

// ######## die Start- und Stop-Methode ###########################################################

func (r *fzeichner) Starte() {
	if r.updaterLaeuft {
		return
	}
	if r.schlicht {
		for _, f := range r.widgets {
			f.SetzeSchlicht()
		}
	}
	r.updater = hilf.NewRoutine("Zeichner",
		func() {
			if !gfx.FensterOffen() {
				println("Öffne Gfx-Fenster")
				gfx.Fenster(r.breite, r.hoehe) //Fenster öffnen
				gfx.Fenstertitel(r.fenstertitel)
			}
			gfx.UpdateAus()
			gfx.Cls()
			for _, f := range r.widgets {
				f.Update()
			}
			for _, f := range r.widgets {
				f.Zeichne()
			}
			if r.layoutModus {
				for _, f := range r.widgets {
					f.ZeichneLayout()
				}
			}
			// zeige die frame rate
			fps := NewInfoText(fmt.Sprintf("%04d fps", r.updater.GibRate()/10*10))
			fps.SetzeKoordinaten(0, 0, r.breite/2, r.hoehe/30)
			fps.SetzeFarben(Fanzeige, Finfos)
			fps.Zeichne()
			// zeige das copyright an
			copy := NewInfoText("(c)2024 Bettina Chang, Thomas Schrader")
			copy.SetzeKoordinaten(2*r.breite/3, 0, r.breite, r.hoehe/30)
			copy.SetzeFarben(Fanzeige, Finfos)
			copy.Zeichne()
			if r.overlay != nil {
				r.overlay.Zeichne()
			}
			gfx.UpdateAn()
		})
	r.updaterLaeuft = true
	r.updater.StarteRate(r.rate)
}

func (r *fzeichner) Stoppe() {
	if !r.updaterLaeuft {
		return
	}
	r.updater.Stoppe()
	if gfx.FensterOffen() {
		println("Schließe Gfx-Fenster")
		gfx.FensterAus()
	}
	r.updaterLaeuft = false
}

// ######## die übrigen Methoden ####################################################
func (r *fzeichner) LayoutAnAus() { r.layoutModus = !r.layoutModus }

func (r *fzeichner) ZeichneSchlicht() { r.schlicht = true }

func (r *fzeichner) DarkmodeAnAus() {
	if !r.darkmode {
		SetzeDarkFarbSchema()
	} else {
		SetzeStandardFarbSchema()
	}
	if r.overlay != nil {
		r.overlay.LadeFarben()
	}
	for _, w := range r.widgets {
		w.LadeFarben()
	}
	r.darkmode = !r.darkmode
}

func (r *fzeichner) Ueberblende(f Widget) {
	r.overlay = f
}

func (r *fzeichner) UeberblendeText(t string, hg, vg FarbID, sg int) {
	r.overlay = NewTextOverlay(t, sg)
	r.overlay.SetzeKoordinaten(0, 0, r.breite, r.hoehe)
	r.overlay.SetzeFarben(hg, vg)
	r.overlay.SetzeTransparenz(20)
}

func (r *fzeichner) UeberblendeAus() {
	r.overlay = nil
}
