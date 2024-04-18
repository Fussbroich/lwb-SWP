package views_controls

import (
	"gfx"

	"../assets"
)

type schreiber struct {
	fontdatei      string
	schriftgroesse int
}

var (
	standardSchriftgroesse int = 12
)

func (f *widget) monoBoldSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      assets.MonoBoldFontDateipfad(),
		schriftgroesse: standardSchriftgroesse}
}

func (f *widget) monoRegularSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      assets.MonoRegularFontDateipfad(),
		schriftgroesse: standardSchriftgroesse}
}

func (f *widget) monoBoldItalicSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      assets.MonoBoldItalicFontDateipfad(),
		schriftgroesse: standardSchriftgroesse}
}

func (s *schreiber) SetzeSchriftgroesse(groesse int) {
	s.schriftgroesse = groesse
}

func (s *schreiber) GibSchriftgroesse() int {
	return s.schriftgroesse
}

func (s *schreiber) Schreibe(x, y uint16, text string) {
	gfx.SetzeFont(s.fontdatei, s.schriftgroesse)
	gfx.SchreibeFont(x, y, text)
}
