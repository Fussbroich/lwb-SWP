package welt

import (
	"../hilf"
)

type MiniBillardSpiel interface {
	Update()
	Anstoß(hilf.Vec2)
	StoßWiederholen()
	IstStillstand() bool
	GibTaschen() []Tasche
	GibKugeln() []Kugel
	GibAktiveKugeln() []Kugel
	GibStoßkugel() Kugel
	GibStößeBisher() uint8
	GibStrafpunkte() uint8
	GibGröße() (float64, float64)
}

type spiel struct {
	länge       float64
	breite      float64
	kugeln      []Kugel
	alteKugeln  []Kugel
	stoßkugel   Kugel
	taschen     []Tasche
	stößeBisher uint8
	strafPunkte uint8
	stillstand  bool
}

func pos(x, y float64) hilf.Vec2 {
	return hilf.V2(x, y)
}

func (s *spiel) setzeTaschen(t ...Tasche) {
	s.taschen = append(s.taschen, t...)
}

func (s *spiel) setzeKugeln(k ...Kugel) {
	s.kugeln = []Kugel{}
	for _, k := range k {
		k.Stop()
		s.kugeln = append(s.kugeln, k)
	}
	s.stoßkugel = s.kugeln[0]
}

func NewSpiel(länge, breite float64) *spiel {
	return &spiel{länge: länge, breite: breite}
}

func New3BallStandardSpiel() *spiel {
	var länge, breite float64 = 1000, 500
	rk := länge / 40 // Kugelradius
	spiel := NewSpiel(länge, breite)
	rt, rtm := 2.0*rk, 1.6*rk
	spiel.setzeTaschen(
		NewTasche(pos(0, 0), rt),
		NewTasche(pos(0, breite), rt),
		NewTasche(pos(länge/2, breite), rtm),
		NewTasche(pos(länge, breite), rt),
		NewTasche(pos(länge, 0), rt),
		NewTasche(pos(länge/2, 0), rtm))
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

func (s *spiel) GibGröße() (float64, float64) {
	return s.länge, s.breite
}

func (s *spiel) GibTaschen() []Tasche {
	return s.taschen
}

func (s *spiel) GibKugeln() []Kugel {
	return s.kugeln
}

func (s *spiel) GibAktiveKugeln() []Kugel {
	ks := []Kugel{}
	for _, k := range s.kugeln {
		if k.IstEingelocht() {
			continue
		}
		ks = append(ks, k)
	}
	return ks
}

func (s *spiel) GibStoßkugel() Kugel {
	return s.stoßkugel
}

func (s *spiel) Anstoß(v hilf.Vec2) {
	// sichere den Zustand vor dem Anstoß
	s.alteKugeln = []Kugel{}
	for _, k := range s.kugeln {
		kNeu := k.GibKopie()
		k.Stop()
		s.alteKugeln = append(s.alteKugeln, kNeu)
	}
	s.stoßkugel.SetzeV(v)
	s.stößeBisher++
	s.stillstand = false
}

func (s *spiel) GibStößeBisher() uint8 {
	return s.stößeBisher
}

func (s *spiel) GibStrafpunkte() uint8 {
	return s.strafPunkte
}

func (s *spiel) StoßWiederholen() {
	s.setzeKugeln(s.alteKugeln...)
	s.strafPunkte++
	s.stillstand = true
}

func (s *spiel) Update() {
	still := true
	for _, k := range s.kugeln {
		k.BewegenIn(s)
		//prüfe Stillstand
		if !k.GibV().IstNull() {
			still = false
		}
	}
	s.stillstand = still
}

func (s *spiel) IstStillstand() bool {
	return s.stillstand
}
