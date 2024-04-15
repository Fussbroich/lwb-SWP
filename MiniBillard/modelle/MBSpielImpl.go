package modelle

import (
	"math/rand"
	"time"

	"../hilf"
	"../klaenge"
)

type mbspiel struct {
	breite       float64
	hoehe        float64
	rk           float64
	kugeln       []MBKugel
	origKugeln   []MBKugel
	vorigeKugeln []MBKugel
	spielkugel   MBKugel
	stossricht   hilf.Vec2
	stosskraft   float64
	taschen      []MBTasche
	eingelochte  []MBKugel
	strafPunkte  uint8
	stillstand   bool
	updater      hilf.Routine
	startzeit    time.Time
	spielzeit    time.Duration
	countdown    Countdown
	zeitlupe     uint64
	regelPruefer func()
}

func NewMini9BallSpiel(br, hö, ra uint16) *mbspiel {
	// Pool-Tisch:  2540 mm × 1270 mm (2:1)
	// Pool-Kugeln: 57,2 mm
	var breite, höhe, rK float64 = float64(br), float64(hö), float64(ra)
	sp := &mbspiel{breite: breite, hoehe: höhe, rk: rK, stossricht: hilf.V2null()}

	rt, rtm := 1.9*sp.rk, 1.5*sp.rk // Radien der Taschen
	sp.setzeTaschen(
		NewTasche(hilf.V2(0, 0), rt),
		NewTasche(hilf.V2(0, höhe), rt),
		NewTasche(hilf.V2(breite/2, höhe), rtm),
		NewTasche(hilf.V2(breite, höhe), rt),
		NewTasche(hilf.V2(breite, 0), rt),
		NewTasche(hilf.V2(breite/2, 0), rtm))
	sp.SetzeKugeln(sp.KugelSatz9Ball()...)
	sp.spielzeit = 4 * time.Minute
	sp.regelPruefer = func() {}
	sp.countdown = NewCountdown(sp.spielzeit)
	return sp
}

// ######## ein paar Hilfsfunktionen #########################################
func (s *mbspiel) setzeTaschen(t ...MBTasche) {
	s.taschen = []MBTasche{}
	s.taschen = append(s.taschen, t...)
}

func (s *mbspiel) SetzeKugeln(k ...MBKugel) {
	s.kugeln = []MBKugel{}
	for _, k := range k {
		k.Stop()
		s.kugeln = append(s.kugeln, k.GibKopie())
	}
	s.spielkugel = s.kugeln[0]

	// sichere den Zustand vor dem Anstoß
	s.origKugeln = []MBKugel{}
	for _, k := range s.kugeln {
		s.origKugeln = append(s.origKugeln, k.GibKopie())
	}
}

func (s *mbspiel) SetzeKugelnTest() {
	pStoß := hilf.V2(s.breite-5*s.rk, s.hoehe-5*s.rk)
	p1 := hilf.V2(s.breite-2*s.rk, s.hoehe-2*s.rk)
	s.SetzeKugeln(
		NewKugel(pStoß, s.rk, 0),
		NewKugel(p1, s.rk, 1))
}

func (sp *mbspiel) SetzeKugeln3er() { sp.SetzeKugeln(sp.kugelSatz3er()...) }

func (s *mbspiel) kugelSatz3er() []MBKugel {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pStoß := hilf.V2(s.breite-r.Float64()*(s.breite/4), r.Float64()*s.hoehe)
	dx, dy := 0.866*(2*s.rk+1), 0.5*(2*s.rk+1)
	p1 := hilf.V2(s.breite/4, s.hoehe/2)
	p2 := p1.Plus(hilf.V2(-dx, -dy))
	p3 := p1.Plus(hilf.V2(-dx, dy))
	return []MBKugel{NewKugel(pStoß, s.rk, 0),
		NewKugel(p1, s.rk, 1),
		NewKugel(p2, s.rk, 2),
		NewKugel(p3, s.rk, 3)}
}

