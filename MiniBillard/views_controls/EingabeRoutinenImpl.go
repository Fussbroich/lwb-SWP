package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

// Eine spezialisierte Routine für die
// Abfrage der blockierenden Eingabe-Funktionen
// von gfx in einer Endlosschleife.
// Siehe die gfx-Dokumentation für Infos zur Benutzung der
// Eingabe-Funktionen.
//
// Konstruktoren entsprechend dem gewünschten Eingabekanal:
// NewMausRoutine(f func(t uint8, s int8, x uint16, y uint16))
// bzw. NewTastenRoutine(f func(uint16, uint8, uint16))
type bpEingabeRoutine struct {
	hilf.Routine
}

func NewMausRoutine(f func(t uint8, s int8, x uint16, y uint16)) *bpEingabeRoutine {
	routine := hilf.NewRoutine("Maussteuerung",
		func() {
			if !gfx.FensterOffen() {
				return
			}
			taste, status, mausX, mausY := gfx.MausLesen1() // blockiert, bis Maus bedient
			f(taste, status, mausX, mausY)
		})
	routine.SetzeAusnahmeHandler(func() {
		if err := recover(); err != nil {
			if fmt.Sprint(err) == "Das gfx-Fenster ist nicht offen!" ||
				fmt.Sprint(err) == "Das Grafikfenster wurde geschlossen! Programmabbruch!!" {
				// Nichts mehr zu tun
				return
			}
			panic(err)
		}
	})
	return &bpEingabeRoutine{Routine: routine}
}

func NewTastenRoutine(f func(uint16, uint8, uint16)) *bpEingabeRoutine {
	routine := hilf.NewRoutine("Tastensteuerung",
		func() {
			if !gfx.FensterOffen() {
				return
			}
			taste, gedrückt, tiefe := gfx.TastaturLesen1() // blockiert, bis Taste gedrückt
			f(taste, gedrückt, tiefe)
		})
	routine.SetzeAusnahmeHandler(func() {
		if err := recover(); err != nil {
			if fmt.Sprint(err) == "Das gfx-Fenster ist nicht offen!" ||
				fmt.Sprint(err) == "Das Grafikfenster wurde geschlossen! Programmabbruch!!" {
				// Nichts mehr zu tun
				return
			}
			panic(err)
		}
	})
	return &bpEingabeRoutine{Routine: routine}
}
