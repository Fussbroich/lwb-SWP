package welt

import (
	"../hilf"
	"../klaenge"
)

type Kugel interface {
	BewegenIn(MiniBillardSpiel)
	SetzeKollidiertMit(Kugel)
	IstEingelocht() bool
	GibV() hilf.Vec2
	SetzeV(hilf.Vec2)
	Stop()
	GibPos() hilf.Vec2
	SetzePos(hilf.Vec2)
	GibRadius() float64
	GibFarbe() (uint8, uint8, uint8)
	GibKopie() Kugel
}

type kugel struct {
	pos, v     hilf.Vec2
	r          float64
	cR         uint8
	cG         uint8
	cB         uint8
	istKollMit Kugel
	eingelocht bool
}

func NewKugel(pos hilf.Vec2, r float64, farbeR, farbeG, farbeB uint8) *kugel {
	return &kugel{
		pos: pos,
		r:   r, cR: farbeR, cG: farbeG, cB: farbeB}
}

func (k *kugel) GibKopie() Kugel {
	return &kugel{
		pos: k.pos, v: hilf.V2(0, 0),
		r:  k.r,
		cR: k.cR, cG: k.cG, cB: k.cB}
}

func (k *kugel) BewegenIn(s MiniBillardSpiel) {
	if k.eingelocht {
		return
	}
	// prüfe Kollisionen
	for _, k2 := range s.GibKugeln() {
		if k != k2 {
			k.prüfeKugelKollision(k2)
		}
	}
	// Prüfe Berührung mit der Bande.
	k.prüfeBandenKollision(s.GibGröße())
	// Bewege Kugel einen Tick weiter.
	k.pos = k.pos.Plus(k.v)
	vabs := k.v.Betrag()
	// Bremse die Kugel etwas ab.
	if vabs > 0.15 {
		k.v = k.v.Mal(1 - 0.02/vabs)
	} else {
		k.v = hilf.V2(0, 0)
	}
	k.istKollMit = nil
	// Prüfe, ob Kugel eingelocht wurde.
	for _, t := range s.GibTaschen() {
		if t.GibPos().Minus(k.GibPos()).Betrag() < t.GibRadius() {
			klaenge.BallInPocketSound()
			k.eingelocht = true
			k.SetzeV(hilf.V2(0, 0))
			break
		}
	}
}

func (k *kugel) IstEingelocht() bool {
	return k.eingelocht
}

func (k *kugel) prüfeBandenKollision(länge, breite float64) {
	if k.IstEingelocht() {
		return
	}
	// Kugel vorher
	vx, vy := k.v.X(), k.v.Y()
	xK, yK := k.pos.X(), k.pos.Y()

	var hit bool
	// reflektiere die Kugel
	if xK+vx < k.r {
		vx *= -1
		hit = true
	}
	if xK+vx > länge-k.r {
		vx *= -1
		hit = true
	}
	if yK+vy < k.r {
		vy *= -1
		hit = true
	}
	if yK+vy > breite-k.r {
		vy *= -1
		hit = true
	}

	if hit {
		klaenge.BallHitsRailSound()
		k.v = hilf.V2(vx, vy)
	}
}

func (k1 *kugel) prüfeKugelKollision(k2 Kugel) {
	if k1.istKollMit == k2 {
		return
	}
	if k1.IstEingelocht() || k2.IstEingelocht() {
		return
	}
	v1 := k1.GibV()
	v2 := k2.GibV()
	dist := k2.GibPos().Plus(v2).Minus(k1.pos.Plus(v1))
	// Kugeln werden sich gar nicht berühren.
	if dist.Betrag() > (k1.r + k2.GibRadius()) {
		return
	}
	// Kugeln überlappen!
	//	/* TODO: was kann man tun, damit sich die Kugeln nicht gegenseitig einfangen?
	//if dist.Betrag() < (k1.r + k2.GibRadius()) {
	//println("Überlappung", k1.r+k2.GibRadius()-dist.Betrag())
	/*
		for dist.Betrag() <= (k1.r + k2.GibRadius()) {
		// treibe Kugeln auseinander
		}
		return
	*/
	//}
	//	*/

	// Kugeln berühren sich
	klaenge.BallHitsBallSound()
	// die Stoßnormale geht durch die Mittelpunkte der Kugeln
	n12 := dist.Normiert()
	// Zerlege Geschwindigkeiten in eine parallele und eine orthogonale Komponente
	v1p := v1.ProjiziertAuf(n12)
	v1o := v1.Minus(v1p)
	v2p := v2.ProjiziertAuf(n12)
	v2o := v2.Minus(v2p)
	// Tausche Geschwindigkeiten aus
	u1 := v2p.Plus(v1o)
	u2 := v1p.Plus(v2o)
	k1.SetzeV(u1)
	k1.istKollMit = k2
	k2.SetzeV(u2)
	k2.SetzeKollidiertMit(k1)
}

func (k1 *kugel) SetzeKollidiertMit(k2 Kugel) {
	k1.istKollMit = k2
}

func (k *kugel) GibV() hilf.Vec2 {
	return k.v
}

func (k *kugel) SetzeV(v hilf.Vec2) {
	k.v = v
}

func (k *kugel) Stop() {
	k.v = hilf.V2(0, 0)
}

func (k *kugel) GibPos() hilf.Vec2 {
	return k.pos
}

func (k *kugel) SetzePos(pos hilf.Vec2) {
	k.pos = pos
}

func (k *kugel) GibRadius() float64 {
	return k.r
}

func (k *kugel) GibFarbe() (uint8, uint8, uint8) {
	return k.cR, k.cG, k.cB
}
