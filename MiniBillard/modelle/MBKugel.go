package modelle

import "../hilf"

// Eine MBKugel repräsentiert eine Billardkugel in 2 Dimensionen. Sie hat eine Position
// (Mittelpunktskoordinaten), einen Radius, eine Bewegungs-Geschwindigkeit.
// Außerdem kennt sie ihre Nummer (Wert) im Billardspiel.
//
//	NewKugel(pos hilf.Vec2, r float64, wert uint8) erzeugt ein Objekt
type MBKugel interface {
	// Die Bewegungsmethode; wird in der Simulation von einem Billard-Spiel-Modell (MBSpiel)
	// in der Spielschleife aufgerufen.
	//
	//	Vor.: keine
	//	Eff.: Die Kugel bewegt sich voran und ändert bei Kollisionen mit dem Rand des Spieltuches
	//	oder mit einer anderen Kugel ihre Richtung.
	//
	//	Hinweis: Bei der derzeitigen Implementierung wird die eigentliche Bewegungsgeschwindigkeit durch
	//	erhöhen oder vermindern der Aufruf-Frequenz beeinflusst.
	BewegenIn(MiniBillardSpiel)
	// Notifier einer anderen Kugel über die Kollision mit ihr.
	SetzeKollidiertMit(MBKugel)
	// Löscht Nachricht über die Kollision mit anderer Kugel.
	SetzeKollidiertZurueck()
	// Getter, ob in einer Tasche.
	IstEingelocht() bool
	// Getter für die vektorielle (2D) Bewegungsgeschwindigkeit.
	GibV() hilf.Vec2
	// Getter für die vektorielle (2D) Bewegungsgeschwindigkeit.
	SetzeV(hilf.Vec2)
	// Setze die vektorielle (2D) Bewegungsgeschwindigkeit auf Null.
	Stop()
	// Getter für die vektorielle (2D) Position auf dem Tuch.
	GibPos() hilf.Vec2
	// Setter für die vektorielle (2D) Position auf dem Tuch.
	SetzePos(hilf.Vec2)
	// Getter für den in der Simulation gedachten Radius in Pixeln.
	GibRadius() float64
	// Getter für die Nummer (Wert) der Kugel
	GibWert() uint8
	// Getter für eine komplett unabhängige Kopie derselben Kugel
	GibKopie() MBKugel
}
