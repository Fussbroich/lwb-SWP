package hilf

import "sync"

type smu struct {
	anzLeser   uint
	gA         sync.Mutex
	gALeserAnz sync.Mutex
}

func NewSchreiberMutex() *smu { return &smu{} }

// Viele Leser dürfen jederzeit "rein"
func (m *smu) LeserEin() {
	m.gALeserAnz.Lock()
	m.anzLeser++
	// if m.anzLeser == 1 {
	// 	m.gA.Lock() // keine Schreiber mehr!
	// }
	m.gALeserAnz.Unlock()
}

func (m *smu) LeserAus() {
	m.gALeserAnz.Lock()
	m.anzLeser--
	// if m.anzLeser == 0 {
	// 	m.gA.Unlock()
	// }
	m.gALeserAnz.Unlock()
}

func (m *smu) GibAnzLeser() uint { return m.anzLeser } // Testzwecke

// nur *Ein* Schreiber darf rein, selbst wenn Leser drin sind
func (m *smu) SchreiberEin() {
	m.gA.Lock()
}

func (m *smu) SchreiberAus() { m.gA.Unlock() }
