package views_controls

import (
	"fmt"
	"time"
)

type zeitAnzeige struct {
	uhrzeit   *time.Time
	schreiber *schreiber
	widget
}

func NewDigitalUhrzeitAnzeiger(t *time.Time) *zeitAnzeige {
	f := *NewFenster()
	return &zeitAnzeige{uhrzeit: t, widget: f}
}

func (f *zeitAnzeige) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()

	if f.schreiber == nil {
		f.schreiber = f.newSchreiber(Bold)
		f.schreiber.SetzeSchriftgroesse(int(höhe) * 3 / 5)
	}
	t_str := fmt.Sprintf("%02d:%02d", f.uhrzeit.Hour(), f.uhrzeit.Minute())
	secsLen := breite * uint16(f.uhrzeit.Second()) / 60
	f.vollRechteckGFX(0, höhe-8, secsLen, 8)
	dx := (breite - uint16(len(t_str)*f.schreiber.GibSchriftgroesse()*3/5)) / 2
	dy := (höhe - uint16(f.schreiber.GibSchriftgroesse())) / 2
	f.schreiber.Schreibe(f.startX+dx, f.startY+dy, t_str)
	f.ZeichneRand()
}
