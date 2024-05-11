package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type bpZeichenRoutine struct {
	hilf.Routine
}

var (
	renderer *bpZeichenRoutine
	fpsInfo  Widget = NewInfoText(func() string { return fmt.Sprintf("%04d fps", renderer.GibRate()/10*10) })
	cRight   Widget = NewInfoText(func() string { return "(c)2024 Bettina Chang, Thomas Schrader" })
)

func NewZeichenRoutine(a App) *bpZeichenRoutine {
	if renderer != nil {
		return renderer
	}
	b, h := a.GibGroesse()
	fpsInfo.SetzeKoordinaten(0, 0, b/2, h/30)
	fpsInfo.SetzeFarben(Fanzeige, Finfos)
	cRight.SetzeKoordinaten(2*b/3, 0, b, h/30)
	cRight.SetzeFarben(Fanzeige, Finfos)

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
			fpsInfo.Zeichne()
			cRight.Zeichne()
			gfx.UpdateAn()
		})

	renderer = &bpZeichenRoutine{Routine: routine}
	return renderer
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
