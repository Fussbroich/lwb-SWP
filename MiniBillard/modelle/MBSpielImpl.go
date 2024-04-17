package modelle

import (
	"math/rand"
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
	foulPruefer  func() bool
	startzeit    time.Time
	spielzeit    time.Duration // zum Spiel gegen die Zeit
	countdown    Countdown
	zeitlupe     uint64 // Testzwecke
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
	sp.foulPruefer = sp.isFoul9Ball
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

/*
Vor. Ein Anstoß war erfolgt und jetzt stehen die Kugeln alle wieder still.
Erg. True bedeutet, dass es ein Foul gibt.
*/
func (sp *mbspiel) isFoul9Ball() bool {
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

// Ein Spiel mit 3 Kugeln
func NewMini3BallSpiel(br, hö uint16) *mbspiel {
	sp := newPoolSpiel(br, hö)
	sp.setzeKugeln(sp.kugelSatz3er()...)
	sp.neusetzenSpielkugel()
	sp.spielzeit = 90 * time.Second
	sp.countdown = NewCountdown(sp.spielzeit)
	sp.foulPruefer = sp.isFoul9Ball
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
	s.taschen = []MBTasche{}
	s.taschen = append(s.taschen, t...)
}

func (s *mbspiel) setzeKugeln(k ...MBKugel) {
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
	s.spielkugel.SetzePos(pStoß)
}

func (s *mbspiel) GibGroesse() (float64, float64) { return s.breite, s.hoehe }

func (s *mbspiel) GibTaschen() []MBTasche { return s.taschen }

func (s *mbspiel) GibTreffer() uint8 { return uint8(len(s.eingelochte)) }

func (s *mbspiel) GibStrafpunkte() uint8 { return s.strafPunkte }

func (s *mbspiel) ReduziereStrafpunkte() {
	if s.strafPunkte > 0 {
		s.strafPunkte--
	}
}

// ######## die Lebens- und Pause-Methode ###########################################################
func (s *mbspiel) Starte() {
	// Kann zwischendrin gestoppt (Pause) und wieder gestartet werden ...
	if s.sollRate == 0 {
		// Todo Simulation ist derzeit von der Auflösung und der Rate abhängig
		s.sollRate = 83
	}
	if s.countdown == nil {
		if s.spielzeit == 0 {
			s.spielzeit = 3 * time.Minute
		}
		s.countdown = NewCountdown(s.spielzeit)
	}
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
				// prüfe Fouls und Sieg nach Stillstand
				if s.angestossen && s.stillstand {
					s.angestossen = false
					if s.foulPruefer != nil && s.foulPruefer() {
						s.strafPunkte++
					}
					if s.spielkugel.IstEingelocht() {
						s.StossWiederholen()
					}
					// Das Spiel ist gewonnen
					if s.AlleEingelocht() {
						s.countdown.Halt()
					}
				}
			})
	}
	// Starte den updater.
	s.countdown.Weiter()
	// ein konstanter Takt regelt die "Geschwindigkeit"
	if s.zeitlupe > 1 {
		s.updater.StarteRate(s.sollRate / uint64(s.zeitlupe))
		//s.updater.StarteLoop(time.Duration(12*s.zeitlupe) * time.Millisecond)
	} else {
		s.updater.StarteRate(s.sollRate)
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

func (s *mbspiel) SetzeStossRichtung(v hilf.Vec2) {
	s.stossricht = v.Normiert()
}

func (s *mbspiel) SetzeStosskraft(v float64) {
	// Die "Geschwindigkeit/Stärke" ist auf 14 (m/s) begrenzt
	if v > 14 {
		v = 14
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
	// stoßen
	s.angestossen = true
	s.angespielte = nil
	if !(s.stosskraft == 0) {
		s.spielkugel.SetzeV(s.stossricht.Mal(s.stosskraft))
		s.stosskraft = 5
		s.stillstand = false
		klaenge.CueHitsBallSound()
	}
}

// ######## die übrigen Methoden ####################################################

func (s *mbspiel) SetzeSpielzeit(t time.Duration) { s.spielzeit = t; s.countdown.Setze(t) }

// Testzwecke
func (s *mbspiel) SetzeRestzeit(t time.Duration) { s.countdown.Setze(t) }

func (s *mbspiel) StoppeCountdown() { s.countdown.Halt() }

func (s *mbspiel) StarteCountdown() { s.countdown.Weiter() }

func (s *mbspiel) GibRestzeit() time.Duration { return s.countdown.GibRestzeit() }

func (s *mbspiel) Reset() {
	s.kugeln = []MBKugel{}
	for _, k := range s.origKugeln {
		s.kugeln = append(s.kugeln, k.GibKopie()) // Kopien stehen still
	}
	s.spielkugel = s.kugeln[0]
	s.neusetzenSpielkugel()
	s.eingelochte = []MBKugel{}
	s.angestossen = false
	s.strafPunkte = 0
	s.stillstand = true
	s.countdown.Setze(s.spielzeit)
	s.countdown.Weiter()
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
	s.angestossen = false
	s.stillstand = true
}

func (s *mbspiel) GibKugeln() []MBKugel { return s.kugeln }

func (s *mbspiel) GibSpielkugel() MBKugel { return s.spielkugel }

func (s *mbspiel) GibAktiveKugeln() []MBKugel {
	ks := []MBKugel{}
	for _, k := range s.kugeln {
		if !k.IstEingelocht() {
			ks = append(ks, k)

		}
	}
	return ks
}

func (s *mbspiel) AlleEingelocht() bool {
	aktive := []MBKugel{}
	for _, k := range s.kugeln {
		if !k.IstEingelocht() {
			aktive = append(aktive, k)
		}
	}
	return len(aktive) == 0 || (len(aktive) == 1 && aktive[0].GibWert() == 0)
}

func (s *mbspiel) IstStillstand() bool { return s.stillstand }

// eine kleine Buchhaltung für Berührungen (z.B. die angespielte Kugel)
func (s *mbspiel) NotiereBerührt(k1 MBKugel, k2 MBKugel) {
	if s.angespielte == nil {
		if k1 == s.spielkugel {
			s.angespielte = k2
		} else if k2 == s.spielkugel {
			s.angespielte = k1
		}
	}
}

// eine kleine Buchhaltung für eingelochte Kugeln
func (s *mbspiel) NotiereEingelocht(k MBKugel) {
	if k == s.spielkugel {
		return
	}
	// eingelochte ein Mal der Reihe nach speichern
	// Todo: Hier eine Menge nehmen
	for _, ke := range s.eingelochte {
		if k.GibWert() == ke.GibWert() {
			return
		}
	}
	s.eingelochte = append(s.eingelochte, k)
}

func (s *mbspiel) GibEingelochteKugeln() []MBKugel { return s.eingelochte }
