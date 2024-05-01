package views_controls

type text_overlay struct {
	text string
	//schriftgroesse int
	schreiber *schreiber
	widget
}

// TextOverlay zeigt den Hintergrund
func NewTextOverlay(t string) *text_overlay {
	w := text_overlay{text: t, widget: *NewFenster()}
	w.schreiber = w.monoBoldItalicSchreiber()
	return &w
}

func (f *text_overlay) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	f.stiftfarbe(f.vg)
	breite, höhe := f.GibGroesse()
	sg := int(höhe) / 5
	f.schreiber.SetzeSchriftgroesse(sg)
	f.schreiber.Schreibe(f.startX+breite/3, f.startY+höhe/4, f.text)
}

type infotext struct {
	text string
	//schriftgroesse int
	schreiber *schreiber
	widget
}

// InfoText hat immer einen transparenten Hintergrund.
func NewInfoText(t string) *infotext {
	w := infotext{text: t, widget: *NewFenster()}
	w.schreiber = w.monoBoldSchreiber()
	w.SetzeTransparenz(255)
	return &w
}

func (f *infotext) Zeichne() {
	f.widget.Zeichne()
	f.stiftfarbe(f.vg)

	_, höhe := f.GibGroesse()
	sg := int(höhe) * 3 / 5
	f.schreiber.SetzeSchriftgroesse(sg)
	d := (höhe - uint16(sg)) / 2
	f.schreiber.Schreibe(f.startX+d, f.startY+d, f.text)
}
