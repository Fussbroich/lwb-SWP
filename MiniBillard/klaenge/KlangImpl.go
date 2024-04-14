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
		// spielt den Klang einmal ab - nicht mehr stoppbar
		go s.play()
	}
}

func (s *klang) StarteLoop() {
	// Spielt den Klang in Dauerschleife.
	if s.player == nil {
		s.player = hilf.NewRoutine(s.titel, s.play)
	}
	s.player.StarteLoop(s.dauer)
}

func (s *klang) Stoppe() {
	// Todo: Der Klang läuft dennoch ganz durch und stoppt dann erst.
	if s.player != nil {
		s.player.Stoppe()
	}
}