func (sp *mbspiel) SetzeKugeln9Ball() { sp.SetzeKugeln(sp.KugelSatz9Ball()...) }

func (s *mbspiel) KugelSatz9Ball() []MBKugel {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pStoß := hilf.V2(s.breite-s.rk-r.Float64()*(s.breite/4-2*s.rk),
		r.Float64()*(s.hoehe-2*s.rk)+s.rk)
	dx, dy := 0.866*(2*s.rk+1), 0.5*(2*s.rk+1)
	p1 := hilf.V2(s.breite/4, s.hoehe/2)
	//
	p2 := p1.Plus(hilf.V2(-dx, -dy))
	p3 := p1.Plus(hilf.V2(-dx, dy))
	//
	p4 := p1.Plus(hilf.V2(-2*dx, -2*dy))
	p9 := p1.Plus(hilf.V2(-2*dx, 0))
	p5 := p1.Plus(hilf.V2(-2*dx, 2*dy))
	//
	p6 := p1.Plus(hilf.V2(-3*dx, -dy))
	p7 := p1.Plus(hilf.V2(-3*dx, dy))
	//
	p8 := p1.Plus(hilf.V2(-4*dx, 0))
	return []MBKugel{NewKugel(pStoß, s.rk, 0),
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
func (s *mbspiel) SetzeRegeln(f func()) {
	s.regelPruefer = f
}

func (s *mbspiel) Starte() {
	// Kann zwischendrin gestoppt (Pause) und wieder gestartet werden ...
	s.startzeit = time.Now()
	s.stosskraft = 5
	if s.updater == nil {
		s.updater = hilf.NewRoutine("Spiel-Logik",
			func() {
				// zähle Zeit herunter
				vergangen := time.Since(s.startzeit)
				if s.zeitlupe > 0 {
					vergangen = vergangen / time.Duration(s.zeitlupe)
				}
				s.countdown.ZieheAb(vergangen)
				// bewege jede Kugel
				s.startzeit = time.Now()
				for _, k := range s.kugeln {
					k.BewegenIn(s)
				}
				// prüfe Stillstand
				still := true
				for _, k := range s.kugeln {
					k.SetzeKollidiertZurueck()
					if !k.GibV().IstNull() {
						still = false
					}
				}
				s.stillstand = still
				// prüfe Regeln
				if s.regelPruefer != nil {
					s.regelPruefer()
				}
			})
	}
	// Starte den updater.
	s.countdown.Weiter()
	// ein konstanter Takt regelt die "Geschwindigkeit"
	if s.zeitlupe > 1 {
		s.updater.StarteRate(83 / uint64(s.zeitlupe))
		//s.updater.StarteLoop(time.Duration(12*s.zeitlupe) * time.Millisecond)
	} else {
		s.updater.StarteRate(83)
		//s.updater.StarteLoop(12 * time.Millisecond)
	}
}

func (s *mbspiel) Laeuft() bool {
	return s.updater != nil && s.updater.Laeuft()
}

func (s *mbspiel) Stoppe() {
	if s.updater != nil {
		s.updater.Stoppe()
	}
}

func (s *mbspiel) ZeitlupeAnAus() {
	if s.zeitlupe > 1 {
		s.zeitlupe = 1 // wieder normal schnell
	} else {
		s.zeitlupe = 5 // langsamer
	}
	if s.updater != nil && s.updater.Laeuft() {
		s.Stoppe()
		s.Starte()
	}
}

func (s *mbspiel) PauseAnAus() {
	if s.Laeuft() {
		s.Stoppe()
	} else {
		s.zeitlupe = 1
		s.Starte()
	}
}

func (s *mbspiel) IstZeitlupe() bool { return s.zeitlupe > 1 }

// ######## die Methoden zum Stoßen #################################################

func (s *mbspiel) GibVStoss() hilf.Vec2 {
	if s.stossricht == nil {
		s.stossricht = hilf.V2null()
	}
	return s.stossricht.Mal(s.stosskraft)
}

func (s *mbspiel) SetzeStossRichtung(v hilf.Vec2) { s.stossricht = v.Normiert() }

func (s *mbspiel) SetzeStosskraft(v float64) {
	// Die "Geschwindigkeit/Stärke" ist auf 12 (m/s) begrenzt
	if v > 12 {
		v = 12
	} else if v < 0 {
		v = 0
	}
	s.stosskraft = v
}

func (s *mbspiel) Stosse() {
	if !s.stillstand {
		println("Fehler: Stoßen während laufender Bewegungen ist verboten!")
		return
	}
	if s.stossricht == nil {
		s.stossricht = hilf.V2null()
	}
	// sichere den Zustand vor dem Stoß
	s.vorigeKugeln = []MBKugel{}
	for _, k := range s.kugeln {
		s.vorigeKugeln = append(s.vorigeKugeln, k.GibKopie()) // Kopien stehen still
	}
	if !(s.stosskraft == 0) {
		// stoße
		s.spielkugel.SetzeV(s.stossricht.Mal(s.stosskraft))
		s.stosskraft = 5
		s.stillstand = false
		klaenge.CueHitsBallSound()
	}
}

// ######## die übrigen Methoden ####################################################

func (s *mbspiel) SetzeSpielzeit(t time.Duration) { s.spielzeit = t }

func (s *mbspiel) SetzeRestzeit(t time.Duration) { s.countdown.Setze(t) }

func (s *mbspiel) StoppeCountdown() { s.countdown.Halt() }

func (s *mbspiel) StarteCountdown() { s.countdown.Weiter() }

func (s *mbspiel) GibRestzeit() time.Duration { return s.countdown.GibRestzeit() }

func (s *mbspiel) Reset() {
	s.kugeln = s.KugelSatz9Ball() // neue Kugeln, die alle still stehen
	s.spielkugel = s.kugeln[0]
	s.eingelochte = []MBKugel{}
	s.strafPunkte = 0
	s.stillstand = true
	s.countdown.Setze(s.spielzeit)
}

func (s *mbspiel) StossWiederholen() {
	if s.vorigeKugeln == nil {
		return
	}
	// stelle den Zustand vor dem letzten Stoß wieder her
	s.kugeln = []MBKugel{}
	for _, k := range s.vorigeKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie()) // Kopien stehen still
	}
	s.spielkugel = s.kugeln[0]
	s.stillstand = true
}

