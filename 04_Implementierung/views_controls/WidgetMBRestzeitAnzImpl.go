package views_controls

import (
	"fmt"
	"time"

	"../modelle"
)

type miniBRestzeit struct {
	billard  modelle.MiniBillardSpiel
	restzeit time.Duration
	rzString string
	widget
}

func NewMBRestzeitAnzeiger(billard modelle.MiniBillardSpiel) *miniBRestzeit {
	return &miniBRestzeit{billard: billard, widget: *NewFenster()}
}

func (f *miniBRestzeit) Update() {
	rz := f.billard.GibRestzeit().Round(time.Second)
	f.veraltet = rz != f.restzeit
	if !f.veraltet {
		return
	}
	f.restzeit = rz
	m := rz / time.Minute
	rz -= m * time.Minute
	s := rz / time.Second
	f.rzString = fmt.Sprintf("%02d:%02d", m%60, s)
}

func (f *miniBRestzeit) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()
	schreiber := f.newSchreiber(Bold)
	schreiber.SetzeSchriftgroesse(int(höhe) * 4 / 5)
	dx := (breite - uint16(len(f.rzString)*schreiber.GibSchriftgroesse()*3/5)) / 2
	dy := (höhe - uint16(schreiber.GibSchriftgroesse())) / 2
	schreiber.Schreibe(f.startX+dx, f.startY+dy, f.rzString)
}
