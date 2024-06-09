package modelle

import (
	"../hilf"
)

type mbkugel struct {
	pos, v     hilf.Vec2
	r          float64
	wert       uint8
	istKollMit MBKugel
	eingelocht bool
}

func NewKugel(pos hilf.Vec2, r float64, wert uint8) *mbkugel {
	return &mbkugel{
		pos: pos, v: hilf.V2null(),
		r:    r,
		wert: wert}
}

func (k *mbkugel) GibKopie() MBKugel {
	return &mbkugel{
		pos: k.pos, v: hilf.V2(0, 0),
		r:          k.r,
		wert:       k.wert,
		eingelocht: k.eingelocht}
}

// Bewege Kugel einen Tick weiter.
func (k *mbkugel) Bewegen() {
	k.pos = k.pos.Plus(k.v)
	vabs := k.v.Betrag()
	// Bremse die Kugel etwas ab.
	if vabs > 0.15 {
		k.v = k.v.Mal(1 - 0.02/vabs)
	} else {
		k.v = hilf.V2(0, 0)
	}

}

func (k *mbkugel) PruefeBandenKollision(breite, hoehe float64, notiereKollision func(MBKugel)) {
	if k.eingelocht {
		return
	}
	// Kugel vorher
	vx, vy := k.v.X(), k.v.Y()
	xK, yK := k.pos.X(), k.pos.Y()

	var willHit bool
	// reflektiere die Kugel
	var berührt bool = !((xK >= k.r) && (xK <= breite-k.r) && (yK >= k.r) && (yK <= hoehe-k.r))
	if !berührt && xK+vx < k.r {
		vx *= -1
		willHit = true
	}
	if !berührt && xK+vx > breite-k.r {
		vx *= -1
		willHit = true
	}
	if !berührt && yK+vy < k.r {
		vy *= -1
		willHit = true
	}
	if !berührt && yK+vy > hoehe-k.r {
		vy *= -1
		willHit = true
	}

	if willHit {
		notiereKollision(k)
		k.v = hilf.V2(vx, vy)
	}
}

// Vor.: Kugeln dürfen *nicht* still stehen
// und gleichzeitig übereinander liegen (überlappen)
func (k1 *mbkugel) PruefeKollisionMit(k2 MBKugel, notiereKollision func(MBKugel)) {
	if k1.istKollMit == k2 {
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
	überlappen := distAkt.Betrag() < k1.r+k2.GibRadius()

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
		// Überlappung lösen
		if distPre.Betrag() < distAkt.Betrag() {
			u1 = v2p.Plus(v1o)
			u2 = v1p.Plus(v2o)
		} else {
			u1 = v1
			u2 = v2
		}
	} else {
		u1 = v2p.Plus(v1o)
		u2 = v1p.Plus(v2o)
		notiereKollision(k2)
	}
	k1.SetzeV(u1)
	k2.SetzeV(u2)
	k1.istKollMit = k2
	k2.SetzeKollidiertMit(k1)
}

func (k1 *mbkugel) SetzeKollidiertMit(k2 MBKugel) {
	k1.istKollMit = k2
}

func (k1 *mbkugel) SetzeKollidiertZurueck() {
	k1.istKollMit = nil
}

func (k *mbkugel) IstEingelocht() bool { return k.eingelocht }

func (k *mbkugel) SetzeEingelocht() { k.eingelocht = true }

func (k *mbkugel) GibV() hilf.Vec2 {
	return k.v
}

func (k *mbkugel) SetzeV(v hilf.Vec2) {
	k.v = v
}

func (k *mbkugel) Stop() {
	k.v = hilf.V2(0, 0)
}

func (k *mbkugel) GibPos() hilf.Vec2 {
	return k.pos
}

func (k *mbkugel) SetzePos(pos hilf.Vec2) {
	k.pos = pos
}

func (k *mbkugel) GibRadius() float64 {
	return k.r
}

func (k *mbkugel) GibWert() uint8 {
	return k.wert
}
