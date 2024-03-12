package views

import (
	"fmt"
	"gfx"

	"../hilf"
	"../welt"
)

type FensterZeichner struct {
	startX, startY uint16
	stopX, stopY   uint16
}

type HintergrundZeichner struct {
	FensterZeichner
}

type MiniBillardSpielfeldZeichner struct {
	FensterZeichner
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
			stopX: stopx, stopY: stopy}}
}

func NewHintergrundZeichner(startx, starty, stopx, stopy uint16) *HintergrundZeichner {
	return &HintergrundZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy}}
}

func NewMBEingelochteZeichner(startx, starty, stopx, stopy uint16) *MiniBillardEingelochteZeichner {
	return &MiniBillardEingelochteZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy}}
}

func NewMBSpielinfoZeichner(startx, starty, stopx, stopy uint16) *MiniBillardSpielinfoZeichner {
	return &MiniBillardSpielinfoZeichner{
		FensterZeichner{
			startX: startx, startY: starty,
			stopX: stopx, stopY: stopy}}
}

func (f *MiniBillardSpielfeldZeichner) Zeichne(spiel welt.MiniBillardSpiel) {
	gfx.Cls()
	l, b := spiel.GibGröße()
	// zeichne den Belag
	f.ZeichneVollRechteck(hilf.V2(0, 0), l, b, 60, 179, 113)
	// zeichne die Taschen
	ts := spiel.GibTaschen()
	f.ZeichneVollKreissektor(ts[0].GibPos(), ts[0].GibRadius(), 270, 0, 0, 0, 0)
	f.ZeichneVollKreissektor(ts[1].GibPos(), ts[1].GibRadius(), 0, 90, 0, 0, 0)
	f.ZeichneVollKreissektor(ts[2].GibPos(), ts[2].GibRadius(), 0, 180, 0, 0, 0)
	f.ZeichneVollKreissektor(ts[3].GibPos(), ts[3].GibRadius(), 90, 180, 0, 0, 0)
	f.ZeichneVollKreissektor(ts[4].GibPos(), ts[4].GibRadius(), 180, 270, 0, 0, 0)
	f.ZeichneVollKreissektor(ts[5].GibPos(), ts[5].GibRadius(), 180, 360, 0, 0, 0)
	// zeichne die Kugeln
	for _, k := range spiel.GibAktiveKugeln() {
		r, g, b := k.GibFarbe()
		f.ZeichneVollKreis(k.GibPos(), k.GibRadius(), 0, 0, 0)
		f.ZeichneVollKreis(k.GibPos(), k.GibRadius()-1, r, g, b)
	}
}

func (f *MiniBillardEingelochteZeichner) Zeichne(spiel welt.MiniBillardSpiel) {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	// zeichne den Hintergrund
	gfx.Stiftfarbe(80, 80, 80)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	gfx.Stiftfarbe(180, 180, 180)
	schriftgröße := höhe / 5
	gfx.SetzeFont("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\src\\gfx\\fonts\\LiberationMono-Bold.ttf",
		int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Eingelocht", len(spiel.GibEingelochteKugeln())))
}

func (f *MiniBillardSpielinfoZeichner) Zeichne(spiel welt.MiniBillardSpiel) {
	breite := f.stopX - f.startX
	höhe := f.stopY - f.startY
	// zeichne den Hintergrund
	gfx.Stiftfarbe(80, 80, 80)
	gfx.Vollrechteck(f.startX, f.startY, breite, höhe)
	//schreibe Stößezahl
	gfx.Stiftfarbe(180, 180, 180)
	schriftgröße := höhe / 5
	gfx.SetzeFont("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\src\\gfx\\fonts\\LiberationMono-Bold.ttf",
		int(schriftgröße))
	gfx.SchreibeFont(f.startX, f.startY,
		fmt.Sprintf("%d Stöße", spiel.GibStößeBisher()))
	gfx.SchreibeFont(f.startX, f.startY+12*schriftgröße/10,
		fmt.Sprintf("%d Strafpunkte", spiel.GibStrafpunkte()))
}

func (f *HintergrundZeichner) Zeichne(r, g, b uint8) {
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(f.startX, f.startY, f.stopX-f.startX, f.stopY-f.startY)
}

func (f *FensterZeichner) ZeichneVollKreis(pos hilf.Vec2, radius float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreis(
		f.startX+uint16(0.5+pos.X()), f.startY+uint16(0.5+pos.Y()),
		uint16(0.5+radius))
}

func (f *FensterZeichner) ZeichneVollKreissektor(pos hilf.Vec2, radius float64, wVon, wBis uint16, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollkreissektor(
		f.startX+uint16(0.5+pos.X()), f.startY+uint16(0.5+pos.Y()),
		uint16(0.5+radius), wVon, wBis)
}

func (f *FensterZeichner) ZeichneVollRechteck(pos hilf.Vec2, breite, höhe float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)
	gfx.Vollrechteck(
		f.startX+uint16(0.5+pos.X()), f.startY+uint16(0.5+pos.Y()),
		uint16(0.5+breite), uint16(0.5+höhe))
}

func (f *FensterZeichner) ZeichneVollDreieck(p1, p2, p3 hilf.Vec2) {
	gfx.Volldreieck(
		f.startX+uint16(0.5+p1.X()), f.startY+uint16(0.5+p1.Y()),
		f.startX+uint16(0.5+p2.X()), f.startY+uint16(0.5+p2.Y()),
		f.startX+uint16(0.5+p3.X()), f.startY+uint16(0.5+p3.Y()))
}

func (f *FensterZeichner) ZeichneBreiteLinie(pV, pN hilf.Vec2, breite float64, cr, cg, cb uint8) {
	gfx.Stiftfarbe(cr, cg, cb)

	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))
	f.ZeichneVollDreieck(pA, pB, pC)
	f.ZeichneVollDreieck(pA, pC, pD)
	gfx.Vollkreis(
		f.startX+uint16(0.5+pV.X()), f.startY+uint16(0.5+pV.Y()), uint16(0.5+breite/2))
	gfx.Vollkreis(
		f.startX+uint16(0.5+pN.X()), f.startX+uint16(0.5+pN.Y()), uint16(0.5+breite/2))
}
