package assets

func klangDateipfad(filename string) string {
	dir := "MiniBillard/assets/soundfiles"
	return assetDateipfad(dir, filename)
}

func MassivePulseDateipfad() string { return klangDateipfad("massivePulseLoop.wav") }

func CoolJazz2641Dateipfad() string { return klangDateipfad("coolJazzLoop2641.wav") }

func BillardPubAmbienceDateipfad() string { return klangDateipfad("billardPubAmbience.wav") }

func CueHitsBallDateipfad() string { return klangDateipfad("cueHitsBall.wav") }

func BallHitsBallDateipfad() string { return klangDateipfad("ballHitsBall.wav") }

func BallInPocketDateipfad() string { return klangDateipfad("ballIntoPocket.wav") }

func BallHitsRailDateipfad() string { return klangDateipfad("ballHitsRail.wav") }
