package views

import (
	"fmt"
	"gfx"
	"math"

	"../hilf"
)

type FensterZeichner interface {
	Starte()
	Stoppe()
	Überblende(Fenster)
	ÜberblendeText(string, Farbe)
	ÜberblendeAus()
}

type fzeichner struct {
	breite, höhe uint16
	fenster      []Fenster
	overlay      Fenster
	updater      hilf.Prozess
	updaterLäuft bool
	rate         uint64
}

func NewFensterZeichner(fenster ...Fenster) *fzeichner {
	var bMax, hMax uint16
	for _, f := range fenster {
		b, h := f.GibGröße()
		bMax = uint16(math.Max(float64(bMax), float64(b)))
		hMax = uint16(math.Max(float64(hMax), float64(h)))
	}
	return &fzeichner{fenster: fenster, breite: bMax, höhe: hMax, rate: 80}
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

func (r *fzeichner) Stoppe() {
	if !r.updaterLäuft {
		return
	}
	r.updater.Stoppe()
	r.updaterLäuft = false
}

// ######## die übrigen Methoden ####################################################

func (r *fzeichner) Überblende(f Fenster) {
	r.overlay = f
	if r.updaterLäuft {
		r.Stoppe()
		r.Starte()
	}
}

func (r *fzeichner) ÜberblendeText(t string, c Farbe) {
	r.overlay = NewTextOverlay(0, 0, r.breite, r.höhe, t, 180, Weiß(), c)
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
