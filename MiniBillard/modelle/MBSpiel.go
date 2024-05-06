package modelle

import (
	"time"

	"../hilf"
)

// Ein MiniBillardSpiel besteht aus einem rechteckigen Tuch (Spieltisch)
// mit am Rand verteilten kreisförmigen Taschen. Ein Spiel bietet eine sehr einfache
// 2D-Simulation von bewegten Billardkugeln gleicher Masse und gleicher Größe. Bewegt
// sich eine Kugel über eine Tasche, so nimmt sie nicht weiter an der Simulation teil
// und "verschwindet". Das Spiel enthält beliebig viele Kugeln (2D-Objekte vom typ MBKugel),
// die sich innerhalb des Tuchbereiches frei bewegen können.
// Das Spiel kontrolliert die Bewegungsmethode jeder Kugel in einer Spielschleife,
// die mit einer kontrollierten Frequenz aufgerufen wird. So wird eine bestimmte
// Zeitordnung und ein bestimmtes Geschwindigkeitsverhalten der Kugeln simuliert.
// Jede Kugel prüft, ob sie andere Kugeln oder den Rand des Tuches berührt und
// ändert entsprechend ihre Bewegungsrichtung (elastischer Stoß).
//
// Konstruktoren: newPoolSpiel(br, hö uint16), bzw. NewMiniXBallSpiel(br, hö uint16) erzeugen
// ein Objekt ohne bzw. mit X Kugeln und einer weißen Spielkugel.
type MiniBillardSpiel interface {
	// Startet die Spielschleife mit einer geregelten Frequenz.
	//
	//	Vor.: keine
	//	Eff.: Jede enthaltene Kugel wird in konstanten Abständen bewegt;
	//	die Spielregeln, sowie der Kugelstillstand werden laufen überwacht.
	Starte()
	// Stoppt die Spielschleife (Pause-Modus) - lässt sich wieder starten.
	//	Vor.: Spielschleife läuft
	//	Eff.: Spielschleife läuft nicht
	Stoppe()
	Laeuft() bool
	// Verringert den Betrag der Kugelbewegungen um einen festen Quotienten
	//	Vor.: keine
	//	Eff.: Kugeln bewegen sich langsamer
	ZeitlupeAnAus()
	PauseAnAus()
	IstZeitlupe() bool
	// Nur wenn alle Kugeln stillstehen, kann die weiße Spielkugel angestossen werden.
	//
	//	Vor.: Alle Kugeln stehen still.
	//	Eff.: Die Spielkugel bewegt sich mit einer vorher gesetzten Geschwindigkeit und Richtung.
	Stosse()
	// Alle Kugeln werden in den Zustand vor dem vergangenen Stoß versetzt.
	//
	//	Eff.: Kugeln stehen still auf den Positionen vor dem vergangenen Stoß.
	StossWiederholen()
	// Das ganze Spiel wird neu begonnen.
	//
	//	Eff.: Alle Kugeln wie zu Beginn des Spiels, Spielkugel in der Küche an einer
	//	zufälligen Position gesetzt. Zeit beginnt neu. Punkte stehen bei 0.
	Reset()
	IstStillstand() bool
	GibTaschen() []MBTasche
	// Getter für alle Kugeln, die im Spiel sind.
	GibKugeln() []MBKugel
	// Getter nur für diejenigen Kugeln, die noch nicht eingelocht sind.
	GibAktiveKugeln() []MBKugel
	// Notifier für die Berührung zwischen zwei Kugeln.
	NotiereBerührt(MBKugel, MBKugel)
	// Notifier für das Einlochen einer Kugel.
	NotiereEingelocht(MBKugel)
	// Getter nur für diejenigen Kugeln, die bereits eingelocht sind in der Reihenfolge des Einlochens.
	GibEingelochteKugeln() []MBKugel
	// Getter nur für die weiße Spielkugel.
	GibSpielkugel() MBKugel
	// Setter für bestimmte Kugelkonstellationen
	SetzeKugeln1BallTest() // Testzwecke
	SetzeKugeln3Ball()     // Testzwecke
	SetzeKugeln9Ball()     // Testzwecke
	// Getter für die aktuell gesetzte vektorielle (2D) Stoßgeschwindigkeit (Betrag und Richtung).
	GibVStoss() hilf.Vec2
	// Setter für die vektorielle (2D) Stoßrichtung.
	SetzeStossRichtung(hilf.Vec2)
	// Setter für den Betrag der Stoßkraft.
	SetzeStosskraft(float64)
	// Setter für die gesamte Zeitdauer zum "Spiel gegen die Zeit".
	SetzeSpielzeit(time.Duration)
	// Setter für die künstliche Beeinflussung der verbleibenden Spieldauer.
	SetzeRestzeit(time.Duration) // Testzwecke
	// Getter für die verbleibende Spieldauer.
	GibRestzeit() time.Duration
	// Getter für die Anzahl der in dem Spiel bereits versenkten (eingelochten) Kugeln
	GibTreffer() uint8
	// Getter für die Anzahl der Fouls in diesem Spiel.
	GibStrafpunkte() uint8
	// Prüft, ob alle Kugeln, bis auf die weiße Spielgkugel, bereits eingelocht sind.
	AlleEingelocht() bool
	// Reduziere künstlich die Anzahl der Fouls in diesem Spiel um 1.
	ReduziereStrafpunkte()
	// Erhöhe künstlich die Anzahl der Fouls in diesem Spiel um 1.
	ErhoeheStrafpunkte() // Testzwecke
	// Getter für die gedachte Pixelgröße der Simulation.
	GibGroesse() (float64, float64)
}
