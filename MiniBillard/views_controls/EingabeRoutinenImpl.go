package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type bpEingabeRoutine struct {
	hilf.Routine
}

func abfangenFensterZu() {
	if r := recover(); r != nil {
		if fmt.Sprint(r) == "Das gfx-Fenster ist nicht offen!" {
			println("Abgefangen: Gfx-Fenster ist schon zu.")
			return
		}
		panic(r)
	}
}

func NewMausRoutine(f func(t uint8, s int8, x uint16, y uint16)) *bpEingabeRoutine {
	ctl := bpEingabeRoutine{}
	rfunc := func() {
		if !gfx.FensterOffen() {
			return
		}
		//defer abfangenFensterZu() // TODO: testen und einbauen
		taste, status, mausX, mausY := gfx.MausLesen1() // blockiert, bis Maus bedient
		f(taste, status, mausX, mausY)
	}
	ctl.Routine = hilf.NewRoutine("Maussteuerung", rfunc)
	return &ctl
}

func NewTastenRoutine(f func(uint16, uint8, uint16)) *bpEingabeRoutine {
	ctl := bpEingabeRoutine{}
	rfunc :=
		func() {
			if !gfx.FensterOffen() {
				return
			}
			//defer abfangenFensterZu() // TODO: testen und einbauen
			taste, gedrückt, tiefe := gfx.TastaturLesen1() // blockiert, bis Taste gedrückt
			f(taste, gedrückt, tiefe)
		}
	ctl.Routine = hilf.NewRoutine("Maussteuerung", rfunc)
	return &ctl
}
