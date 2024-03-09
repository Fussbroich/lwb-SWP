package welt

import (
	"../hilf"
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
	hatBerührt Bande
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
	for _, b := range s.GibBanden() {
		k.prüfeBandenKollision(b)
	}
	// Bewege Kugel einen Tick weiter.
	k.pos = k.pos.Plus(k.v)
	vabs := k.v.Betrag()
	// Bremse die Kugel etwas ab.
	if vabs > 0.15 {
		k.v = k.v.Mal(1 - 0.02/vabs)
	} else {
		k.v = hilf.V2(0, 0)
	}
	k.hatBerührt = nil
	k.istKollMit = nil
	// Prüfe, ob Kugel eingelocht wurde.
	for _, t := range s.GibTaschen() {
		if t.GibPos().Minus(k.GibPos()).Betrag() < t.GibRadius() {
			hilf.BallInPocketSound()
			k.eingelocht = true
			k.SetzeV(hilf.V2(0, 0))
			break
		}
	}
}

func (k *kugel) IstEingelocht() bool {
	return k.eingelocht
}

func (k *kugel) prüfeBandenKollision(b Bande) {
	if k.hatBerührt == b {
		return
	}
	if k.IstEingelocht() {
		return
	}
	// Kugel vorher
	pK := k.GibPos()
	vK := k.GibV()
	//Bande
	pBVon := b.GibVon()
	pBNach := b.GibNach()
	vBRicht := pBNach.Minus(pBVon).Normiert()
	//Stoßnormale (Lot von Kugel K zur Bande b)
	t := -((pBVon.Minus(pK)).Punkt(vBRicht) / vBRicht.Punkt(vBRicht))
	lotFußpunkt := pBVon.Plus(vBRicht.Mal(t))
	lot := lotFußpunkt.Minus(pK)

	// Kugel berührt gar nicht
	if lot.Betrag() > (k.GibRadius()) {
		return //Kugel zu weit weg
	}
	// TODO: hier passiert Mist
	//	bLänge := pBNach.Minus(pBVon).Betrag()
	//	if t < 0 || lotFußpunkt.Minus(pBVon).Betrag() > bLänge {
	//		println("Kugel nicht zwischen den Endpunkten der Bande")
	//		return
	//	}
	// reflektiere die Kugel 1mal
	hilf.BallHitsRailSound()
	norm := lot.Normiert()
	vp := vK.ProjiziertAuf(norm)
	vo := vK.Minus(vp)
	// Kugel nachher
	u := vo.Plus(vp.Mal(-1))
	k.SetzeV(u)
	k.hatBerührt = b
}

func (k1 *kugel) prüfeKugelKollision(k2 Kugel) {
	if k1.istKollMit == k2 {
		return
	}
	if k1.eingelocht || k2.IstEingelocht() {
		return
	}
	pos1 := k1.pos
	pos2 := k2.GibPos()
	v1 := k1.v
	v2 := k2.GibV()
	norm12 := pos2.Plus(v2).Minus(pos1.Plus(v1))
	// Kugeln berühren sich gar nicht.
	if norm12.Betrag() >= (k1.r + k2.GibRadius()) {
		return
	}
	hilf.BallHitsBallSound()
	n12 := norm12.Normiert()
	v1p := v1.ProjiziertAuf(n12)
	v1o := v1.Minus(v1p)
	v2p := v2.ProjiziertAuf(n12)
	v2o := v2.Minus(v2p)
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
