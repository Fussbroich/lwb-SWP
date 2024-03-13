package views

import (
	"errors"
	"fmt"
	"gfx"
	"os"
	"path/filepath"

	"../hilf"
	"../welt"
)

var (
	kugelPalette *[16]Farbe
)

type FensterZeichner struct {
	startX, startY uint16
	stopX, stopY   uint16
	hg, vg         Farbe
}

type HintergrundZeichner struct {
	FensterZeichner
}

type MiniBillardSpielfeldZeichner struct {
	FensterZeichner
}

func (f *FensterZeichner) KugelPalette() *[16]Farbe {
	if kugelPalette == nil {
		kugelPalette = &[16]Farbe{
			F(252, 253, 242), // weiß
			F(255, 201, 78),  // gelb
			F(34, 88, 175),   // blau
			F(249, 73, 68),   // hellrot
			F(84, 73, 149),   // violett
			F(255, 139, 33),  // orange
			F(47, 159, 52),   // grün
			F(155, 53, 30),   // dunkelrot
			F(48, 49, 54),    // schwarz
			F(255, 201, 78),  // gelb
			F(34, 88, 175),   // blau
			F(249, 73, 68),   // hellrot
			F(84, 73, 149),   // violett
			F(255, 139, 33),  // orange
			F(47, 159, 52),   // grün
			F(155, 53, 30)}   // dunkelrot
	}
	return kugelPalette
}

type MiniBillardSpielinfoZeichner struct {
	FensterZeichner
}

type MiniBillardEingelochteZeichner struct {
	FensterZeichner
}

func NewMBSpielfeldZeichner(startx, starty, stopx, stopy uint16) *MiniBillardSpielfeldZeichner {
	return &MiniBillardSpielfeldZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy,
			hg: Weiß(), vg: Schwarz()}}
}

func NewHintergrundZeichner(startx, starty, stopx, stopy uint16, hg Farbe) *HintergrundZeichner {
	return &HintergrundZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy,
			hg: hg, vg: Schwarz()}}
}

func NewMBEingelochteZeichner(startx, starty, stopx, stopy uint16, hg Farbe) *MiniBillardEingelochteZeichner {
	return &MiniBillardEingelochteZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy,
			hg: hg, vg: Schwarz()}}
}

func NewMBSpielinfoZeichner(startx, starty, stopx, stopy uint16, hg, vg Farbe) *MiniBillardSpielinfoZeichner {
	return &MiniBillardSpielinfoZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy,
			hg: hg, vg: vg}}
}

func fontDateipfad(filename string) string {
	fontsDir := "MiniBillard/fonts"
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}
	wdir := filepath.Dir(wd)
	fp := filepath.Join(wdir, "lwb-SWP", fontsDir, filename)
	//	println("wdir:", wdir)
	//	println("klaengeDir:", klaengeDir)
	//	println("filename", filename)
	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return fp
}

func (f *MiniBillardSpielfeldZeichner) Zeichne(spiel welt.MiniBillardSpiel) {
	gfx.Cls()
	l, b := spiel.GibGröße()
	// zeichne den Belag
	f.ZeichneVollRechteck(hilf.V2(0, 0), l, b, F(60, 179, 113))
	// zeichne die Taschen
	ts := spiel.GibTaschen()
	f.ZeichneVollKreissektor(ts[0].GibPos(), ts[0].GibRadius(), 270, 0, Schwarz())
	f.ZeichneVollKreissektor(ts[1].GibPos(), ts[1].GibRadius(), 0, 90, Schwarz())
	f.ZeichneVollKreissektor(ts[2].GibPos(), ts[2].GibRadius(), 0, 180, Schwarz())
	f.ZeichneVollKreissektor(ts[3].GibPos(), ts[3].GibRadius(), 90, 180, Schwarz())
	f.ZeichneVollKreissektor(ts[4].GibPos(), ts[4].GibRadius(), 180, 270, Schwarz())
	f.ZeichneVollKreissektor(ts[5].GibPos(), ts[5].GibRadius(), 180, 360, Schwarz())
	// zeichne die Kugeln
	for _, k := range spiel.GibAktiveKugeln() {
		f.zeichneKugel(k)
	}

}

