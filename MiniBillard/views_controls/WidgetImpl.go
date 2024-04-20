package views_controls

import (
	"fmt"
	"gfx"

	"../hilf"
)

type widget struct {
	aktiv          bool
	hg, vg         Farbe
	hgName, vgName string
	startX, startY uint16
	stopX, stopY   uint16
	trans          uint8
	eckra          uint16
}

// ########## Methoden für die Konstruktion ############################################

func NewFenster() *widget {
	w := widget{aktiv: true}
	w.SetzeFarben(Fhintergrund(), Ftext())
	return &w
}

func (f *widget) SetzeKoordinaten(startx, starty, stopx, stopy uint16) {
	f.startX, f.startY, f.stopX, f.stopY = startx, starty, stopx, stopy
}

func (f *widget) SetzeFarben(hg, vg string) {
	f.hgName, f.vgName = hg, vg
	f.hg, f.vg = gibFarbe(f.hgName), gibFarbe(f.vgName)
}

func (f *widget) LadeFarben() {
	f.hg, f.vg = gibFarbe(f.hgName), gibFarbe(f.vgName)
}

func (f *widget) SetzeTransparenz(tr uint8) {
	f.trans = tr
}

func (f *widget) SetzeEckradius(ra uint16) {
	f.eckra = ra
}

func (f *widget) GibStartkoordinaten() (uint16, uint16) { return f.startX, f.startY }

func (f *widget) GibGroesse() (uint16, uint16) { return f.stopX - f.startX, f.stopY - f.startY }

// ########## Methoden für die Darstellung ############################################

// Zeichnet einen roten Rand um das Widget herum - zu Testzwecken
func (f *widget) ZeichneLayout() {
	if !f.IstAktiv() {
		return
	}
	br, ho := f.GibGroesse()
	f.stiftfarbe(Rot())
	f.transparenz(0)
	f.rechteckGFX(0, 0, br, ho)
	f.stiftfarbe(f.vg)
	f.transparenz(0)
}

// Zeichnet den Hintergrund und danach den Inhalt
// Aufrufer - die den Inhalt ergänzen wollen, müssen
// ihren Inhalt erst danach zeichnen, sonst wird
// dieser ggf. überdeckt. Transparenz wird beachtet.
func (f *widget) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.stiftfarbe(f.hg)
	f.transparenz(f.trans)
	br, ho := f.GibGroesse()
	if f.eckra > 0 {
		f.vollRechteckGFX(f.eckra, 0, br-2*f.eckra, ho)
		f.vollRechteckGFX(0, f.eckra, br, ho-2*f.eckra)
		f.vollKreisGFX(f.eckra, f.eckra, f.eckra)
		f.vollKreisGFX(f.eckra, ho-f.eckra, f.eckra)
		f.vollKreisGFX(br-f.eckra, ho-f.eckra, f.eckra)
		f.vollKreisGFX(br-f.eckra, f.eckra, f.eckra)
	} else {
		f.vollRechteckGFX(0, 0, br, ho)
	}
	f.stiftfarbe(f.vg)
	f.transparenz(0)
}

// Wie Zeichne, jedoch wird der Hintergrund etwas nach innen eingerückt.
// Kann als Ersatz für Zeichne aufgerufen werden.
func (f *widget) ZeichneOffset(offset uint16) {
	if !f.IstAktiv() {
		return
	}
	f.stiftfarbe(f.hg)
	f.transparenz(f.trans)
	br, ho := f.GibGroesse()
	if f.eckra > 0 {
		f.vollRechteckGFX(f.eckra, offset, br-2*f.eckra, ho-2*offset)
		f.vollRechteckGFX(offset, f.eckra, br-2*offset, ho-2*f.eckra)
		f.vollKreisGFX(f.eckra, f.eckra, f.eckra-offset)
		f.vollKreisGFX(f.eckra, ho-f.eckra, f.eckra-offset)
		f.vollKreisGFX(br-f.eckra, ho-f.eckra, f.eckra-offset)
		f.vollKreisGFX(br-f.eckra, f.eckra, f.eckra-offset)
	} else {
		f.vollRechteckGFX(offset, offset, br-2*offset, ho-2*offset)
	}
	f.stiftfarbe(f.vg)
	f.transparenz(0)
}

// Muss nicht aufgerufen werden, falls kein Rand erscheinen soll.
// Todo: Teilkreise sind als Rand ungeeignet, da die Radien
// von Gfx mit gezeichnet werden. in dem Fall muss der
// Rand *zuerst* gezeichnet werden und der Hintergrund ohne
// Transparenz mit Offset danach.
func (f *widget) ZeichneRand() {
	if !f.IstAktiv() {
		return
	}
	f.stiftfarbe(f.vg)
	f.transparenz(0)
	br, ho := f.GibGroesse()
	if f.eckra > 0 {
		f.kreissektorGFX(f.eckra, f.eckra, f.eckra, 90, 180)
		f.kreissektorGFX(f.eckra, ho-f.eckra, f.eckra, 180, 270)
		f.kreissektorGFX(br-f.eckra, ho-f.eckra, f.eckra, 270, 0)
		f.kreissektorGFX(br-f.eckra, f.eckra, f.eckra, 0, 90)
		f.LinieGFX(0, f.eckra, 0, ho-f.eckra)
		f.LinieGFX(f.eckra, ho, br-f.eckra, ho)
		f.LinieGFX(br, ho-f.eckra, br, f.eckra)
		f.LinieGFX(br-f.eckra, 0, f.eckra, 0)
	} else {
		f.LinieGFX(0, 0, br, 0)
		f.LinieGFX(0, ho, br, ho)
		f.LinieGFX(0, 0, 0, ho)
		f.LinieGFX(br, 0, br, ho)
	}
}

