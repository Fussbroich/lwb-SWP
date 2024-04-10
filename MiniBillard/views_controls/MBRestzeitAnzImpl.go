package views_controls

import (
	"fmt"
	"gfx"
	"time"

	"../fonts"
	"../modelle"
)

type miniBRestzeit struct {
	billard modelle.MiniBillardSpiel
	widget
}

func NewMBRestzeitAnzeiger(billard modelle.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBRestzeit {
	fenster := widget{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &miniBRestzeit{billard: billard, widget: fenster}
}

func fmtRestzeit(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m%60, s)
}

func (f *miniBRestzeit) Zeichne() {
	f.widget.Zeichne()
	breite, höhe := f.GibGröße()
	font := fonts.LiberationMonoBold(int(höhe) * 4 / 5)
	anzeige := fmtRestzeit(f.billard.GibRestzeit())
	dx := (breite - uint16(len(anzeige)*font.GibSchriftgröße()*3/5)) / 2
	dy := (höhe - uint16(font.GibSchriftgröße())) / 2
	gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgröße())
	gfx.SchreibeFont(f.startX+dx, f.startY+dy, anzeige)
}
