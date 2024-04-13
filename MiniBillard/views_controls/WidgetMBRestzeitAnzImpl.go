package views_controls

import (
	"fmt"
	"time"

	"../modelle"
)

type miniBRestzeit struct {
	billard modelle.MiniBillardSpiel
	widget
}

func NewMBRestzeitAnzeiger(billard modelle.MiniBillardSpiel) *miniBRestzeit {
	return &miniBRestzeit{billard: billard, widget: *NewFenster()}
}

func fmtRestzeit(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m%60, s)
}

func (f *miniBRestzeit) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()
	schreiber := f.LiberationMonoBoldSchreiber()
	schreiber.SetzeSchriftgroesse(int(höhe) * 4 / 5)
	anzeige := fmtRestzeit(f.billard.GibRestzeit())
	dx := (breite - uint16(len(anzeige)*schreiber.GibSchriftgroesse()*3/5)) / 2
	dy := (höhe - uint16(schreiber.GibSchriftgroesse())) / 2
	schreiber.Schreibe(f.startX+dx, f.startY+dy, anzeige)
}
