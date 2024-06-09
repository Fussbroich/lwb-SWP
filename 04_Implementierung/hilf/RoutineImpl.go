package hilf

import "time"

type routine struct {
	name         string
	frun         func()
	fausnahme    func()
	rate         uint64 // Hertz
	verzoegerung time.Duration
	stop         chan bool
	laeuft       bool
}

func NewRoutine(name string, f_run func()) *routine {
	return &routine{
		name:      name,
		frun:      f_run,
		fausnahme: func() {}}
}

func (r *routine) SetzeAusnahmeHandler(f func()) { r.fausnahme = f }

// Prüfe, ob die Routine noch läuft
func (r *routine) Laeuft() bool { return r.laeuft }

// Starte eine Funktion als goroutine.
// Lasse sie in einem fest vorgegebenen Takt loopen (feste Rate je Sekunde).
func (r *routine) StarteMitTakt(tick time.Duration) {
	if r.laeuft {
		println("Fehler:", r.name, "läuft bereits.")
		return
	}
	r.laeuft = true
	println("Starte Takt für", r.name)
	r.rate = 1e9 / uint64(tick.Nanoseconds())
	takt := time.NewTicker(tick)
	r.stop = make(chan bool)
	runner := func() {
		defer func() { takt.Stop() }()
		defer r.fausnahme()
		for {
			r.frun() // Starte die Funktion sofort, ohne auf den ersten Tick zu warten.
			select {
			case <-r.stop:
				println("Stoppe", r.name)
				r.laeuft = false
				return
			case <-takt.C:
				continue // Warte auf den nächsten Tick.
			}
		}
	}
	// starte Loop
	go runner()
}

// Starte eine Funktion als goroutine.
// Lasse sie - falls möglich - mit einer bestimmten Rate je Sekunde loopen.
// Die Rate wird laufend durch eine veränderliche Verzögerung nachgeführt.
func (r *routine) StarteMitRate(sollRate uint64) {
	if r.laeuft {
		println("Fehler:", r.name, "läuft bereits.")
		return
	}
	r.laeuft = true
	println("Starte", r.name, "(soll:", sollRate, "Hz)")
	r.rate = sollRate
	r.verzoegerung = 0
	maxVerzögerung := time.Second / 5
	var minRate, maxRate uint64 = sollRate * 4 / 5, sollRate * 6 / 5
	r.stop = make(chan bool)
	runner := func() {
		defer r.fausnahme()
		var startzeit time.Time = time.Now()
		var laufzeit time.Duration
		var läufe float64
		for {
			laufzeit = time.Since(startzeit)
			if laufzeit >= time.Second/50 { // Rate alle 50stel Sekunde anpassen
				r.rate = uint64(0.5 + läufe/laufzeit.Seconds()) // Rate ist "Läufe je Sekunde"
				if r.rate < minRate {
					if r.verzoegerung > 0 {
						r.verzoegerung -= time.Millisecond
					}
				}
				if r.rate > maxRate {
					if r.verzoegerung < maxVerzögerung {
						r.verzoegerung += time.Millisecond
					}
				}
				startzeit = time.Now()
				läufe = 0
			}
			select {
			case <-r.stop:
				println("Stoppe", r.name)
				r.laeuft = false
				return
			default:
				r.frun() // Starte Funktion, falls der loop nicht gestoppt wurde ...
				läufe++
				if r.verzoegerung > 0 {
					time.Sleep(time.Duration(r.verzoegerung)) // und warte etwas, damit die Rate konstant bleibt.
				}
			}
		}
	}
	// starte loop
	go runner()
}

// Starte eine Funktion als goroutine.
// Lasse sie so schnell wie möglich loopen und bestimme laufend die Rate je Sekunde.
func (r *routine) Starte() {
	if r.laeuft {
		println("Fehler:", r.name, "läuft bereits.")
		return
	}
	r.laeuft = true
	println("Starte", r.name)
	r.rate = 1e9
	r.verzoegerung = 0
	r.stop = make(chan bool)
	runner := func() {
		defer r.fausnahme()
		var startzeit time.Time = time.Now()
		var laufzeit time.Duration
		var läufe float64
		for {
			laufzeit = time.Since(startzeit)
			if laufzeit >= time.Second/5 { // Rate alle 5tel Sekunde messen:
				r.rate = uint64(0.5 + läufe/laufzeit.Seconds()) // Rate ist "Läufe je Sekunde".
				startzeit = time.Now()
				läufe = 0
			}
			select {
			case <-r.stop:
				println("Stoppe", r.name)
				r.laeuft = false
				return
			default:
				r.frun()
				läufe++
			}
		}
	}
	// starte Prozess
	go runner()
}

// Starte eine Funktion lokal in einer Endlosschleife.
// Lasse sie so schnell wie möglich loopen und bestimme laufend die Rate je Sekunde.
// Hinweis: blockiert, bis sie von außerhalb gestoppt wird.
func (r *routine) LoopeHier() {
	if r.laeuft {
		println("Fehler:", r.name, "läuft bereits.")
		return
	}
	r.laeuft = true
	println("Starte lokal:", r.name)
	r.rate = 1e9
	r.verzoegerung = 0
	r.stop = make(chan bool)
	defer r.fausnahme()
	var startzeit time.Time = time.Now()
	var laufzeit time.Duration
	var läufe float64
	for {
		laufzeit = time.Since(startzeit)
		if laufzeit >= time.Second/5 { // Rate alle 5tel Sekunde messen:
			r.rate = uint64(0.5 + läufe/laufzeit.Seconds()) // Rate ist "Läufe je Sekunde".
			startzeit = time.Now()
			läufe = 0
		}
		select {
		case <-r.stop:
			println("Stoppe", r.name)
			r.laeuft = false
			return
		default:
			r.frun()
			läufe++
		}
	}
}

func (r *routine) GibRate() uint64 { return r.rate }

func (r *routine) GibName() string { return r.name }

func (r *routine) Stoppe() {
	if r.stop == nil {
		println("Fehler:", r.name, "läuft gar nicht.")
		return
	}
	r.stop <- true
	r.stop = nil
}
