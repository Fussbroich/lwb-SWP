package views

type Fenster interface {
	Zeichne()
	ZeichneRand()
	ZeichneLayout()
	GibStartkoordinaten() (uint16, uint16)
	GibGröße() (uint16, uint16)
}
