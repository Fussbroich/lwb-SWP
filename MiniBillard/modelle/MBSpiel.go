package modelle

import (
	"time"

	"../hilf"
)

// TODO Ein MiniBillardSpiel ist ...
//
//	New...()
type MiniBillardSpiel interface {
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	Starte()
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	Stoppe()
	Laeuft() bool
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	ZeitlupeAnAus()
	PauseAnAus()
	IstZeitlupe() bool
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	Stosse()
	StossWiederholen()
	Reset()
	IstStillstand() bool
	GibTaschen() []MBTasche
	GibKugeln() []MBKugel
	GibAktiveKugeln() []MBKugel
	NotiereBerührt(MBKugel, MBKugel)
	NotiereEingelocht(MBKugel)
	GibEingelochteKugeln() []MBKugel
	GibSpielkugel() MBKugel
	// TODO Spezifikation
	// ...
	//
	//	Vor.: ...
	//	Eff.: ...
	SetzeKugeln1BallTest() // Testzwecke
	SetzeKugeln3Ball()     // Testzwecke
	SetzeKugeln9Ball()     // Testzwecke
	GibVStoss() hilf.Vec2
	SetzeStossRichtung(hilf.Vec2)
	SetzeStosskraft(float64)
	SetzeSpielzeit(time.Duration)
	SetzeRestzeit(time.Duration) // Testzwecke
	GibRestzeit() time.Duration
	GibTreffer() uint8
	GibStrafpunkte() uint8
	AlleEingelocht() bool
	ReduziereStrafpunkte()
	ErhoeheStrafpunkte() // Testzwecke
	GibGroesse() (float64, float64)
}
