package modelle

import (
	"time"

	"../hilf"
)

type MiniBillardSpiel interface {
	Starte()
	Stoppe()
	Läuft() bool
	ZeitlupeAnAus()
	IstZeitlupe() bool
	Stoße()
	StoßWiederholen()
	Reset()
	IstStillstand() bool
	GibTaschen() []MBTasche
	GibKugeln() []MBKugel
	GibAktiveKugeln() []MBKugel
	Einlochen(MBKugel)
	GibEingelochteKugeln() []MBKugel
	GibStoßkugel() MBKugel
	GibVStoß() hilf.Vec2
	SetzeStoßRichtung(hilf.Vec2)
	SetzeStoßStärke(float64)
	SetzeRestzeit(time.Duration)
	GibRestzeit() time.Duration
	GibTreffer() uint8
	GibStrafpunkte() uint8
	ReduziereStrafpunkte()
	GibGröße() (float64, float64)
}
