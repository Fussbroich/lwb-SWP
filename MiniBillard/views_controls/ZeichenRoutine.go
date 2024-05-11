package views_controls

import "../hilf"

// Eine spezialisierte Routine, die die Darstellungs-Methode einer App aufruft
// und damit die Bildwiederholung steuert.
// Im "unmittelbaren Modus" wird die App in einem regelmäßigen Takt in ein
// einziges Fenster gezeichnet.
//
//	Vor.: Das Grafikpaket gfx muss im GOPATH installiert sein.
//
//	NewZeichenRoutine(App) erzeugt ein Objekt.
type ZeichenRoutine interface {
	hilf.Routine
}
