package views

import (
	"fmt"
	"gfx"
	"time"

	"../welt"
)

type miniBRestzeit struct {
	billard welt.MiniBillardSpiel
	fenster
}

func NewMBRestzeitAnzeiger(billard welt.MiniBillardSpiel, startx, starty, stopx, stopy uint16, hg, vg Farbe, tr uint8) *miniBRestzeit {
	fenster := fenster{startX: startx, startY: starty, stopX: stopx, stopY: stopy, hg: hg, vg: vg, transparenz: tr}
	return &miniBRestzeit{billard: billard, fenster: fenster}
}

func fmtRestzeit(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m%60, s)
}
func (f *miniBRestzeit) Zeichne() {
	f.fenster.Zeichne()
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	breite, höhe := f.GibGröße()
	schriftgröße := int(höhe) * 4 / 5
	anzeige := fmtRestzeit(f.billard.GibRestzeit())
	dx := (breite - uint16(len(anzeige)*schriftgröße*3/5)) / 2
	dy := (höhe - uint16(schriftgröße)) / 2
	gfx.SetzeFont(fp, schriftgröße)
	gfx.SchreibeFont(f.startX+dx, f.startY+dy, anzeige)
}
