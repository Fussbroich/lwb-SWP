package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type fzeichner struct {
	hintergrund   Widget
	widgets       []Widget
	overlay       Widget
	updater       hilf.Routine
	updaterLaeuft bool
	layoutModus   bool
	darkmode      bool
	rate          uint64
}

func NewFensterZeichner() *fzeichner {
	return &fzeichner{rate: 80, hintergrund: NewFenster()}
}

func (r *fzeichner) SetzeFensterHintergrund(w Widget) {
	r.hintergrund = w
}

func (r *fzeichner) SetzeWidgets(w ...Widget) {
	r.widgets = w
}

// ######## die Start- und Stop-Methode ###########################################################

func (r *fzeichner) Starte() {
	if r.updaterLaeuft {
		return
	}
	r.updater = hilf.NewRoutine("Zeichner",
		func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.UpdateAus()
			gfx.Cls()
			r.hintergrund.Zeichne()
			for _, f := range r.widgets {
				if f.IstAktiv() {
					f.Zeichne()
					if r.layoutModus {
						f.ZeichneLayout()
					}
				}
			}
			b, h := r.hintergrund.GibGroesse()
			// zeige die frame rate
			fps := NewInfoText(fmt.Sprintf("%04d fps", r.updater.GibRate()/10*10))
			fps.SetzeKoordinaten(0, 0, b/2, h/30)
			fps.SetzeFarben(Fanzeige(), Finfos())
			fps.Zeichne()
			// zeige das copyright an
			copy := NewInfoText("(c)2024 Bettina Chang, Thomas Schrader")
			copy.SetzeKoordinaten(2*b/3, 0, b, h/30)
			copy.SetzeFarben(Fanzeige(), Finfos())
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
	r.updaterLaeuft = false
}

// ######## die übrigen Methoden ####################################################
func (r *fzeichner) LayoutAnAus() { r.layoutModus = !r.layoutModus }

func (r *fzeichner) DarkmodeAnAus() {
	if !r.darkmode {
		DarkFarbSchema()
	} else {
		StandardFarbSchema()
	}
	r.hintergrund.LadeFarben()
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

func (r *fzeichner) UeberblendeText(t string, hg, vg string, tr uint8) {
	b, h := r.hintergrund.GibGroesse()
	r.overlay = NewTextOverlay(t)
	r.overlay.SetzeKoordinaten(0, 0, b, h)
	r.overlay.SetzeFarben(hg, vg)
	r.overlay.SetzeTransparenz(tr)
}

func (r *fzeichner) UeberblendeAus() {
	r.overlay = nil
}
