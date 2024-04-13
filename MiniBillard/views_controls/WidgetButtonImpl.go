package views_controls

type button struct {
	text   string
	action func()
	widget
}

// Buttons haben einen Text in der Mitte
func NewButton(t string, action func()) *button {
	return &button{text: t, action: action, widget: *NewFenster()}
}

func (f *button) MausklickBei(mausX, mausY uint16) {
	if f.ImFenster(mausX, mausY) {
		f.action()
	}
}

func (f *button) Zeichne() {
	f.ZeichneRand()
	f.widget.Zeichne()
	breite, höhe := f.GibGroesse()

	schreiber := f.liberationMonoRegularSchreiber()
	schreiber.SetzeSchriftgroesse(int(höhe) * 3 / 5)
	f.stiftfarbe(f.vg)

	d := (höhe - uint16(schreiber.GibSchriftgroesse())) / 2

	schreiber.Schreibe(
		f.startX+(breite/2)-uint16(len(f.text)*schreiber.GibSchriftgroesse()*7/24), f.startY+d,
		f.text)
}
