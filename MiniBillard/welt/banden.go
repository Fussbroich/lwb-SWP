package welt

import "../hilf"

type Bande interface {
	GibVon() hilf.Vec2
	GibNach() hilf.Vec2
}

type bande struct {
	pV hilf.Vec2
	pN hilf.Vec2
}

func NewBande(pV, pN hilf.Vec2) *bande {
	return &bande{
		pV: pV, pN: pN}
}

func (b *bande) GibVon() hilf.Vec2 {
	return b.pV
}

func (b *bande) GibNach() hilf.Vec2 {
	return b.pN
}
