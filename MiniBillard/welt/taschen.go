package welt

import "../hilf"

type Tasche interface {
	GibPos() hilf.Vec2
	GibRadius() float64
}

type tasche struct {
	pos hilf.Vec2
	r   float64
}

func NewTasche(pos hilf.Vec2, r float64) *tasche {
	return &tasche{
		pos: pos, r: r}
}

func (t *tasche) GibPos() hilf.Vec2 {
	return t.pos
}

func (t *tasche) GibRadius() float64 {
	return t.r
}
