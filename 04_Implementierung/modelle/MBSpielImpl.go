package modelle

import (
	"math/rand"
	"sync"
	"time"

	"../hilf"
	"../klaenge"
)

type mbspiel struct {
	breite       float64 // Länge des Tuchs in der Simulation
	hoehe        float64 // Breite des Tuchs in der Simulation
	rk           float64 // der Radius aller Kugeln in der Simulation
	kugeln       []MBKugel
	origKugeln   []MBKugel // der ganze Kugelsatz vor dem Spielbeginn
	vorigeKugeln []MBKugel // der Kugelsatz vor dem letzten Stoß
	spielkugel   MBKugel   // die weiße Kugel
	angespielte  MBKugel   // die zuerst berührte Kugel, falls überhaupt
	angestossen  bool      // ein Stoß hat stattgefunden
	stossricht   hilf.Vec2
	stosskraft   float64
	taschen      []MBTasche
	eingelochte  []MBKugel // eine simple Buchhaltung der eingelochten
	strafPunkte  uint8
	stillstand   bool // alle Kugeln stehen still
	updater      hilf.Routine
	sollRate     uint64 // die Wunschgeschwindigkeit der Simulation
	vorigeZeit   time.Time
	spielzeit    time.Duration // zum Spiel gegen die Zeit
	countdown    Countdown
	rwZustand    sync.RWMutex // insbesondere Fouls *müssen* korrekt gezählt sein
}

// Ein Pool-Spiel ohne Kugeln
func newPoolSpiel(br, hö uint16) *mbspiel {
	// echter Pool-Tisch:  2540 mm × 1270 mm (2:1)
	// Der Aufrufer darf sich hier auch ein anderes Seitenverhältnis wünschen.
	// Pool-Kugeln: 57,2 mm
	// Radius der Kugeln
	var breite, höhe float64 = float64(br), float64(hö)
	var rK float64 = breite * 57.2 / 2540
	sp := &mbspiel{breite: breite, hoehe: höhe, rk: rK, stossricht: hilf.V2null()}

	// Radien der Taschen sind groß, damit die Kugeln auch reingehen
	rt, rtm := 1.9*sp.rk, 1.5*sp.rk
	sp.setzeTaschen(
		NewTasche(hilf.V2(0, 0), rt),
		NewTasche(hilf.V2(0, höhe), rt),
		NewTasche(hilf.V2(breite/2, höhe), rtm),
		NewTasche(hilf.V2(breite, höhe), rt),
		NewTasche(hilf.V2(breite, 0), rt),
		NewTasche(hilf.V2(breite/2, 0), rtm))
	return sp
}

// Eine schwarze Kugel und eine weiße zu Testzwecken
func (s *mbspiel) SetzeKugeln1BallTest() {
	pStoß := hilf.V2(s.breite-5*s.rk, s.hoehe-5*s.rk)
	p8 := hilf.V2(s.breite-2*s.rk, s.hoehe-2*s.rk)
	s.setzeKugeln(
		NewKugel(pStoß, s.rk, 0),
		NewKugel(p8, s.rk, 8))
}

// Ein Spiel mit 9 Kugeln
func (s *mbspiel) SetzeKugeln9Ball() {
	s.setzeKugeln(s.kugelSatz9Ball()...)
}

func NewMini9BallSpiel(br, hö uint16) *mbspiel {
	sp := newPoolSpiel(br, hö)
	sp.setzeKugeln(sp.kugelSatz9Ball()...)
	sp.neusetzenSpielkugel()
	sp.spielzeit = 4 * time.Minute
	sp.countdown = NewCountdown(sp.spielzeit)
	return sp
}

