package views_controls

import (
	"gfx"

	"../hilf"
)

type bpZeichenRoutine struct {
	hilf.Routine
}

func NewZeichenRoutine(a App) *bpZeichenRoutine {
	b, h := a.GibGroesse()
	routine := hilf.NewRoutine("Zeichner",
		func() {
			if !gfx.FensterOffen() {
				println("Öffne Gfx-Fenster")
				gfx.Fenster(b, h) //Fenster öffnen
				gfx.Fenstertitel(a.GibTitel())
			}
			gfx.UpdateAus()
			gfx.Cls()
			a.Zeichne()
			gfx.UpdateAn()
		})
	return &bpZeichenRoutine{Routine: routine}
}

// ######## die Stop-Methode schließt das Gfx-Fenster ###################################

func (r *bpZeichenRoutine) Stoppe() {
	if !r.Routine.Laeuft() {
		return
	}
	r.Routine.Stoppe()
	if gfx.FensterOffen() {
		println("Schließe Gfx-Fenster")
		gfx.FensterAus()
	}
}
