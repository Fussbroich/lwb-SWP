package views_controls

import (
	"fmt"
	"gfx"

	"../modelle"
)

type miniBSpielinfo struct {
	billard                     modelle.MiniBillardSpiel
	anzKugeln, treffer, strafen uint16
	widget
}

func NewMBPunkteAnzeiger(billard modelle.MiniBillardSpiel) *miniBSpielinfo {
	return &miniBSpielinfo{billard: billard, widget: *NewFenster()}
}

func uMin16(a, b uint16) uint16 {
	if a <= b {
		return a
	}
	return b
}

func (f *miniBSpielinfo) Update() {
	anz := uint16(len(f.billard.GibKugeln()) - 1)
	tr, st := uint16(f.billard.GibTreffer()), uint16(f.billard.GibStrafpunkte())
	f.veraltet = anz != f.anzKugeln || tr != f.treffer || st != f.strafen
	if !f.veraltet {
		return
	}
	f.anzKugeln, f.treffer, f.strafen = anz, tr, st
}

func (f *miniBSpielinfo) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	schreiber := f.monoBoldSchreiber()
	breite, höhe := f.GibGroesse()
	ra := höhe / 2

	// dieses Widget zeigt Treffer und Strafen an
	var zeilenhöhe uint16 = höhe / 2
	schreiber.SetzeSchriftgroesse(int(zeilenhöhe) * 3 / 5)
	d := (zeilenhöhe - uint16(schreiber.GibSchriftgroesse())) / 2
	var bBalken uint16 = breite - 2*ra - 5*uint16(schreiber.GibSchriftgroesse())
	var xSBalken uint16 = f.startX + 2*ra + 5*uint16(schreiber.GibSchriftgroesse())

	f.stiftfarbe(f.hg)
	f.vollKreisGFX(ra, ra, ra)
	f.vollRechteckGFX(ra, 0, xSBalken-f.startX-ra, höhe)

	f.stiftfarbe(f.vg)
	schreiber.Schreibe(f.startX+2*ra+d, f.startY+d, "Treffer")
	schreiber.Schreibe(f.startX+2*ra+d, f.startY+zeilenhöhe+d, "Fouls")

	// zeichne beide Fortschritts-Balken
	f.stiftfarbe(gibFarbe(FanzTreffer())) // Treffer
	gfx.Vollrechteck(xSBalken, f.startY+1, uMin16(bBalken*f.treffer/f.anzKugeln, bBalken), zeilenhöhe-2)
	f.stiftfarbe(gibFarbe(FanzFouls())) // Fouls
	gfx.Vollrechteck(xSBalken, f.startY+zeilenhöhe+1, uMin16(bBalken*f.strafen/f.anzKugeln, bBalken), zeilenhöhe-2)
	if f.treffer > 0 {
		f.stiftfarbe(gibFarbe(FanzFouls()))
		schreiber.Schreibe(xSBalken+d, f.startY+d, fmt.Sprint(f.treffer))
	}
	if f.strafen > 0 {
		f.stiftfarbe(gibFarbe(FanzTreffer()))
		schreiber.Schreibe(xSBalken+d, f.startY+zeilenhöhe+d, fmt.Sprint(f.strafen))
	}

	// Kreis links zeigt Treffer an
	schreiber.SetzeSchriftgroesse(int(ra * 4 / 3))
	f.stiftfarbe(f.vg)
	f.kreisGFX(ra, ra, ra)
	var x, y uint16
	if f.treffer > 9 {
		x = ra + f.startX - uint16(schreiber.GibSchriftgroesse())*2/5
		y = ra + f.startY - uint16(schreiber.GibSchriftgroesse())/2
	} else {
		x = ra + f.startX - uint16(schreiber.GibSchriftgroesse())/4
		y = ra + f.startY - uint16(schreiber.GibSchriftgroesse())/2
	}
	schreiber.Schreibe(x, y, fmt.Sprintf("%d", f.treffer))
}
