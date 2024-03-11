package hilf

import "time"

type Prozess interface {
	Run()
	StarteLoop()
	StoppeLoop()
}

type prozess struct {
	name string
	run  func()
	stop chan bool
}

func NewProzess(name string, run func()) *prozess {
	return &prozess{
		name: name,
		run:  run}
}

func (proc *prozess) Run() {
	println("Starte", proc.name)
	proc.run()
}

func (proc *prozess) StarteLoop(tick time.Duration) {
	if proc.stop != nil {
		println("Fehler: Prozess", proc.name, "läuft bereits.")
		return
	}
	println("Starte Takt für", proc.name)
	takt := time.NewTicker(tick)
	proc.stop = make(chan bool)
	proc.run()
	runner := func() {
		defer func() { takt.Stop(); println("Stoppe", proc.name) }()
		for {
			select {
			case <-proc.stop:
				return
			case <-takt.C:
				proc.run()
			}
		}
	}
	// starte Prozess
	go runner()
}

func (proc *prozess) StoppeLoop() {
	proc.stop <- true
	proc.stop = nil
}
