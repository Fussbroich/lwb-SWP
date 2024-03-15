package views

import (
	"fmt"
	"gfx"

	"../hilf"
)

type Zeichner interface {
	Starte()
	Stoppe()
	Überblende(Fenster)
	ÜberblendeText(string, Farbe)
	ÜberblendeAus()
}

type zeichner struct {
	breite, höhe uint16
	fenster      []Fenster
	overlay      Fenster
	updater      hilf.Prozess
	updaterLäuft bool
	rate         uint64
}

func NewZeichner(fenster ...Fenster) *zeichner {
	var bMax, hMax uint16
	for _, f := range fenster {
		b, h := f.GibGröße()
		bMax = max(bMax, b)
		hMax = max(hMax, h)
	}
	return &zeichner{fenster: fenster, breite: bMax, höhe: hMax, rate: 80}
}

// ######## die Start- und Stop-Methode ###########################################################

func (r *zeichner) Starte() {
	if r.updaterLäuft {
		return
	}

	r.updater = hilf.NewProzess("Renderer",
		func() {
			gfx.UpdateAus()
			gfx.Cls()
			for _, f := range r.fenster {
				f.Zeichne()
			}
			// zeige die frame rate
			info := fmt.Sprintf("%04d fps", r.updater.GibRate()/10*10)
			NewInfoText(r.breite/50, r.höhe/50, r.breite/8, r.höhe/10, info, F(249, 73, 68)).Zeichne()
			if r.overlay != nil {
				r.overlay.Zeichne()
			}
			gfx.UpdateAn()
		})
	r.updaterLäuft = true
	r.updater.StarteRate(r.rate)
}

func (r *zeichner) Stoppe() {
	if !r.updaterLäuft {
		return
	}
	r.updater.Stoppe()
	r.updaterLäuft = false
}

// ######## die übrigen Methoden ####################################################

func (r *zeichner) Überblende(f Fenster) {
	r.overlay = f
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}

func (r *zeichner) ÜberblendeText(t string, c Farbe) {
	r.overlay = NewTextOverlay(0, 0, r.breite, r.höhe, t, 180, Weiß(), c)
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}

func (r *zeichner) ÜberblendeAus() {
	r.overlay = nil
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}
