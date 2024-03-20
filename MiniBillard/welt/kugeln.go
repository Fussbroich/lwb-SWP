package welt

import (
	"../hilf"
	"../klaenge"
)

type Kugel interface {
	BewegenIn(MiniBillardSpiel)
	SetzeKollidiertMit(Kugel)
	SetzeKollidiertZurück()
	IstEingelocht() bool
	GibV() hilf.Vec2
	SetzeV(hilf.Vec2)
	Stop()
	GibPos() hilf.Vec2
	SetzePos(hilf.Vec2)
	GibRadius() float64
	GibWert() uint8
	GibKopie() Kugel
}

type kugel struct {
	pos, v     hilf.Vec2
	r          float64
	wert       uint8
	istKollMit Kugel
	eingelocht bool
}

func NewKugel(pos hilf.Vec2, r float64, wert uint8) *kugel {
	return &kugel{
		pos:  pos,
		r:    r,
		wert: wert}
}

func (k *kugel) GibKopie() Kugel {
	return &kugel{
		pos: k.pos, v: hilf.V2(0, 0),
		r:    k.r,
		wert: k.wert}
}

func (k *kugel) BewegenIn(s MiniBillardSpiel) {
	if k.eingelocht {
		return
	}
	// prüfe Kollisionen
	for _, k2 := range s.GibKugeln() {
		if (k != k2) && !k2.IstEingelocht() {
			k.prüfeKugelKollision(k2)
		}
	}
	// setze kollidierte zurück
	k.istKollMit = nil
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
	// Prüfe, ob Kugel eingelocht wurde.
	for _, t := range s.GibTaschen() {
		if t.GibPos().Minus(k.GibPos()).Betrag() < t.GibRadius() {
			klaenge.BallInPocketSound().Play()
			k.eingelocht = true
			s.Einlochen(k)
			k.SetzeV(hilf.V2(0, 0))
			break
		}
	}
}

func (k *kugel) IstEingelocht() bool { return k.eingelocht }

func (k *kugel) prüfeBandenKollision(länge, breite float64) {
	if k.eingelocht {
		return
	}
	// Kugel vorher
	vx, vy := k.v.X(), k.v.Y()
	xK, yK := k.pos.X(), k.pos.Y()

	var willHit bool
	// reflektiere die Kugel
	var berührt bool = !((xK >= k.r) && (xK <= länge-k.r) && (yK >= k.r) && (yK <= breite-k.r))
	if !berührt && xK+vx < k.r {
		vx *= -1
		willHit = true
	}
	if !berührt && xK+vx > länge-k.r {
		vx *= -1
		willHit = true
	}
	if !berührt && yK+vy < k.r {
		vy *= -1
		willHit = true
	}
	if !berührt && yK+vy > breite-k.r {
		vy *= -1
		willHit = true
	}

	if willHit {
		klaenge.BallHitsRailSound().Play()
		k.v = hilf.V2(vx, vy)
	}
}

func (k1 *kugel) prüfeKugelKollision(k2 Kugel) {
	if k1.istKollMit == k2 {
		//if (k1.wert == 7 && k2.GibWert() == 2) || (k1.wert == 2 && k2.GibWert() == 7) {
		//	fmt.Printf("schon geprüft: %d->%d\n", k1.wert, k2.GibWert())
		//}
		return
	}

	v1 := k1.GibV()
	v2 := k2.GibV()
	distPre := k2.GibPos().Plus(v2).Minus(k1.pos.Plus(v1))
	distAkt := k2.GibPos().Minus(k1.pos)

	// Kugeln werden sich gar nicht berühren.
	if distPre.Betrag() > (k1.r + k2.GibRadius()) {
		return
	}

	// Kugeln überlappen!
	//	TODO: darf nicht sein - darf zumindest nicht so bleiben
	überlappen := distAkt.Betrag() < k1.r+k2.GibRadius()
	//if überlappen {
	//if (k1.wert == 7 && k2.GibWert() == 2) || (k1.wert == 2 && k2.GibWert() == 7) {
	//	fmt.Printf("   --> überlappen um %04.1f\n", k1.r+k2.GibRadius()-distAkt.Betrag())
	//}

	// die Stoßnormale geht durch die Mittelpunkte der Kugeln
	n12 := distPre.Normiert()
	// Zerlege Geschwindigkeiten in eine parallele und eine orthogonale Komponente
	v1p := v1.ProjiziertAuf(n12)
	v1o := v1.Minus(v1p)
	v2p := v2.ProjiziertAuf(n12)
	v2o := v2.Minus(v2p)
	// Tausche Geschwindigkeiten parallel zur Normalen aus
	var u1, u2 hilf.Vec2
	if überlappen {
		// Überlappung lösen, sonst rattern die Kugeln zusammen
		if distPre.Betrag() < distAkt.Betrag() {
			u1 = v2p.Plus(v1o)
			u2 = v1p.Plus(v2o)
		} else {
			u1 = v1
			u2 = v2
		}
	} else {
		klaenge.BallHitsBallSound().Play()
		u1 = v2p.Plus(v1o)
		u2 = v1p.Plus(v2o)
	}
	k1.SetzeV(u1)
	k2.SetzeV(u2)
	k1.istKollMit = k2
	k2.SetzeKollidiertMit(k1)
}

func (k1 *kugel) SetzeKollidiertMit(k2 Kugel) {
	k1.istKollMit = k2
}

func (k1 *kugel) SetzeKollidiertZurück() {
	k1.istKollMit = nil
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

func (k *kugel) GibWert() uint8 {
	return k.wert
}
