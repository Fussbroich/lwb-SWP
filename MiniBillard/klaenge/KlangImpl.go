package klaenge

import (
	"time"

	"../hilf"
)

type klang struct {
	titel  string
	dauer  time.Duration
	autor  string
	player hilf.Prozess
	play   func()
}

func (s *klang) Play() {
	go s.play()
}

func (s *klang) StarteLoop() {
	if s.player == nil {
		s.player = hilf.NewProzess(s.titel, s.play)
	}
	s.player.StarteLoop(s.dauer)
}

func (s *klang) Stoppe() {
	if s.player != nil {
		s.player.Stoppe()
	}
}
