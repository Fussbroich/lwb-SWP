package main

import (
	"gfx"

	"./klaenge"
)

func main() {
	gfx.Fenster(100, 100)
	sound := klaenge.BallHitsBallSound()

	läuft := false
	for {
		taste, gedrückt, _ := gfx.TastaturLesen1()
		if gedrückt == 1 {
			switch taste {
			case 's': // quit
				if läuft {
					sound.StoppeLoop()
				} else {
					sound.StarteLoop()
				}
				läuft = !läuft
			case 'q':
				return
			}
		}
	}
}
