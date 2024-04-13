package views_controls

type text_overlay struct {
	text string
	widget
}

type infotext struct {
	text string
	widget
}

// TextOverlay zeigt den Hintergrund
func NewTextOverlay(t string) *text_overlay {
	return &text_overlay{text: t, widget: *NewFenster()}
}

func (f *text_overlay) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	f.stiftfarbe(f.vg)
	schreiber := f.liberationMonoBoldItalicSchreiber()
	schreiber.SetzeSchriftgroesse(int(f.stopY-f.startY) / 5)
	schreiber.Schreibe((f.stopX-f.startX)/3, (f.stopY-f.startY)/4, f.text)
}

// InfoText ist immer Transparent
func NewInfoText(t string) *infotext {
	return &infotext{text: t, widget: widget{hg: Weiss(), trans: 255}}
}

func (f *infotext) Zeichne() {
	f.widget.Zeichne()
	f.stiftfarbe(f.vg)

	_, höhe := f.GibGroesse()
	schreiber := f.liberationMonoBoldItalicSchreiber()
	schreiber.SetzeSchriftgroesse(int(höhe) * 3 / 5)
	d := (höhe - uint16(schreiber.GibSchriftgroesse())) / 2
	schreiber.Schreibe(f.startX+d, f.startY+d, f.text)
}
