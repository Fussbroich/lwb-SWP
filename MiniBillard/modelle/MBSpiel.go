package modelle

import (
	"time"

	"../hilf"
)

type MiniBillardSpiel interface {
	Starte()
	Stoppe()
	Laeuft() bool
	ZeitlupeAnAus()
	PauseAnAus()
	IstZeitlupe() bool
	Stosse()
	StossWiederholen()
	Reset()
	IstStillstand() bool
	GibTaschen() []MBTasche
	GibKugeln() []MBKugel
	GibAktiveKugeln() []MBKugel
	Einlochen(MBKugel)
	GibEingelochteKugeln() []MBKugel
	GibSpielkugel() MBKugel
	GibVStoss() hilf.Vec2
	SetzeStossRichtung(hilf.Vec2)
	SetzeStosskraft(float64)
	SetzeRestzeit(time.Duration)
	GibRestzeit() time.Duration
	GibTreffer() uint8
	GibStrafpunkte() uint8
	ReduziereStrafpunkte()
	GibGroesse() (float64, float64)
}
