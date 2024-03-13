package welt

import (
	"../hilf"
)

type MiniBillardSpiel interface {
	Update()
	Stoße(hilf.Vec2)
	StoßWiederholen()
	Reset()
	IstStillstand() bool
	GibTaschen() []Tasche
	GibKugeln() []Kugel
	GibAktiveKugeln() []Kugel
	GibEingelochteKugeln() []Kugel
	GibStoßkugel() Kugel
	GibStößeBisher() uint8
	GibStrafpunkte() uint8
	GibGröße() (float64, float64)
}

type spiel struct {
	breite       float64
	höhe         float64
	rk           float64
	kugeln       []Kugel
	origKugeln   []Kugel
	vorigeKugeln []Kugel
	stoßkugel    Kugel
	taschen      []Tasche
	stößeBisher  uint8
	strafPunkte  uint8
	stillstand   bool
}

func pos(x, y float64) hilf.Vec2 {
	return hilf.V2(x, y)
}

func (s *spiel) setzeTaschen(t ...Tasche) {
	s.taschen = []Tasche{}
	s.taschen = append(s.taschen, t...)
}

func (s *spiel) setzeKugeln(k ...Kugel) {
	s.kugeln = []Kugel{}
	for _, k := range k {
		k.Stop()
		s.kugeln = append(s.kugeln, k.GibKopie())
	}
	s.stoßkugel = s.kugeln[0]

	// sichere den Zustand vor dem Anstoß
	s.origKugeln = []Kugel{}
	for _, k := range s.kugeln {
		s.origKugeln = append(s.origKugeln, k.GibKopie())
	}
}

func (sp *spiel) SetzeKugeln9Ball() {
	sp.setzeKugeln(sp.kugelSatz9Ball()...)
}

func (s *spiel) kugelSatz3er() []Kugel {
	pStoß := pos(4*s.breite/5, s.höhe/3)
	p1 := pos(3*s.breite/5, s.höhe/2)
	dx, dy := 0.866*(2*s.rk+2), 0.5*(2*s.rk+2)
	p2 := p1.Plus(pos(-dx, -dy))
	p3 := p1.Plus(pos(-dx, dy))
	return []Kugel{NewKugel(pStoß, s.rk, 0),
		NewKugel(p1, s.rk, 1),
		NewKugel(p2, s.rk, 2),
		NewKugel(p3, s.rk, 3)}
}

func (s *spiel) kugelSatz9Ball() []Kugel {
	pStoß := pos(4*s.breite/5, s.höhe/2)
	dx, dy := 0.866*(2*s.rk+1), 0.5*(2*s.rk+1)
	p1 := pos(3*s.breite/5, s.höhe/2)
	p5 := p1.Plus(pos(-dx, -dy))
	p3 := p1.Plus(pos(-dx, dy))

	p2 := p1.Plus(pos(-2*dx, -2*dy))
	p9 := p1.Plus(pos(-2*dx, 0))
	p8 := p1.Plus(pos(-2*dx, 2*dy))

	p7 := p1.Plus(pos(-3*dx, -dy))
	p6 := p1.Plus(pos(-3*dx, dy))
	p4 := p1.Plus(pos(-4*dx, 0))
	return []Kugel{NewKugel(pStoß, s.rk, 0),
		NewKugel(p1, s.rk, 1),
		NewKugel(p2, s.rk, 2),
		NewKugel(p3, s.rk, 3),
		NewKugel(p4, s.rk, 4),
		NewKugel(p5, s.rk, 5),
		NewKugel(p6, s.rk, 6),
		NewKugel(p7, s.rk, 7),
		NewKugel(p8, s.rk, 8),
		NewKugel(p9, s.rk, 9),
	}
}

func NewMiniBillardSpiel(br uint16) *spiel {
	var breite, höhe float64 = float64(br), float64(br) / 2
	sp := &spiel{breite: breite, höhe: höhe, rk: breite / 38}
	rt, rtm := 2.0*sp.rk, 1.6*sp.rk
	sp.setzeTaschen(
		NewTasche(pos(0, 0), rt),
		NewTasche(pos(0, höhe), rt),
		NewTasche(pos(breite/2, höhe), rtm),
		NewTasche(pos(breite, höhe), rt),
		NewTasche(pos(breite, 0), rt),
		NewTasche(pos(breite/2, 0), rtm))
	sp.setzeKugeln(sp.kugelSatz9Ball()...)
	return sp
}

func (s *spiel) Reset() {
	s.kugeln = []Kugel{}
	for _, k := range s.origKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie())
	}
	s.stoßkugel = s.kugeln[0]
	s.stößeBisher = 0
	s.strafPunkte = 0
	s.stillstand = true
}

func (s *spiel) StoßWiederholen() {
	// stelle den Zustand vor dem letzten Stoß wieder her
	s.kugeln = []Kugel{}
	for _, k := range s.vorigeKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie())
	}
	s.stoßkugel = s.kugeln[0]
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

func (s *spiel) Stoße(v hilf.Vec2) {
	// sichere den Zustand vor dem Stoß
	s.vorigeKugeln = []Kugel{}
	for _, k := range s.kugeln {
		s.vorigeKugeln = append(s.vorigeKugeln, k.GibKopie())
	}
	// stoße
	s.stoßkugel.SetzeV(v)
	s.stößeBisher++
	s.stillstand = false
}

func (s *spiel) GibGröße() (float64, float64) {
	return s.breite, s.höhe
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

func (s *spiel) GibEingelochteKugeln() []Kugel {
	ks := []Kugel{}
	for _, k := range s.kugeln {
		if !k.IstEingelocht() {
			continue
		}
		ks = append(ks, k)
	}
	return ks
}

func (s *spiel) GibStoßkugel() Kugel {
	return s.stoßkugel
}

func (s *spiel) IstStillstand() bool {
	return s.stillstand
}

func (s *spiel) GibStößeBisher() uint8 {
	return s.stößeBisher
}

func (s *spiel) GibStrafpunkte() uint8 {
	return s.strafPunkte
}
