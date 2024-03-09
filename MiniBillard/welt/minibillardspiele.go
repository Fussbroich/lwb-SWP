package welt

import (
	"../hilf"
)

type MiniBillardSpiel interface {
	BewegeKugeln()
	Anstoß(hilf.Vec2)
	IstStillstand() bool
	GibGröße() (float64, float64)
	GibBanden() []Bande
	GibTaschen() []Tasche
	GibKugeln() []Kugel
	GibStoßkugel() Kugel
	GibBahnDreiecke() [][3]hilf.Vec2
}

type spiel struct {
	länge        float64
	breite       float64
	banden       []Bande
	bahndreiecke [][3]hilf.Vec2
	kugeln       []Kugel
	stoßkugel    Kugel
	taschen      []Tasche
	stillstand   bool
}

func NewSpiel(länge, breite float64) *spiel {
	return &spiel{länge: länge, breite: breite}
}

func (s *spiel) GibGröße() (float64, float64) {
	return s.länge, s.breite
}

func (s *spiel) setzeBahnform(ps ...hilf.Vec2) {
	von := ps[0]
	for i := 1; i < len(ps); i++ {
		s.banden = append(s.banden, NewBande(von, ps[i]))
		von = ps[i]
	}
	s.banden = append(s.banden, NewBande(von, ps[0]))
}

func (s *spiel) setzeBahndreiecke(ds ...[3]hilf.Vec2) {
	s.bahndreiecke = append(s.bahndreiecke, ds...)
}

func (s *spiel) GibBanden() []Bande {
	return s.banden
}

func (s *spiel) setzeTaschen(t ...Tasche) {
	s.taschen = append(s.taschen, t...)
}

func (s *spiel) GibTaschen() []Tasche {
	return s.taschen
}

func (s *spiel) setzeKugeln(k ...Kugel) {
	s.kugeln = k
	s.stoßkugel = k[0]
}

func (s *spiel) GibKugeln() []Kugel {
	return s.kugeln
}

func (s *spiel) GibStoßkugel() Kugel {
	return s.stoßkugel
}

func (s *spiel) Anstoß(v hilf.Vec2) {
	s.stoßkugel.SetzeV(v)
	s.stillstand = false
}

func (s *spiel) GibBahnDreiecke() [][3]hilf.Vec2 {
	return s.bahndreiecke
}

func (s *spiel) BewegeKugeln() {
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
