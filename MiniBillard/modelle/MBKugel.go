package modelle

import "../hilf"

// TODO Eine MBKugel ist ...
//
//	New...()
type MBKugel interface {
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
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
