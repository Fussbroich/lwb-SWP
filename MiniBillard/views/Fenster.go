package views

type Fenster interface {
	Zeichne()
	ZeichneLayout()
	GibStartkoordinaten() (uint16, uint16)
	GibGröße() (uint16, uint16)
}