func (s *mbspiel) kugelSatz9Ball() []MBKugel {
	// Lege eine Raute mit der 9 in der Mitte
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
	return []MBKugel{NewKugel(hilf.V2(3*s.breite/4, s.hoehe/2), s.rk, 0),
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

// Ein Spiel mit 3 Kugeln
func NewMini3BallSpiel(br, hö uint16) *mbspiel {
	sp := newPoolSpiel(br, hö)
	sp.setzeKugeln(sp.kugelSatz3er()...)
	sp.neusetzenSpielkugel()
	sp.spielzeit = 90 * time.Second
	sp.countdown = NewCountdown(sp.spielzeit)
	return sp
}

func (s *mbspiel) SetzeKugeln3Ball() {
	s.setzeKugeln(s.kugelSatz3er()...)
}

func (s *mbspiel) kugelSatz3er() []MBKugel {
	dx, dy := 0.866*(2*s.rk+1), 0.5*(2*s.rk+1)
	p1 := hilf.V2(s.breite/4, s.hoehe/2)
	p2 := p1.Plus(hilf.V2(-dx, -dy))
	p3 := p1.Plus(hilf.V2(-dx, dy))
	return []MBKugel{NewKugel(hilf.V2(3*s.breite/4, s.hoehe/2), s.rk, 0),
		NewKugel(p1, s.rk, 1),
		NewKugel(p2, s.rk, 2),
		NewKugel(p3, s.rk, 3)}
}

// ######## ein paar Hilfsfunktionen #########################################
func (s *mbspiel) setzeTaschen(t ...MBTasche) {
	//s.rwZustand.Lock()
	//defer s.rwZustand.Unlock()

	s.taschen = []MBTasche{}
	s.taschen = append(s.taschen, t...)
}

func (s *mbspiel) setzeKugeln(k ...MBKugel) {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.kugeln = []MBKugel{}
	for _, k := range k {
		s.kugeln = append(s.kugeln, k.GibKopie())
	}
	s.spielkugel = s.kugeln[0]

	// sichere den Zustand vor dem Break
	s.origKugeln = []MBKugel{}
	for _, k := range s.kugeln {
		s.origKugeln = append(s.origKugeln, k.GibKopie())
	}
}

func (s *mbspiel) neusetzenSpielkugel() {
	if s.spielkugel == nil {
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pStoß := hilf.V2(s.breite-s.rk-r.Float64()*(s.breite/4-2*s.rk),
		r.Float64()*(s.hoehe-2*s.rk)+s.rk)

	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.spielkugel.SetzePos(pStoß)
}

func (s *mbspiel) GibGroesse() (float64, float64) {
	return s.breite, s.hoehe
}

func (s *mbspiel) GibTaschen() []MBTasche {
	return s.taschen
}

func (s *mbspiel) GibTreffer() uint8 {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return uint8(len(s.eingelochte))

}

func (s *mbspiel) GibStrafpunkte() uint8 {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return s.strafPunkte
}

func (s *mbspiel) ReduziereStrafpunkte() {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	if s.strafPunkte > 0 {
		s.strafPunkte--
	}
}

func (s *mbspiel) ErhoeheStrafpunkte() {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.strafPunkte++
}

// ######## die Lebens-Methode ###########################################################
/*
Vor. Ein Anstoß war erfolgt und jetzt stehen die Kugeln alle wieder still.
Erg. True bedeutet, dass es ein Foul gibt.
*/
func (sp *mbspiel) isFoul() bool {
	sp.rwZustand.RLock()
	defer sp.rwZustand.RUnlock()

	// Spielkugel einlochen ist ein Foul
	if sp.spielkugel.IstEingelocht() {
		return true
	}
	// gar keine Kugel anspielen ist ein Foul
	if sp.angespielte == nil {
		return true
	}
	// gar keine versenken ist ein Foul
	var anzVor, anzNeu int
	for _, k := range sp.vorigeKugeln {
		if k.IstEingelocht() {
			anzVor++
		}
	}
	anzNeu = len(sp.eingelochte)
	return !(anzVor < anzNeu)
}

// eine kleine Buchhaltung für (Erst-)Berührungen
func (s *mbspiel) notiereBerührt(k1 MBKugel, k2 MBKugel) {
	// nur das Anspielen einer Kugel durch die weiße
	// merken und alle weiteren Berührungen ignorieren
	if s.angespielte == nil {
		if k1 == s.spielkugel {
			s.angespielte = k2
		} else if k2 == s.spielkugel {
			s.angespielte = k1
		}
	}
}

// eine kleine Buchhaltung für eingelochte Kugeln
func (s *mbspiel) notiereEingelocht(k MBKugel) {
	if k == s.spielkugel {
		return
	}
	// eingelochte der Reihe nach merken
	for _, ke := range s.eingelochte {
		if k.GibWert() == ke.GibWert() { // TODO Hier eine Menge nehmen
			return
		}
	}
	s.eingelochte = append(s.eingelochte, k)
}

func (s *mbspiel) alleEingelocht() bool {
	var aktive uint8
	for _, k := range s.kugeln {
		if !k.IstEingelocht() && !(k == s.spielkugel) {
			aktive++
		}
	}
	return aktive == 0
}

// Die Lebensmethode - wird bei jedem Tick einmal aufgerufen
func (s *mbspiel) Update() {
	s.rwZustand.Lock()
	if s.countdown == nil {
		if s.spielzeit == 0 {
			s.spielzeit = 4 * time.Minute
		}
		s.countdown = NewCountdown(s.spielzeit)
	}
	// zähle Zeit herunter
	s.countdown.ZieheAb(time.Since(s.vorigeZeit))
	s.vorigeZeit = time.Now()

	// update jede Kugel
	for _, k := range s.kugeln {
		if !k.IstEingelocht() {
			for _, k2 := range s.kugeln {
				if (k == k2) || k2.IstEingelocht() {
					continue
				}
				k.PruefeKollisionMit(k2,
					func(n MBKugel) {
						klaenge.BallHitsBallSound().Play()
						s.notiereBerührt(k, n)
					})
			}
			k.SetzeKollidiertZurueck()
			l, b := s.GibGroesse()
			k.PruefeBandenKollision(l, b,
				func(MBKugel) {
					klaenge.BallHitsRailSound().Play()
				})
			// Kugel einen Tick weiter bewegen
			k.Bewegen()
			// Prüfe, ob Kugel eingelocht wurde.
			for _, t := range s.GibTaschen() {
				if t.GibPos().Minus(k.GibPos()).Betrag() < t.GibRadius() {
					klaenge.BallInPocketSound().Play()
					k.SetzeEingelocht()
					k.Stop()
					s.notiereEingelocht(k)
					break
				}
			}
		}
	}

	// prüfe Stillstand
	still := true
	for _, k := range s.kugeln {
		k.SetzeKollidiertZurueck()
		if !k.GibV().IstNull() {
			still = false
		}
	}

	if still {
		s.stillstand = true
	}
	s.rwZustand.Unlock()

	// prüfe Fouls und Sieg nach Stillstand
	if s.angestossen && s.stillstand {
		s.angestossen = false

		if s.isFoul() {
			s.ErhoeheStrafpunkte()
		}
		if s.spielkugel.IstEingelocht() {
			s.StossWiederholen()
		}
		// Das Spiel ist gewonnen
		if s.alleEingelocht() {
			s.countdown.Halt()
		}
	}
}

func (s *mbspiel) Starte() {
	// Kann zwischendrin gestoppt (Pause) und wieder gestartet werden ...
	s.vorigeZeit = time.Now()
	s.stosskraft = 5
	if s.updater == nil {
		s.updater = hilf.NewRoutine("Spiel-Logik", s.Update)
	}
	// Starte den updater.
	if !s.alleEingelocht() {
		s.countdown.Weiter()
	}
	// ein konstanter Takt regelt die "Geschwindigkeit"
	if s.sollRate == 0 {
		s.sollRate = 83
	}
	s.updater.StarteMitRate(s.sollRate)
}

func (s *mbspiel) Laeuft() bool {
	return s.updater != nil && s.updater.Laeuft()
}

func (s *mbspiel) Stoppe() {
	if s.updater != nil {
		s.updater.Stoppe()
	}
}

func (s *mbspiel) GetTicksPS() uint64 {
	if s.updater != nil {
		return s.updater.GibRate()
	} else {
		return 0
	}
}

// ######## die Methoden zum Stoßen #################################################

func (s *mbspiel) GibVStoss() hilf.Vec2 {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	if s.stossricht == nil {
		s.stossricht = hilf.V2null()
	}
	return s.stossricht.Mal(s.stosskraft)
}

func (s *mbspiel) SetzeStossRichtung(v hilf.Vec2) {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.stossricht = v.Normiert()
}

func (s *mbspiel) SetzeStosskraft(v float64) {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	// Die "Geschwindigkeit/Stärke" ist auf 14 (m/s) begrenzt
	if v > 14 {
		v = 14
	} else if v < 0 {
		v = 0
	}
	s.stosskraft = v
}

func (s *mbspiel) Stosse() {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	if !s.stillstand {
		println("Fehler: Stoßen während laufender Bewegungen ist verboten!")
		return
	}
	if s.stosskraft <= 0.15 {
		return
	}

	if s.stossricht == nil || s.stossricht.IstNull() {
		return
	}

	klaenge.CueHitsBallSound().Play()

	// sichere den Zustand vor dem Stoß
	s.vorigeKugeln = []MBKugel{}
	for _, k := range s.kugeln {
		s.vorigeKugeln = append(s.vorigeKugeln, k.GibKopie()) // Kopien stehen still
	}
	// stoßen
	s.angestossen = true
	s.angespielte = nil
	s.spielkugel.SetzeV(s.stossricht.Mal(s.stosskraft))
	s.stosskraft = 5
	s.stillstand = false
}

// ######## die übrigen Methoden ####################################################

func (s *mbspiel) SetzeSpielzeit(t time.Duration) {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.spielzeit = t
	s.countdown.Setze(t)
}

// Testzwecke
func (s *mbspiel) SetzeRestzeit(t time.Duration) {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.countdown.Setze(t)
}

func (s *mbspiel) StoppeCountdown() {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.countdown.Halt()
}

func (s *mbspiel) StarteCountdown() {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	s.countdown.Weiter()
}

func (s *mbspiel) GibRestzeit() time.Duration {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return s.countdown.GibRestzeit()
}

func (s *mbspiel) Reset() {
	s.rwZustand.Lock()
	s.kugeln = []MBKugel{}
	for _, k := range s.origKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie()) // Kopien stehen still
	}
	s.spielkugel = s.kugeln[0]
	s.eingelochte = []MBKugel{}
	s.angestossen = false
	s.strafPunkte = 0
	s.stillstand = true
	s.countdown.Setze(s.spielzeit)
	s.countdown.Weiter()
	s.rwZustand.Unlock()
	s.neusetzenSpielkugel()
}

func (s *mbspiel) StossWiederholen() {
	s.rwZustand.Lock()
	defer s.rwZustand.Unlock()

	if s.vorigeKugeln == nil {
		return
	}
	// stelle den Zustand vor dem letzten Stoß wieder her
	s.kugeln = []MBKugel{}
	for _, k := range s.vorigeKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie()) // Kopien stehen still
	}
	s.spielkugel = s.kugeln[0]
	s.angestossen = false
	s.stillstand = true
}

func (s *mbspiel) GibKugeln() []MBKugel {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return s.kugeln
}

func (s *mbspiel) GibSpielkugel() MBKugel {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return s.spielkugel
}

func (s *mbspiel) GibAktiveKugeln() []MBKugel {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	ks := []MBKugel{}
	for _, k := range s.kugeln {
		if !k.IstEingelocht() {
			ks = append(ks, k)

		}
	}
	return ks
}

func (s *mbspiel) IstStillstand() bool {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return s.stillstand
}

func (s *mbspiel) GibEingelochteKugeln() []MBKugel {
	s.rwZustand.RLock()
	defer s.rwZustand.RUnlock()

	return s.eingelochte
}
