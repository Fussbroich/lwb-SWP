package views

type Fenster interface {
	Zeichne()
	GibStartkoordinaten() (uint16, uint16)
	GibGröße() (uint16, uint16)
}
