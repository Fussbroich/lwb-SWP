package views_controls

import (
	"fmt"
	"time"
)

type uhrzeit struct {
	uhrzeit *time.Time
	t_str   string
	widget
}

func NewDigitalUhrzeitAnzeiger(t *time.Time) *uhrzeit {
	return &uhrzeit{uhrzeit: t, widget: *NewFenster()}
}

func (f *uhrzeit) Update() {
	h, m, s := f.uhrzeit.Hour(), f.uhrzeit.Minute(), f.uhrzeit.Second()
	f.t_str = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func (f *uhrzeit) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()
	schreiber := f.newSchreiber(Bold)
	schreiber.SetzeSchriftgroesse(int(höhe) * 2 / 5)
	dx := (breite - uint16(len(f.t_str)*schreiber.GibSchriftgroesse()*3/5)) / 2
	dy := (höhe - uint16(schreiber.GibSchriftgroesse())) / 2
	f.stiftfarbe(f.vg)
	schreiber.Schreibe(f.startX+dx, f.startY+dy, f.t_str)
}
