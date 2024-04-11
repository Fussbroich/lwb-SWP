package views_controls

import (
	"gfx"

	"../hilf"
)

func NewMausProzess(f func(t uint8, s int8, x uint16, y uint16)) *bpMausProzess {
	return &bpMausProzess{steuerFunktion: f}
}

type bpMausProzess struct {
	steuerProzess  hilf.Prozess
	steuerFunktion func(uint8, int8, uint16, uint16)
}

func (ctl *bpMausProzess) mausSteuerung() {
	taste, status, mausX, mausY := gfx.MausLesen1()
	ctl.steuerFunktion(taste, status, mausX, mausY)
}

func (ctl *bpMausProzess) StarteRate(rate uint64) {
	if ctl.steuerProzess == nil {
		ctl.steuerProzess = hilf.NewProzess(
			"Maussteuerung",
			ctl.mausSteuerung)
	}
	ctl.steuerProzess.StarteRate(rate) // gewünschte Abtastrate je Sekunde
}

func (ctl *bpMausProzess) Starte() {
	if ctl.steuerProzess == nil {
		ctl.steuerProzess = hilf.NewProzess(
			"Maussteuerung",
			ctl.mausSteuerung)
	}
	ctl.steuerProzess.Starte()
}

func (ctl *bpMausProzess) Stoppe() {
	if ctl.steuerProzess != nil {
		ctl.steuerProzess.Stoppe()
	}
}
