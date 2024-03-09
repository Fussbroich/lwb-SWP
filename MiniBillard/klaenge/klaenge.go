package klaenge

import "gfx"

// SetzeKlangparameter(rate uint32, aufloesung, kanaele, signal uint8, p float64)
// 	rate      ist die Abtastrate, z.B. 11025, 22050 oder 44100.
// 	auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
// 	kanaele   ist 1 für mono oder 2 für stereo.
// 	signal    gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
// 	pulsweite (für Rechtecksignale) gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.

func MassivePulseLoopSound() {
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 0.3)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\massivePulseLoop.wav")
}

func CueHitsBallSound() {
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\cueHitsBall.wav")
}

func BallHitsBallSound() {
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\ballHitsBall.wav")
}

func BallInPocketSound() {
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\ballIntoPocket.wav")
}

func BallHitsRailSound() {
	gfx.SetzeKlangparameter(22050, 2, 2, 1, 1.0)
	gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\klaenge\\ballHitsRail.wav")
}
