package views_controls

import "../assets"

var (
	monoBoldFont       = assets.FontDateipfad("LiberationMono-Bold.ttf")
	monoRegularFont    = assets.FontDateipfad("LiberationMono-Regular.ttf")
	monoBoldItalicFont = assets.FontDateipfad("LiberationMono-BoldItalic.ttf")
)

func (f *widget) monoBoldSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      monoBoldFont,
		schriftgroesse: 24}
}

func (f *widget) monoRegularSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      monoRegularFont,
		schriftgroesse: 24}
}

func (f *widget) monoBoldItalicSchreiber() *schreiber {
	return &schreiber{
		fontdatei:      monoBoldItalicFont,
		schriftgroesse: 24}
}