func (f *MiniBillardSpielfeldZeichner) zeichneKugel(k welt.Kugel) {
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	schriftgröße := int(k.GibRadius()) - 3
	p := k.GibPos()
	f.ZeichneVollKreis(p, k.GibRadius(), F(48, 49, 54))
	f.ZeichneVollKreis(p, k.GibRadius()-1, F(252, 253, 242))
	c := f.KugelPalette()[k.GibWert()]
	if k.GibWert() <= 8 {
		f.ZeichneVollKreis(p, k.GibRadius()-1, c)
	} else {
		r, g, b := c.RGB()
		gfx.Stiftfarbe(r, g, b)
		gfx.Vollrechteck(f.startX+uint16(p.X()-k.GibRadius()*0.818+0.5), f.startY+uint16(p.Y()-k.GibRadius()*0.61+0.5),
			uint16(2*0.818*k.GibRadius()+0.5), uint16(2*0.61*k.GibRadius()+0.5))
		f.ZeichneVollKreissektor(p, k.GibRadius()-1, 325, 35, c)
		f.ZeichneVollKreissektor(p, k.GibRadius()-1, 145, 215, c)
	}
	if k.GibWert() != 0 {
		f.ZeichneVollKreis(p, (k.GibRadius()-1)/2, F(252, 253, 242))
		gfx.Stiftfarbe(0, 0, 0)
		gfx.SetzeFont(fp, schriftgröße)
		if k.GibWert() < 10 {
			gfx.SchreibeFont(
				f.startX-uint16(schriftgröße)/4+uint16(p.X()+0.5),
				f.startY-uint16(schriftgröße)/2+uint16(p.Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		} else {
			gfx.SchreibeFont(
				f.startX-uint16(schriftgröße)/2+uint16(p.X()+0.5),
				f.startY-uint16(schriftgröße)/2+uint16(p.Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		}
	}
}

func (f *MiniBillardEingelochteZeichner) Zeichne(spiel welt.MiniBillardSpiel) {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	// zeichne den Hintergrund
	gfx.Stiftfarbe(80, 80, 80)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	gfx.Stiftfarbe(180, 180, 180)
	schriftgröße := höhe / 3
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Eingelocht", len(spiel.GibEingelochteKugeln())))
}

func (f *MiniBillardSpielinfoZeichner) Zeichne(spiel welt.MiniBillardSpiel) {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	fp := fontDateipfad("LiberationMono-Bold.ttf")
	// zeichne den Hintergrund
	cr, cg, cb := f.hg.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	cr, cg, cb = f.vg.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	schriftgröße := höhe / 4
	gfx.SetzeFont(fp, int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Stöße", spiel.GibStößeBisher()))
	gfx.SchreibeFont(f.startX, f.startY+12*schriftgröße/10,
		fmt.Sprintf("%d Strafpunkte", spiel.GibStrafpunkte()))
}

func (f *HintergrundZeichner) Zeichne() {
	r, g, b := f.hg.RGB()
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
}

func (f *FensterZeichner) ZeichneVollKreis(pos hilf.Vec2, radius float64, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		f.startX+uint16(0.5+pos.X()), f.startY+uint16(0.5+pos.Y()),
		uint16(0.5+radius))
}

func (f *FensterZeichner) ZeichneVollKreissektor(pos hilf.Vec2, radius float64, wVon, wBis uint16, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreissektor(
		f.startX+uint16(0.5+pos.X()), f.startY+uint16(0.5+pos.Y()),
		uint16(0.5+radius), wVon, wBis)
}

func (f *FensterZeichner) ZeichneVollRechteck(pos hilf.Vec2, breite, höhe float64, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(
		f.startX+uint16(0.5+pos.X()), f.startY+uint16(0.5+pos.Y()),
		uint16(0.5+breite), uint16(0.5+höhe))
}

func (f *FensterZeichner) ZeichneVollDreieck(p1, p2, p3 hilf.Vec2, c Farbe) {
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Volldreieck(
		f.startX+uint16(0.5+p1.X()), f.startY+uint16(0.5+p1.Y()),
		f.startX+uint16(0.5+p2.X()), f.startY+uint16(0.5+p2.Y()),
		f.startX+uint16(0.5+p3.X()), f.startY+uint16(0.5+p3.Y()))
}

func (f *FensterZeichner) ZeichneBreiteLinie(pV, pN hilf.Vec2, breite float64, c Farbe) {
	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	f.ZeichneVollDreieck(pA, pB, pC, c)
	f.ZeichneVollDreieck(pA, pC, pD, c)
	cr, cg, cb := c.RGB()
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		f.startX+uint16(0.5+pV.X()), f.startY+uint16(0.5+pV.Y()), uint16(0.5+breite/2))
	gfx.Vollkreis(
		f.startX+uint16(0.5+pN.X()), f.startX+uint16(0.5+pN.Y()), uint16(0.5+breite/2))
}