// ########## Methoden zum Ein- und Ausblenden ############################################

func (f *widget) IstAktiv() bool { return f.aktiv }

func (f *widget) DarstellenAnAus() { f.aktiv = !f.aktiv }

func (f *widget) Einblenden() { f.aktiv = true }

func (f *widget) Ausblenden() { f.aktiv = false }

// ########## Methoden für die Maussteuerung ############################################

func (f *widget) ImFenster(x, y uint16) bool {
	if !f.IstAktiv() {
		return false
	}
	xs, ys := f.startX+f.eckra*3/10, f.startY+f.eckra*3/10
	b, h := f.stopX-f.startX-f.eckra*6/10, f.stopY-f.startY-f.eckra*6/10
	return x > xs && x < xs+b && y > ys && y < ys+h
}

func (f *widget) MausklickBei(x, y uint16) {
	if !f.IstAktiv() {
		return
	}
	fmt.Println("Unbeachteter Mausklick bei", x, y)
}

func (f *widget) MausBei(x, y uint16) {
	if !f.IstAktiv() {
		return
	}
	fmt.Println("Unbeachtete Mausbewegung bei", x, y)
}

func (f *widget) MausScrolltHoch() {
	if !f.IstAktiv() {
		return
	}
	fmt.Println("Unbeachtetes Mausscrollen")
}

func (f *widget) MausScrolltRunter() {
	if !f.IstAktiv() {
		return
	}
	fmt.Println("Unbeachtetes Mausscrollen")
}

// ######## Hilfsmethoden zum Zeichnen ############################################################

func (f *widget) stiftfarbe(c Farbe) {
	r, g, b := c.RGB()
	f.stiftfarbeGFX(r, g, b)
}

func (f *widget) stiftfarbeGFX(r, g, b uint8) {
	gfx.Stiftfarbe(r, g, b)
}

func (f *widget) transparenz(tr uint8) {
	gfx.Transparenz(tr)
}

func (f *widget) LinieGFX(xV, yV, xN, yN uint16) {
	gfx.Linie(f.startX+xV, f.startY+yV, f.startX+xN, f.startY+yN)
}

func (f *widget) vollRechteckGFX(xV, yV, b, h uint16) {
	gfx.Vollrechteck(f.startX+xV, f.startY+yV, b, h)
}

func (f *widget) rechteckGFX(xV, yV, b, h uint16) {
	gfx.Rechteck(f.startX+xV, f.startY+yV, b, h)
}

func (f *widget) vollKreis(pos hilf.Vec2, radius float64, c Farbe) {
	f.stiftfarbe(c)
	f.vollKreisGFX(uint16(0.5+pos.X()), uint16(0.5+pos.Y()), uint16(0.5+radius))
}

func (f *widget) vollKreisGFX(x, y, ra uint16) {
	gfx.Vollkreis(f.startX+x, f.startY+y, ra)
}

func (f *widget) kreisGFX(x, y, ra uint16) {
	gfx.Kreis(f.startX+x, f.startY+y, ra)
}

func (f *widget) kreissektorGFX(x, y, ra, wVon, wBis uint16) {
	gfx.Kreissektor(f.startX+x, f.startY+y, ra, wVon, wBis)
}

func (f *widget) vollKreissektor(pos hilf.Vec2, radius float64, wVon, wBis uint16, c Farbe) {
	f.stiftfarbe(c)
	f.vollKreissektorGFX(uint16(0.5+pos.X()), uint16(0.5+pos.Y()), uint16(0.5+radius), wVon, wBis)
}

func (f *widget) vollKreissektorGFX(x, y, ra, wVon, wBis uint16) {
	gfx.Vollkreissektor(f.startX+x, f.startY+y, ra, wVon, wBis)
}

func (f *widget) vollDreieck(pA, pB, pC hilf.Vec2) {
	f.vollDreieckGFX(
		uint16(0.5+pA.X()), uint16(0.5+pA.Y()),
		uint16(0.5+pB.X()), uint16(0.5+pB.Y()),
		uint16(0.5+pC.X()), uint16(0.5+pC.Y()))
}

func (f *widget) vollDreieckGFX(x1, y1, x2, y2, x3, y3 uint16) {
	gfx.Volldreieck(
		f.startX+x1, f.startY+y1,
		f.startX+x2, f.startY+y2,
		f.startX+x3, f.startY+y3)
}

func (f *widget) breiteLinie(pV, pN hilf.Vec2, breite float64, c Farbe) {
	richt := pN.Minus(pV).Normiert()
	d := hilf.V2(richt.Y(), -richt.X())
	f.stiftfarbe(c)

	pA := pV.Minus(d.Mal(breite / 2))
	pB := pV.Plus(d.Mal(breite / 2))
	pC := pN.Plus(d.Mal(breite / 2))
	pD := pN.Minus(d.Mal(breite / 2))

	f.vollDreieck(pA, pB, pC)
	f.vollDreieck(pA, pC, pD)
}
