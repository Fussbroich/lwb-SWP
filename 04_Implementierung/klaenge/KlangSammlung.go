package klaenge

import (
	"brainpool/assets"
	"brainpool/gfx"
	"time"
)

func MassivePulseSound() *klang {
	daten := assets.MassivePulseDaten()
	return &klang{
		titel: "Massive Pulse",
		dauer: 16 * time.Second,
		autor: "unknown",
		play: func() {
			for !gfx.FensterOffen() {
				time.Sleep(100 * time.Millisecond)
			}
			gfx.SpieleSoundDaten(daten, "massivePulse")
		}}
}

func CoolJazz2641SOUND() *klang {
	daten := assets.CoolJazz2641Daten()
	return &klang{
		titel: "Cool Jazz 2641",
		dauer: 2*time.Minute + 8*time.Second,
		autor: "Julius H. (pixabay)",
		play: func() {
			for !gfx.FensterOffen() {
				time.Sleep(100 * time.Millisecond)
			}
			gfx.SpieleSoundDaten(daten, "coolJazz2641")
		}}
}

func BillardPubAmbienceSOUND() *klang {
	daten := assets.BillardPubAmbienceDaten()
	return &klang{
		titel: "Billard Pub Ambience",
		dauer: time.Minute + 13*time.Second,
		autor: "unknown (directory.audio)",
		play: func() {
			for !gfx.FensterOffen() {
				time.Sleep(100 * time.Millisecond)
			}
			gfx.SpieleSoundDaten(daten, "billardPubAmbience")
		}}
}

func CueHitsBallSound() *klang {
	daten := assets.CueHitsBallDaten()
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() { return }
			gfx.SpieleSoundDaten(daten, "cueHitsBall")
		}}
}

func BallHitsBallSound() *klang {
	daten := assets.BallHitsBallDaten()
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() { return }
			gfx.SpieleSoundDaten(daten, "ballHitsBall")
		}}
}

func BallInPocketSound() *klang {
	daten := assets.BallInPocketDaten()
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() { return }
			gfx.SpieleSoundDaten(daten, "ballInPocket")
		}}
}

func BallHitsRailSound() *klang {
	daten := assets.BallHitsRailDaten()
	return &klang{
		dauer: 300 * time.Millisecond,
		autor: "freesman (directory.audio)",
		play: func() {
			if !gfx.FensterOffen() { return }
			gfx.SpieleSoundDaten(daten, "ballHitsRail")
		}}
}
