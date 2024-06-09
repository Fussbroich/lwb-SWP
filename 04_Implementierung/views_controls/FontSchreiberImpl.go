package views_controls

import (
	"gfx"

	"../assets"
)

// Ein FontSchreiber ist ein Wrapper-Objekt für einen Schrifttyp und eine Schriftgröße und
// integriert die Funktionen gfx.SetzeFont und SchreibeFont.
type schreiber struct {
	font           string
	schriftgroesse int
}

type FontStyle uint8

const (
	Bold FontStyle = iota
	Regular
	BoldItalic
)

var (
	standardFonts = map[FontStyle]string{
		Bold:       assets.MonoBoldFontDateipfad(),
		Regular:    assets.MonoRegularFontDateipfad(),
		BoldItalic: assets.MonoBoldItalicFontDateipfad()}
	standardSchriftgroesse int = 12
)

func (f *widget) newSchreiber(style FontStyle) *schreiber {
	return &schreiber{
		font:           standardFonts[style],
		schriftgroesse: standardSchriftgroesse}
}

func (s *schreiber) SetzeSchriftgroesse(groesse int) {
	s.schriftgroesse = groesse
}

func (s *schreiber) GibSchriftgroesse() int {
	return s.schriftgroesse
}

func (s *schreiber) Schreibe(x, y uint16, text string) {
	gfx.SetzeFont(s.font, s.schriftgroesse)
	gfx.SchreibeFont(x, y, text)
}
