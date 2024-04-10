package hilf

import "time"

type Prozess interface {
	StarteLoop(time.Duration)
	StarteRate(uint64)
	Starte()
	GibRate() uint64
	Stoppe()
	Läuft() bool
}

type prozess struct {
	name        string
	frun        func()
	rate        uint64 // Hertz
	verzögerung time.Duration
	stop        chan bool
}

func NewProzess(name string, f func()) *prozess {
	return &prozess{
		name: name,
		frun: f}
}

func (proc *prozess) Läuft() bool { return proc.stop != nil }

func (proc *prozess) StarteLoop(tick time.Duration) {
	proc.rate = 1e9 / uint64(tick.Nanoseconds())
	if proc.stop != nil {
		println("Fehler:", proc.name, "läuft bereits.")
		return
	}
	println("Starte Takt für", proc.name)
	takt := time.NewTicker(tick)
	proc.stop = make(chan bool)
	proc.frun()
	runner := func() {
		defer func() { takt.Stop(); println("Stoppe", proc.name) }()
		for {
			select {
			case <-proc.stop:
				return
			case <-takt.C:
				proc.frun()
			}
		}
	}
	// starte Prozess
	go runner()
}

func (proc *prozess) StarteRate(sollRate uint64) {
	if proc.stop != nil {
		println("Fehler:", proc.name, "läuft bereits.")
		return
	}
	proc.rate = sollRate
	proc.verzögerung = 0
	maxVerzögerung := time.Second / 5
	var minRate, maxRate uint64 = sollRate * 4 / 5, sollRate * 6 / 5
	proc.stop = make(chan bool)
	runner := func() {
		var startzeit time.Time = time.Now()
		var laufzeit time.Duration
		var läufe float64
		for {
			laufzeit = time.Since(startzeit)
			if laufzeit >= time.Second/20 { // Rate alle 20stel Sekunde anpassen
				proc.rate = uint64(läufe / laufzeit.Seconds()) // Rate ist Läufe je Sekunde
				println(proc.name, "läuft mit", proc.rate, "Hz")
				if proc.rate < minRate {
					if proc.verzögerung > 0 {
						proc.verzögerung -= time.Millisecond
					}
				}
				if proc.rate > maxRate {
					if proc.verzögerung < maxVerzögerung {
						proc.verzögerung += time.Millisecond
					}
				}
				startzeit = time.Now()
				läufe = 0
			}
			select {
			case <-proc.stop:
				println("Stoppe", proc.name)
				return
			default:
				proc.frun()
				läufe++
				time.Sleep(time.Duration(proc.verzögerung))
			}
		}
	}
	// starte Prozess
	println("Starte", proc.name, "(soll:", sollRate, "Hz)")
	go runner()
}

func (proc *prozess) Starte() {
	if proc.stop != nil {
		println("Fehler:", proc.name, "läuft bereits.")
		return
	}
	proc.rate = 1e9
	proc.verzögerung = 0
	proc.stop = make(chan bool)
	runner := func() {
		var startzeit time.Time = time.Now()
		var laufzeit time.Duration
		var läufe float64
		for {
			laufzeit = time.Since(startzeit)
			if laufzeit >= time.Second/20 { // Rate alle 20stel Sekunde messen
				proc.rate = uint64(läufe / laufzeit.Seconds()) // Rate ist Läufe je Sekunde
				println("Maus", proc.rate)
				startzeit = time.Now()
				läufe = 0
			}
			select {
			case <-proc.stop:
				println("Stoppe", proc.name)
				return
			default:
				proc.frun()
				läufe++
			}
		}
	}
	// starte Prozess
	println("Starte", proc.name)
	go runner()
}

func (proc *prozess) GibRate() uint64 {
	return uint64(proc.rate)
}

func (proc *prozess) Stoppe() {
	if proc.stop == nil {
		println("Fehler:", proc.name, "läuft gar nicht.")
		return
	}
	proc.stop <- true
	proc.stop = nil
}