func (s *mbspiel) GibKugeln() []MBKugel { return s.kugeln }

func (s *mbspiel) GibSpielkugel() MBKugel { return s.spielkugel }

func (s *mbspiel) GibAktiveKugeln() []MBKugel {
	ks := []MBKugel{}
	for _, k := range s.kugeln {
		if k.IstEingelocht() {
			continue
		}
		ks = append(ks, k)
	}
	return ks
}

func (s *mbspiel) Einlochen(k MBKugel) {
	if k == s.spielkugel {
		return
	}
	//TODO: Mengen benutzen
	for _, ke := range s.eingelochte {
		if k.GibWert() == ke.GibWert() {
			return
		}
	}
	s.eingelochte = append(s.eingelochte, k)
	if s.regelPruefer != nil {
		s.regelPruefer()
	}
}

func (s *mbspiel) GibEingelochteKugeln() []MBKugel { return s.eingelochte }

func (s *mbspiel) GibGroesse() (float64, float64) { return s.breite, s.hoehe }

func (s *mbspiel) GibTaschen() []MBTasche { return s.taschen }

func (s *mbspiel) IstStillstand() bool { return s.stillstand }

func (s *mbspiel) GibTreffer() uint8 { return uint8(len(s.eingelochte)) }

func (s *mbspiel) GibStrafpunkte() uint8 { return s.strafPunkte }

func (s *mbspiel) ReduziereStrafpunkte() {
	if s.strafPunkte > 0 {
		s.strafPunkte--
	}
}

func (s *mbspiel) ErhoeheStrafpunkte() {
	s.strafPunkte++
}
