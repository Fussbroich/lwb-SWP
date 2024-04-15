package modelle

import "time"

type Countdown interface {
	GibRestzeit() time.Duration
	Setze(d time.Duration)
	ZieheAb(d time.Duration)
	IstAbgelaufen() bool
	Halt()
	Weiter()
}
