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

func NewMiniBillardSpiel(br uint16) *spiel {
	var breite, höhe float64 = float64(br), float64(br) / 2
	rk := breite / 38 // Kugelradius
	rt, rtm := 2.0*rk, 1.6*rk
	sp := &spiel{breite: breite, höhe: höhe}
	sp.setzeTaschen(
		NewTasche(pos(0, 0), rt),
		NewTasche(pos(0, höhe), rt),
		NewTasche(pos(breite/2, höhe), rtm),
		NewTasche(pos(breite, höhe), rt),
		NewTasche(pos(breite, 0), rt),
		NewTasche(pos(breite/2, 0), rtm))
	pStoß := pos(4*breite/5, höhe/3)
	p1 := pos(3*breite/5, höhe/2)
	p11 := p1.Plus(pos(-2*(rk+1), -(rk + 1)))
	p2 := p1.Plus(pos(-2*(rk+1), (rk + 1)))
	sp.setzeKugeln(
		NewKugel(pStoß, rk, 0),
		NewKugel(p1, rk, 1),
		NewKugel(p11, rk, 11),
		NewKugel(p2, rk, 2))
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
