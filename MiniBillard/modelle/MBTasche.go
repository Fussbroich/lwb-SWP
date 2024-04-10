package modelle

import "../hilf"

type MBTasche interface {
	GibPos() hilf.Vec2
	GibRadius() float64
}
