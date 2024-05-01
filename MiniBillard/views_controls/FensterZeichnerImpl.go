package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type fzeichner struct {
	hintergrund   Widget
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
	return &fzeichner{rate: 80, hintergrund: hintergrund}
}

func (r *fzeichner) SetzeFensterHintergrund(w Widget) { r.hintergrund = w }

func (r *fzeichner) SetzeFensterTitel(t string) { r.fenstertitel = t }

func (r *fzeichner) SetzeWidgets(w ...Widget) { r.widgets = w }

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
	b, h := r.hintergrund.GibGroesse()
	println("Öffne Gfx-Fenster")
	gfx.Fenster(b, h) //Fenster öffnen
	gfx.Fenstertitel(r.fenstertitel)
	r.updater = hilf.NewRoutine("Zeichner",
		func() {
			if !gfx.FensterOffen() {
				return
			}
			gfx.UpdateAus()
			gfx.Cls()
			if r.hintergrund != nil {
				r.hintergrund.Update()
				r.hintergrund.Zeichne()
			}
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
			if r.hintergrund != nil {
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
			}
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
	if r.hintergrund != nil {
		r.hintergrund.LadeFarben()
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

func (r *fzeichner) UeberblendeText(t string, hg, vg string, sg int) {
	b, h := r.hintergrund.GibGroesse()
	r.overlay = NewTextOverlay(t)
	r.overlay.SetzeKoordinaten(0, 0, b, h)
	r.overlay.SetzeFarben(hg, vg)
	r.overlay.SetzeTransparenz(30)
}

func (r *fzeichner) UeberblendeAus() {
	r.overlay = nil
}
