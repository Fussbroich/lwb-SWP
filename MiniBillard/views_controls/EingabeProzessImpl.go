package views_controls

import (
	"../hilf"
)

type bpEingabeProzess struct {
	steuerProzess hilf.Prozess
}

func (ctl *bpEingabeProzess) Stoppe() {
	if ctl.steuerProzess != nil {
		ctl.steuerProzess.Stoppe()
	}
}
