package apps

import (
	"fmt"
	"gfx"

	"../hilf"
	vc "../views_controls"
)

// Eine spezialisierte Routine, die die Darstellungs-Methode einer App aufruft
// und damit die Bildwiederholung steuert.
// Im "unmittelbaren Modus" wird die App in einem regelmäßigen Takt in ein
// einziges Fenster gezeichnet.
//
//	Vor.: Das Grafikpaket gfx muss im GOPATH installiert sein.
//
//	NewZeichenRoutine(App) erzeugt ein Objekt.
type bpZeichenRoutine struct {
	hilf.Routine
}

var (
	renderer *bpZeichenRoutine
	fpsInfo  vc.Widget = vc.NewInfoText(func() string { return fmt.Sprintf("%04d fps", renderer.GibRate()/10*10) })
	cRight   vc.Widget = vc.NewInfoText(func() string { return "(c)2024 Bettina Chang, Thomas Schrader" })
)

func NewZeichenRoutine(a App) *bpZeichenRoutine {
	if renderer != nil {
		return renderer
	}
	b, h := a.GibGroesse()
	fpsInfo.SetzeKoordinaten(0, 0, b/2, h/30)
	fpsInfo.SetzeFarben(vc.Fanzeige, vc.Finfos)
	cRight.SetzeKoordinaten(2*b/3, 0, b, h/30)
	cRight.SetzeFarben(vc.Fanzeige, vc.Finfos)

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
	r.Routine.Stoppe()
	if gfx.FensterOffen() {
		println("Schließe Gfx-Fenster")
		gfx.FensterAus()
	}
}
