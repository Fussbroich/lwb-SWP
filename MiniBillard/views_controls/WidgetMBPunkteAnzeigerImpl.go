package views_controls

import (
	"fmt"
	"gfx"
	"math"

	"../modelle"
)

type miniBSpielinfo struct {
	billard modelle.MiniBillardSpiel
	widget
}

func NewMBPunkteAnzeiger(billard modelle.MiniBillardSpiel) *miniBSpielinfo {
	return &miniBSpielinfo{billard: billard, widget: *NewFenster()}
}

func (f *miniBSpielinfo) Zeichne() {
	if !f.IstAktiv() {
		return
	}
	f.widget.Zeichne()
	schreiber := f.liberationMonoBoldSchreiber()
	breite, höhe := f.GibGroesse()
	ra := höhe / 2

	var zeilenhöhe uint16
	// dieses Widget zeigt Treffer und Strafen an
	var anzKugeln uint8 = uint8(len(f.billard.GibKugeln()) - 1)
	tr, st := f.billard.GibTreffer(), math.Min(float64(anzKugeln), float64(f.billard.GibStrafpunkte()))
	zeilenhöhe = höhe / 2
	schreiber.SetzeSchriftgroesse(int(zeilenhöhe) * 3 / 5)
	d := (zeilenhöhe - uint16(schreiber.GibSchriftgroesse())) / 2
	var bBalken, xSBalken uint16 = breite - 2*ra - 5*uint16(schreiber.GibSchriftgroesse()), f.startX + 2*ra + 5*uint16(schreiber.GibSchriftgroesse())

	f.stiftfarbe(f.hg)
	f.vollKreisGFX(ra, ra, ra)
	f.vollRechteckGFX(ra, 0, xSBalken-f.startX-ra, höhe)
	f.stiftfarbe(f.vg)

	schreiber.Schreibe(f.startX+2*ra+d, f.startY+d, "Treffer")
	schreiber.Schreibe(f.startX+2*ra+d, f.startY+zeilenhöhe+d, "Fouls")

	// zeichne beide Fortschritts-Balken
	f.stiftfarbe(gibFarbe(FanzTreffer())) // Treffer gelb
	gfx.Vollrechteck(xSBalken, f.startY+1, bBalken*uint16(tr)/uint16(anzKugeln), zeilenhöhe-2)
	f.stiftfarbe(gibFarbe(FanzFouls())) // Fouls in Warnrot
	gfx.Vollrechteck(xSBalken, f.startY+zeilenhöhe+1, bBalken*uint16(st)/uint16(anzKugeln), zeilenhöhe-2)

	// Kreis links zeigt Treffer an
	schreiber.SetzeSchriftgroesse(int(ra * 4 / 3))
	f.stiftfarbe(f.vg)
	f.kreisGFX(ra, ra, ra)
	var x, y uint16
	if tr > 9 {
		x = ra + f.startX - uint16(schreiber.GibSchriftgroesse())*2/5
		y = ra + f.startY - uint16(schreiber.GibSchriftgroesse())/2
	} else {
		x = ra + f.startX - uint16(schreiber.GibSchriftgroesse())/4
		y = ra + f.startY - uint16(schreiber.GibSchriftgroesse())/2
	}
	schreiber.Schreibe(x, y, fmt.Sprintf("%d", tr))
}
