package modelle

import "../hilf"

// Eine MBKugel repräsentiert eine Billardkugel in 2 Dimensionen. Sie hat eine Position
// (Mittelpunktskoordinaten), einen Radius, eine Bewegungs-Geschwindigkeit.
// Außerdem kennt sie ihre Nummer (Wert) im Billardspiel.
//
//	NewKugel(pos hilf.Vec2, r float64, wert uint8) erzeugt ein Objekt
type MBKugel interface {
	// TODO Die Bewegungsmethode; wird in der Simulation von einem Billard-Spiel-Modell (MBSpiel)
	// in einem Loop aufgerufen.
	//
	//	Vor.: ...
	//	Eff.: ...
	//	Hinweis: Bei der derzeitigen Implementierung wird die eigentliche Bewegungsgeschwindigkeit durch
	//	erhöhen oder vermindern der Aufruf-Frequenz beeinflusst.
	BewegenIn(MiniBillardSpiel)
	SetzeKollidiertMit(MBKugel)
	SetzeKollidiertZurueck()
	IstEingelocht() bool
	GibV() hilf.Vec2
	SetzeV(hilf.Vec2)
	Stop()
	GibPos() hilf.Vec2
	SetzePos(hilf.Vec2)
	GibRadius() float64
	GibWert() uint8
	GibKopie() MBKugel
}
