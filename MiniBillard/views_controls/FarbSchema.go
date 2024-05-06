package views_controls

// Ein komplettes Schema an Farben für eine ganze App.
// Die Farben eines bestimmten Schemas werden über IDs
// zugegriffen, die in der App-Konstruktion verwendbar sind.
// Siehe Implementierung für verfügbare Farben und IDs.
type FarbSchema interface {
	// Getter für eine bestimmte Farbe aus einem Schema.
	// Hinweis: Schema-Farben haben IDs. Wird eine unbekannte
	// Farbe angefordert, so gibt die Methode rot zurück. Für die bekannten
	// Farben der hier vorinstallierten Schemata gibt es die Konstanten F...
	GibFarbe(FarbID) Farbe
}
