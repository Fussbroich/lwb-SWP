package views_controls

import (
	"fmt"
	"gfx"

	"../fonts"
	"../modelle"
)

type KugelZeichner struct {
	kugelPalette *[16]Farbe
}

var (
	standardPoolPalette = [16]Farbe{
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

	englishPoolPalette = [16]Farbe{
		F(252, 253, 242), // weiß
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(249, 73, 68),   // hellrot
		F(48, 49, 54),    // schwarz
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78),  // gelb
		F(255, 201, 78)}  // gelb
)

func (w *KugelZeichner) GibKugelPalette() *[16]Farbe {
	if w.kugelPalette == nil {
		w.kugelPalette = &standardPoolPalette
	}
	return w.kugelPalette
}

func (w *KugelZeichner) SetzeStandardPoolPalette() {
	w.kugelPalette = &standardPoolPalette
}

func (w *KugelZeichner) SetzeEnglishPoolPalette() {
	w.kugelPalette = &englishPoolPalette
}

func (w *KugelZeichner) ZeichneKugel(startX, startY uint16, k modelle.MBKugel) {
	font := fonts.LiberationMonoBold(int(k.GibRadius()) - 3)
	gfxVollKreis(startX, startY, k.GibPos(), k.GibRadius(), F(48, 49, 54))
	gfxVollKreis(startX, startY, k.GibPos(), k.GibRadius()-1, F(252, 253, 242))
	c := w.GibKugelPalette()[k.GibWert()]
	if k.GibWert() <= 8 {
		gfxVollKreis(startX, startY, k.GibPos(), k.GibRadius()-1, c)
	} else {
		r, g, b := c.RGB()
		gfx.Stiftfarbe(r, g, b)
		gfx.Vollrechteck(startX+uint16(k.GibPos().X()-k.GibRadius()*0.75+0.5), startY+uint16(k.GibPos().Y()-k.GibRadius()*0.6+0.5),
			uint16(2*0.75*k.GibRadius()+0.5), uint16(2*0.6*k.GibRadius()+0.5))
		gfxVollKreissektor(startX, startY, k.GibPos(), k.GibRadius()-1, 325, 35, c)
		gfxVollKreissektor(startX, startY, k.GibPos(), k.GibRadius()-1, 145, 215, c)
	}
	// Nur die weiße erhält keine Nummer.
	if k.GibWert() != 0 {
		gfxVollKreis(startX, startY, k.GibPos(), (k.GibRadius()-1)/2, F(252, 253, 242))
		gfx.Stiftfarbe(0, 0, 0)
		gfx.SetzeFont(font.GibDateipfad(), font.GibSchriftgroesse())
		if k.GibWert() < 10 {
			gfx.SchreibeFont(
				startX-uint16(font.GibSchriftgroesse())/4+uint16(k.GibPos().X()+0.5),
				startY-uint16(font.GibSchriftgroesse())/2+uint16(k.GibPos().Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		} else {
			gfx.SchreibeFont(
				startX-uint16(font.GibSchriftgroesse())/2+uint16(k.GibPos().X()+0.5),
				startY-uint16(font.GibSchriftgroesse())/2+uint16(k.GibPos().Y()+0.5),
				fmt.Sprintf("%d", k.GibWert()))
		}
	}
}
