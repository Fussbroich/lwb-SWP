package views_controls

import (
	"fmt"
	"time"
)

type uhrzeit struct {
	uhrzeit   *time.Time
	schreiber *schreiber
	widget
}

func NewDigitalUhrzeitAnzeiger(t *time.Time) *uhrzeit {
	f := *NewFenster()
	return &uhrzeit{uhrzeit: t, widget: f}
}

func (f *uhrzeit) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()
	if f.schreiber == nil {
		f.schreiber = f.newSchreiber(Bold)
		f.schreiber.SetzeSchriftgroesse(int(höhe) * 3 / 5)
	}
	h, m, s := f.uhrzeit.Hour(), f.uhrzeit.Minute(), f.uhrzeit.Second()
	t_str := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	dx := (breite - uint16(len(t_str)*f.schreiber.GibSchriftgroesse()*3/5)) / 2
	dy := (höhe - uint16(f.schreiber.GibSchriftgroesse())) / 2
	f.schreiber.Schreibe(f.startX+dx, f.startY+dy, t_str)
}
