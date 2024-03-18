package welt

import (
	"time"

	"../hilf"
	"../klaenge"
)

type MiniBillardSpiel interface {
	Starte()
	Stoppe()
	Läuft() bool
	ZeitlupeAnAus()
	IstZeitlupe() bool
	Stoße()
	StoßWiederholen()
	Reset()
	IstStillstand() bool
	GibTaschen() []Tasche
	GibKugeln() []Kugel
	GibAktiveKugeln() []Kugel
	GibEingelochteKugeln() []Kugel
	GibStoßkugel() Kugel
	GibVStoß() hilf.Vec2
	SetzeVStoß(hilf.Vec2)
	SetzeRestzeit(time.Duration)
	GibRestzeit() time.Duration
	GibTreffer() uint8
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
	vStoß        hilf.Vec2
	taschen      []Tasche
	stößeBisher  uint8
	strafPunkte  uint8
	stillstand   bool
	updater      hilf.Prozess
	startzeit    time.Time
	restzeit     time.Duration
	zeitlupe     uint64
}

func NewMiniPoolSpiel(br, hö, ra uint16) *spiel {
	// Pool-Tisch:  2540 mm × 1270 mm (2:1)
	// Pool-Kugeln: 57,2 mm
	var breite, höhe, rK float64 = float64(br), float64(hö), float64(ra)
	sp := &spiel{breite: breite, höhe: höhe, rk: rK}
	rt, rtm := 1.9*sp.rk, 1.5*sp.rk
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

// ######## ein paar Hilfsfunktionen #########################################
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
	//
	p2 := p1.Plus(pos(-dx, -dy))
	p3 := p1.Plus(pos(-dx, dy))
	//
	p4 := p1.Plus(pos(-2*dx, -2*dy))
	p9 := p1.Plus(pos(-2*dx, 0))
	p5 := p1.Plus(pos(-2*dx, 2*dy))
	//
	p6 := p1.Plus(pos(-3*dx, -dy))
	p7 := p1.Plus(pos(-3*dx, dy))
	//
	p8 := p1.Plus(pos(-4*dx, 0))
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

// ######## die Lebens- und Pause-Methode ###########################################################
func (s *spiel) Starte() {
	s.startzeit = time.Now()
	if s.updater == nil {
		s.updater = hilf.NewProzess("Spiel-Logik",
			func() {
				vergangen := time.Since(s.startzeit)
				if s.zeitlupe > 0 {
					vergangen = vergangen / time.Duration(s.zeitlupe)
				}
				if s.restzeit <= vergangen {
					s.restzeit = 0
				} else {
					s.restzeit -= vergangen
				}
				s.startzeit = time.Now()
				for _, k := range s.kugeln {
					k.BewegenIn(s)
				}
				still := true
				for _, k := range s.kugeln {
					k.SetzeKollidiertZurück()
					//prüfe Stillstand
					if !k.GibV().IstNull() {
						still = false
					}
				}
				s.stillstand = still
			})
	}
	// ein konstanter Takt regelt die "Geschwindigkeit"
	if s.zeitlupe > 1 {
		s.updater.StarteRate(83 / uint64(s.zeitlupe))
		//s.updater.StarteLoop(time.Duration(12*s.zeitlupe) * time.Millisecond)
	} else {
		s.updater.StarteRate(83)
		//s.updater.StarteLoop(12 * time.Millisecond)
	}
}

func (s *spiel) Läuft() bool {
	return s.updater != nil && s.updater.Läuft()
}

func (s *spiel) Stoppe() {
	if s.updater != nil {
		s.updater.Stoppe()
	}
}

func (s *spiel) ZeitlupeAnAus() {
	if s.zeitlupe > 1 {
		s.zeitlupe = 1 // wieder normal schnell
	} else {
		s.zeitlupe = 10 // langsamer
	}
	if s.updater != nil && s.updater.Läuft() {
		s.Stoppe()
		s.Starte()
	}
}

func (s *spiel) IstZeitlupe() bool { return s.zeitlupe > 1 }

// ######## die Methoden zum Stoßen #################################################

func (s *spiel) GibVStoß() hilf.Vec2 { return s.vStoß }

func (s *spiel) SetzeVStoß(v hilf.Vec2) {
	vabs := v.Betrag()
	// Die "Geschwindigkeit/Stärke" ist auf 17 (m/s) begrenzt
	if vabs > 12 {
		s.vStoß = v.Mal(12 / vabs)
	} else {
		s.vStoß = v
	}
}

func (s *spiel) Stoße() {
	if !s.stillstand {
		println("Fehler: Stoßen während laufender Bewegungen ist verboten!")
		return
	}
	// sichere den Zustand vor dem Stoß
	s.vorigeKugeln = []Kugel{}
	for _, k := range s.kugeln {
		s.vorigeKugeln = append(s.vorigeKugeln, k.GibKopie()) // Kopien stehen still
	}
	if !s.vStoß.IstNull() {
		// stoße
		s.stoßkugel.SetzeV(s.vStoß)
		s.vStoß = hilf.V2(0, 0)
		klaenge.CueHitsBallSound()
		s.stößeBisher++
		s.stillstand = false
	}
}

// ######## die übrigen Methoden ####################################################

func (s *spiel) SetzeRestzeit(t time.Duration) { s.restzeit = t }

func (s *spiel) GibRestzeit() time.Duration { return s.restzeit }

func (s *spiel) Reset() {
	s.kugeln = []Kugel{}
	for _, k := range s.origKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie()) // Kopien stehen still
	}
	s.stoßkugel = s.kugeln[0]
	s.stößeBisher = 0
	s.strafPunkte = 0
	s.stillstand = true
}

func (s *spiel) StoßWiederholen() {
	if s.vorigeKugeln == nil {
		return
	}
	// stelle den Zustand vor dem letzten Stoß wieder her
	s.kugeln = []Kugel{}
	for _, k := range s.vorigeKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie()) // Kopien stehen still
	}
	s.stoßkugel = s.kugeln[0]
	s.strafPunkte++
	s.stillstand = true
}

func (s *spiel) GibKugeln() []Kugel { return s.kugeln }

func (s *spiel) GibStoßkugel() Kugel { return s.stoßkugel }

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

func (s *spiel) GibGröße() (float64, float64) { return s.breite, s.höhe }

func (s *spiel) GibTaschen() []Tasche { return s.taschen }

func (s *spiel) IstStillstand() bool { return s.stillstand }

func (s *spiel) GibTreffer() (treffer uint8) {
	for _, k := range s.kugeln {
		if !(k == s.stoßkugel) && k.IstEingelocht() {
			treffer++
		}
	}
	return
}

func (s *spiel) GibStrafpunkte() uint8 { return s.strafPunkte }
