package views_controls

import "gfx"

type schreiber struct {
	fontdatei      string
	schriftgroesse int
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
