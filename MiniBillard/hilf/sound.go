package hilf

import "gfx"

// treffen zwei Kugeln aufeinander
func BallHitsBallSound() {
	// SetzeKlangparameter(rate uint32, aufloesung, kanaele, signal uint8, p float64)
	// rate      ist die Abtastrate, z.B. 11025, 22050 oder 44100.
	// auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
	// kanaele   ist 1 für mono oder 2 für stereo.
	// signal    gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
	// pulsweite (für Rechtecksignale) gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\hilf\\ballHitsBall.wav")
}

func BallInPocketSound() {
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\hilf\\ballIntoPocket.wav")
}

func Bobb() {
	// gfx.SetzeHuellkurve(0.03, 2, 0.5, 0)
	// gfx.SpieleNote("1F", 0.1, 0)
}
