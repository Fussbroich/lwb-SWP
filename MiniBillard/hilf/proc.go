package hilf

import "time"

type Prozess interface {
	StarteLoop(time.Duration)
	StarteRate(uint8)
	Stoppe()
}

type prozess struct {
	name  string
	frun  func()
	stop  chan bool
	läuft bool // eigentlich überflüssig
}

func NewProzess(name string, f func()) *prozess {
	return &prozess{
		name: name,
		frun: f}
}

func (proc *prozess) StarteLoop(tick time.Duration) {
	if proc.läuft {
		println("Fehler: Prozess", proc.name, "läuft bereits.")
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
	proc.läuft = true
	go runner()
}

func (proc *prozess) StarteRate(rate uint8) {
	if proc.läuft {
		println("Fehler: Prozess", proc.name, "läuft bereits.")
		return
	}
	println("Starte", proc.name, "(max", rate, "Hz)")
	takt := time.NewTicker(time.Second / time.Duration(rate))
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
	proc.läuft = true
	go runner()
}

func (proc *prozess) Stoppe() {
	if !proc.läuft {
		println("Fehler: Prozess", proc.name, "steht bereits.")
		return
	}
	proc.stop <- true
	proc.stop = nil
	proc.läuft = false
}
