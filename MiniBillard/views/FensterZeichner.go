package views

import (
	"fmt"
	"gfx"

	"../hilf"
)

type FensterZeichner interface {
	Starte()
	Stoppe()
	ZeigeLayout()
	Überblende(Widget)
	ÜberblendeText(string, Farbe, Farbe, uint8)
	ÜberblendeAus()
}

type fzeichner struct {
	breite, höhe uint16
	widgets      []Widget
	overlay      Widget
	updater      hilf.Prozess
	updaterLäuft bool
	rate         uint64
}

func NewFensterZeichner(w ...Widget) *fzeichner {
	bMax, hMax := w[0].GibGröße()
	return &fzeichner{widgets: w, breite: bMax, höhe: hMax, rate: 80}
}

// ######## die Start- und Stop-Methode ###########################################################

func (r *fzeichner) Starte() {
	if r.updaterLäuft {
		return
	}

	r.updater = hilf.NewProzess("Zeichner",
		func() {
			gfx.UpdateAus()
			gfx.Cls()
			for _, f := range r.widgets {
				f.Zeichne()
			}
			// zeige die frame rate
			info := fmt.Sprintf("%04d fps", r.updater.GibRate()/10*10)
			NewInfoText(0, 0, r.breite, r.höhe/30, info, F(240, 255, 255)).Zeichne()
			if r.overlay != nil {
				r.overlay.Zeichne()
			}
			gfx.UpdateAn()
		})
	r.updaterLäuft = true
	r.updater.StarteRate(r.rate)
	//r.updater.Starte()
}

func (r *fzeichner) Stoppe() {
	if !r.updaterLäuft {
		return
	}
	r.updater.Stoppe()
	r.updaterLäuft = false
}

// ######## die übrigen Methoden ####################################################

func (r *fzeichner) ZeigeLayout() {
	for _, f := range r.widgets {
		f.ZeichneLayout()
	}
	if r.overlay != nil {
		r.overlay.Zeichne()
	}
	NewInfoText(r.breite/2, 0, r.breite/2, r.höhe/10, "Layout-Ansicht", F(240, 255, 255)).Zeichne()
}

func (r *fzeichner) Überblende(f Widget) {
	r.overlay = f
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}

func (r *fzeichner) ÜberblendeText(t string, hg, vg Farbe, tr uint8) {
	r.overlay = NewTextOverlay(0, 0, r.breite, r.höhe, t, hg, vg, tr)
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}

func (r *fzeichner) ÜberblendeAus() {
	r.overlay = nil
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}
