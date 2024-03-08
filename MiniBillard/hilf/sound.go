package hilf

import "gfx"

/*
var (
	anschlagzeit  float64 // 0ms <= a <= 1ms
	abschwellzeit float64 // 0 <= d <= 5ms
	haltepegel    float64 // in Prozent vom Maximum mit 0 <= s <= 1.0
	ausklingzeit  float64 // 0 <= r <= 5ms
)
*/
// Das erste Zeichen von tonname ist eine Ziffer von 0 bis 9 und gibt die Oktave an.
// Erlaubte weitere Zeichen für den Notennamen sind "C","D","E","F","G","A","H","C#","D#","F#","G#","A#".
// 0 < laenge <= 1;  laenge 1: volle Note; 1.0/2: halbe Note, ..., 1.0/16: sechzehntel Note
// 0.0<=wartedauer; Die Wartedauer gibt die Dauer in Notenlänge an, nach der nach dem Anspielen der
// Note im Programmablauf fortgefahren wird. 0: keine Wartedauer; 1.0/2: Dauer einer halben Note, ...
// Es werden gerade höchstens 9 Noten oder WAV-Dateien abgespielt.
//       Der voreingestellte Standard ist aus 'GibHuellkurve ()' und 'GibKlangParameter()' ersichtlich.
//       Die Einstellungen mit 'SetzeHuellkurve' und 'SetzeKlangparameter' haben Einfluss auf den "Ton".
// SpieleNote (tonname string, laenge float64, wartedauer float64)

// Vor.: Das Grafikfenster ist offen.
//       rate ist die Abtastrate, z.B. 11025, 22050 oder 44100.
//       auflösung ist 1 für 8 Bit oder 2 für 16 Bit.
//       kanaele ist 1 für mono oder 2 für stereo.
//       signal gibt die Signalform an: 0: Sinus, 1: Rechteck, 2:Dreieck, 3: Sägezahn
//       p ist die Pulsweite für Rechtecksignale und gibt den Prozentsatz (0<=p<=1) für den HIGH-Teil an.
// Eff.: Die klangparameter sind auf die angegebenen Werte gesetzt.
// SetzeKlangparameter(rate uint32, aufloesung,kanaele,signal uint8, p float64)
func Klack() {
	gfx.SetzeHuellkurve(0.01, 0.1, 0.1, 5)
	gfx.SpieleNote("5D#", 0.03, 0)
	//gfx.SetzeKlangparameter(2, 3, 2, 0, 0)
	//gfx.SpieleSound("C:\\Users\\fussb\\OneDrive\\Arbeitsplatz privat\\bbSt-Inf\\11 SWP\\lwb-SWP\\MiniBillard\\hilf\\White-pool-ball-hitting-a-solid-ball.wav")
}

func Klonk() {
	gfx.SetzeHuellkurve(0.1, 0.2, 0.1, 2)
	gfx.SpieleNote("2D", 0.5, 0)
}

func Bobb() {
	// gfx.SetzeHuellkurve(0.03, 2, 0.5, 0)
	// gfx.SpieleNote("1F", 0.1, 0)
}
