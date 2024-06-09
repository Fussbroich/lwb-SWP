package klaenge

import (
	"time"

	"../hilf"
)

type klang struct {
	titel  string
	dauer  time.Duration
	autor  string
	player hilf.Routine
	play   func()
}

func (s *klang) Play() {
	if s.play != nil {
		// spielt den Klang einmal ab - bei gfx: nicht mehr stoppbar
		go s.play()
	}
}

func (s *klang) StarteLoop() {
	// Spielt den Klang in Dauerschleife.
	if s.player == nil {
		s.player = hilf.NewRoutine(s.titel, s.play)
	}
	s.player.StarteMitTakt(s.dauer)
}

func (s *klang) Stoppe() {
	// Todo: Ein gfx-Klang läuft dennoch ganz durch und stoppt dann erst.
	if s.player != nil {
		s.player.Stoppe()
	}
}
