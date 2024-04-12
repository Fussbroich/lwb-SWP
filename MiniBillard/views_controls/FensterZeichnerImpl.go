package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type fzeichner struct {
	breite, hoehe uint16
	widgets       []Widget
	overlay       Widget
	updater       hilf.Prozess
	updaterLaeuft bool
	rate          uint64
}

func NewFensterZeichner(w ...Widget) *fzeichner {
	bMax, hMax := w[0].GibGroesse()
	return &fzeichner{widgets: w, breite: bMax, hoehe: hMax, rate: 80}
}

// ######## die Start- und Stop-Methode ###########################################################

func (r *fzeichner) Starte() {
	if r.updaterLaeuft {
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
			fps := NewInfoText(fmt.Sprintf("%04d fps", r.updater.GibRate()/10*10))
			fps.SetzeKoordinaten(0, 0, r.breite/2, r.hoehe/30)
			fps.SetzeFarben(Fanzeige(), Finfos())
			fps.Zeichne()
			// zeige das copyright an
			copy := NewInfoText("(c)2024 Bettina Chang, Thomas Schrader")
			copy.SetzeKoordinaten(2*r.breite/3, 0, r.breite, r.hoehe/30)
			copy.SetzeFarben(Fanzeige(), Finfos())
			copy.Zeichne()
			if r.overlay != nil {
				r.overlay.Zeichne()
			}
			gfx.UpdateAn()
		})
	r.updaterLaeuft = true
	r.updater.StarteRate(r.rate)
	//r.updater.Starte()
}

func (r *fzeichner) Stoppe() {
	if !r.updaterLaeuft {
		return
	}
	r.updater.Stoppe()
	r.updaterLaeuft = false
}

// ######## die übrigen Methoden ####################################################

func (r *fzeichner) ZeigeLayout() {
	for _, f := range r.widgets {
		f.ZeichneLayout()
	}
	if r.overlay != nil {
		r.overlay.Zeichne()
	}
	info := NewInfoText("Layout-Ansicht")
	info.SetzeKoordinaten(r.breite/2, 0, r.breite/2, r.hoehe/10)
	info.SetzeFarben(Fanzeige(), Finfos())
	info.Zeichne()
}

func (r *fzeichner) Ueberblende(f Widget) {
	r.overlay = f
}

func (r *fzeichner) UeberblendeText(t string, hg, vg Farbe, tr uint8) {
	r.overlay = NewTextOverlay(t)
	r.overlay.SetzeKoordinaten(0, 0, r.breite, r.hoehe)
	r.overlay.SetzeFarben(hg, vg)
	r.overlay.SetzeTransparenz(tr)
}

func (r *fzeichner) UeberblendeAus() {
	r.overlay = nil
}
