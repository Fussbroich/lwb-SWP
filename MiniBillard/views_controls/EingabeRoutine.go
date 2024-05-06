package views_controls

import "../hilf"

// Eine spezialisierte Routine für die
// Abfrage der blockierenden Eingabe-Funktionen
// von gfx in einer Endlosschleife.
// Siehe die gfx-Dokumentation für Infos zur Benutzung der
// Eingabe-Funktionen.
//
// Konstruktoren entsprechend dem gewünschten Eingabekanal:
// NewMausRoutine(f func(t uint8, s int8, x uint16, y uint16))
// bzw. NewTastenRoutine(f func(uint16, uint8, uint16))
type EingabeRoutine interface {
	hilf.Routine
}
