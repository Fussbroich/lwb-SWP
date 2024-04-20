package views_controls

import (
	"gfx"

	"../hilf"
)

type bpMausRoutine struct {
	steuerFunktion func(uint8, int8, uint16, uint16) bool
	steuerRoutine  hilf.Routine
}

func NewMausRoutine(f func(t uint8, s int8, x uint16, y uint16) (quitScan bool)) *bpMausRoutine {
	return &bpMausRoutine{steuerFunktion: f}
}

func (ctl *bpMausRoutine) mausLesenUndAuswerten() {
	taste, status, mausX, mausY := gfx.MausLesen1()
	// quitScan wird ignoriert, wenn die Funktion in einer Routine läuft
	ctl.steuerFunktion(taste, status, mausX, mausY)
}

func (ctl *bpMausRoutine) StarteRate(rate uint64) {
	if ctl.steuerRoutine == nil {
		ctl.steuerRoutine = hilf.NewRoutine("Maussteuerung", ctl.mausLesenUndAuswerten)
	}
	ctl.steuerRoutine.StarteRate(rate) // gewünschte Abtastrate je Sekunde
}

func (ctl *bpMausRoutine) Starte() {
	if ctl.steuerRoutine == nil {
		ctl.steuerRoutine = hilf.NewRoutine("Maussteuerung", ctl.mausLesenUndAuswerten)
	}
	ctl.steuerRoutine.Starte()
}

func (ctl *bpMausRoutine) Stoppe() {
	if ctl.steuerRoutine != nil {
		ctl.steuerRoutine.Stoppe()
	}
}
