package welt

import (
	"../hilf"
)

type MiniBillardSpiel interface {
	GibGröße() (float64, float64)
	GibBanden() []Bande
	GibTaschen() []Tasche
	GibKugeln() []Kugel
	GibStoßkugel() Kugel
	EntferneKugel(Kugel)
	GibBahnDreiecke() [][3]hilf.Vec2
	UpdateObjekte()
	Anstoß(hilf.Vec2)
	IstStillstand() bool
	SetzeVerloren()
	IstVerloren() bool
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
	verloren     bool
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

func (s *spiel) EntferneKugel(k Kugel) {
	if k == s.stoßkugel && len(s.kugeln) > 1 {
		s.verloren = true
	}
	var ksNeu []Kugel
	for _, k2 := range s.kugeln {
		if !(k2 == k) {
			ksNeu = append(ksNeu, k)
		}
	}
	s.kugeln = ksNeu
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

func (s *spiel) UpdateObjekte() {
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

func (s *spiel) IstVerloren() bool {
	return s.verloren
}

func (s *spiel) SetzeVerloren() {
	s.verloren = true
}

func pos(x, y float64) hilf.Vec2 {
	return hilf.V2(x, y)
}

func dr(p1, p2, p3 hilf.Vec2) [3]hilf.Vec2 {
	return [3]hilf.Vec2{p1, p2, p3}
}

func NewStandardSpiel(länge, breite float64) *spiel {
	var (
		raKugel float64 = länge / 40 // Kugelradius
		bn      float64 = länge / 60 // Dicke der Banden
	)
	spiel := NewSpiel(länge, breite)

	pA, pB, pC, pD := pos(bn, bn), pos(bn, breite-bn), pos(länge-bn, breite-bn), pos(länge-bn, bn)

	spiel.setzeBahnform(pA, pB, pC, pD)
	spiel.setzeBahndreiecke(dr(pA, pB, pC), dr(pA, pC, pD))
	spiel.setzeTaschen(
		NewTasche(pA, 1.7*raKugel),
		NewTasche(pos(länge/2, bn), 1.4*raKugel),
		NewTasche(pD, 1.7*raKugel),
		NewTasche(pB, 1.7*raKugel),
		NewTasche(pos(länge/2, breite-bn), 1.4*raKugel),
		NewTasche(pC, 1.7*raKugel))
	spiel.setzeKugeln(
		NewKugel(pos(4*länge/5, breite/2), raKugel, 255, 255, 255),
		NewKugel(pos(2*länge/5+2, breite/2-raKugel-1), raKugel, 255, 0, 0),
		NewKugel(pos(2*länge/5-2, breite/2+raKugel+1), raKugel, 255, 0, 0))
	return spiel
}

func NewStandardSpielNewtonLinie(länge, breite float64) *spiel {
	spiel := NewStandardSpiel(länge, breite)
	var raKugel, bn float64 = länge / 40, länge / 60
	pB, pD := pos(bn, breite-bn), pos(länge-bn, bn)
	linie := hilf.G2(pD, pB.Minus(pD))
	lDiag := pB.Minus(pD).Betrag()
	spiel.setzeKugeln(
		NewKugel(linie.GibPosFür(4*2*raKugel/lDiag), raKugel, 255, 255, 255),
		NewKugel(linie.GibPosFür(8*(2*raKugel+1)/lDiag), raKugel, 255, 0, 0),
		NewKugel(linie.GibPosFür(9*(2*raKugel+1)/lDiag), raKugel, 255, 0, 0),
		NewKugel(linie.GibPosFür(10*(2*raKugel+1)/lDiag), raKugel, 255, 0, 0))
	return spiel
}

func NewNewtonRauteSpiel(länge, breite float64) *spiel {
	spiel := NewSpiel(länge, breite)
	var raKugel, bn float64 = länge / 40, länge / 60
	pA, pB, pC, pD := pos(länge/2, bn), pos(bn, breite/2), pos(länge/2, breite-bn), pos(länge-bn, breite/2)
	spiel.setzeBahnform(pA, pB, pC, pD)
	spiel.setzeBahndreiecke(dr(pA, pB, pC), dr(pA, pC, pD))
	spiel.setzeTaschen(NewTasche(pos(länge/4, breite/2), 1.4*raKugel))
	spiel.setzeKugeln(
		NewKugel(pos(6*länge/8, breite/2+3*raKugel), raKugel, 255, 255, 255),
		NewKugel(pos(länge/2, breite/2+3*raKugel), raKugel, 255, 0, 0),
		NewKugel(pos(länge/2+2*(1+raKugel), breite/2+3*raKugel), raKugel, 255, 0, 0),
		NewKugel(pos(länge/2+4*(1+raKugel), breite/2+3*raKugel), raKugel, 255, 0, 0))
	return spiel
}

func NewLBahnSpiel(länge float64) *spiel {
	spiel := NewSpiel(länge, länge)
	var t, raKugel float64 = länge / 5, länge / 40
	pA, pB, pC, pD, pE, pF, pG, pH :=
		pos(0, 0), pos(0, 2*t), pos(2*t, 2*t), pos(3*t, 3*t), pos(3*t, 5*t), pos(5*t, 5*t), pos(5*t, 2*t), pos(3*t, 0)
	spiel.setzeBahnform(pA, pB, pC, pD, pE, pF, pG, pH)
	spiel.setzeBahndreiecke(dr(pA, pB, pC), dr(pC, pD, pH), dr(pA, pC, pH), dr(pH, pD, pG), dr(pG, pD, pF), dr(pD, pE, pF))
	spiel.setzeTaschen(NewTasche(pos(t, t), 1.4*raKugel))
	spiel.setzeKugeln(
		NewKugel(pos(4*t, 4*t), raKugel, 255, 255, 255),
		NewKugel(pos(4*t, 3*t), raKugel, 255, 0, 0))
	return spiel
}

func NewHexaBahnSpiel(länge float64) *spiel {
	spiel := NewSpiel(länge, länge)
	var t, raKugel float64 = länge / 16, länge / 40
	pA, pB, pC, pD, pE, pF :=
		pos(4*t, 1*t), pos(0, 8*t), pos(4*t, 15*t), pos(12*t, 15*t), pos(16*t, 8*t), pos(12*t, 1*t)
	spiel.setzeBahnform(pA, pB, pC, pD, pE, pF)
	spiel.setzeBahndreiecke(dr(pA, pB, pF), dr(pB, pC, pF), dr(pC, pE, pF), dr(pC, pD, pE))
	spiel.setzeTaschen(NewTasche(pos(8*t, 8*t), 1.4*raKugel))
	spiel.setzeKugeln(
		NewKugel(pos(6*t, 13*t), raKugel, 255, 255, 255),
		NewKugel(pos(11*t, 7*t), raKugel, 255, 0, 0))
	return spiel
}
