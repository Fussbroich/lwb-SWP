package welt

import "../hilf"

func pos(x, y float64) hilf.Vec2 {
	return hilf.V2(x, y)
}

func New3BallStandardSpiel(länge float64) *spiel {
	breite := länge / 2
	rk := länge / 40 // Kugelradius
	spiel := NewSpiel(länge, breite)
	rt, rtm := 2.0*rk, 1.6*rk
	spiel.setzeTaschen(
		NewTasche(pos(0, 0), rt),
		NewTasche(pos(länge/2, 0), rtm),
		NewTasche(pos(länge, 0), rt),
		NewTasche(pos(0, breite), rt),
		NewTasche(pos(länge/2, breite), rtm),
		NewTasche(pos(länge, breite), rt))
	pWeiß := pos(4*länge/5, breite/3)
	pGelb := pos(3*länge/5, breite/2)
	pRot := pGelb.Plus(pos(-2*(rk+1), -(rk + 1)))
	pBlau := pGelb.Plus(pos(-2*(rk+1), (rk + 1)))
	spiel.setzeKugeln(
		NewKugel(pWeiß, rk, 255, 255, 255),
		NewKugel(pGelb, rk, 240, 20, 50),
		NewKugel(pRot, rk, 255, 215, 0),
		NewKugel(pBlau, rk, 70, 140, 250))
	return spiel
}
