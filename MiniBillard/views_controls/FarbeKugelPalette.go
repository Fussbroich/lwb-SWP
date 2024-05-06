package views_controls

// Eine Kugelpalette bietet 16 verschiedene
// Farben für die Kugeln eines Billardspiels.
// Siehe Implementierung für die verwendeten Farbwerte.
// Kein öffentlicher Konstruktor. Stattdessen kann mit den Paketfunktionen
// StandardKugelPalette()
// EnglishKugelPalette() Eine Palette gewählt werden, die dann global gilt.

type KugelPalette interface {
	// Getter für die Farbe der Kugel mit dem Wert (0 bis 16).
	// Hinweis: Kugel #0 ist die weiße Spielkugel.
	GibFarbe(uint8) Farbe
}
