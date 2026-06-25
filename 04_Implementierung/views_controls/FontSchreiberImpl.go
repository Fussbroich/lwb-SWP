package views_controls

import (
	"brainpool/assets"
	"brainpool/gfx"
)

type schreiber struct {
	fontDaten      []byte
	fontName       string
	schriftgroesse int
}

type FontStyle uint8

const (
	Bold FontStyle = iota
	Regular
	BoldItalic
)

type fontInfo struct {
	daten []byte
	name  string
}

var (
	standardFonts = map[FontStyle]fontInfo{
		Bold:       {assets.MonoBoldFontDaten(), "MonoBold"},
		Regular:    {assets.MonoRegularFontDaten(), "MonoRegular"},
		BoldItalic: {assets.MonoBoldItalicFontDaten(), "MonoBoldItalic"}}
	standardSchriftgroesse int = 12
)

func (f *widget) newSchreiber(style FontStyle) *schreiber {
	fi := standardFonts[style]
	return &schreiber{
		fontDaten:      fi.daten,
		fontName:       fi.name,
		schriftgroesse: standardSchriftgroesse}
}

func (s *schreiber) SetzeSchriftgroesse(groesse int) {
	s.schriftgroesse = groesse
}

func (s *schreiber) GibSchriftgroesse() int {
	return s.schriftgroesse
}

func (s *schreiber) Schreibe(x, y uint16, text string) {
	gfx.SetzeFontDaten(s.fontDaten, s.fontName, s.schriftgroesse)
	gfx.SchreibeFont(x, y, text)
}
