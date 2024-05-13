// Autoren:
//
//	Thomas Schrader
//	Bettina Chang
//
//	Zweck:
//		Das Spielprogramm BrainPool -
//		ein Softwareprojekt im Rahmen der Lehrerweiterbildung Berlin
//
//	Notwendige Software: Linux, Go ab 1.18
//		es läuft auch unter Windows, jedoch wird der gfx-Server
//		sehr gefordert.
//	verwendete Pakete:
//		gfx, fmt, math, math/rand, strconv, strings, unicode/utf8, time,
//		runtime, os, errors, path/filepath, encoding/csv
//	Notwendige Hardware:
//		PC, Bildschirm, Tastatur, Maus mit Scrollrad
//
//	Datum: 01.05.2024
package main

import (
	"./apps"
)

// ####### der Startpunkt ##################################################
func main() {
	// Die gewünschte Fensterbreite in Pixeln wird übergeben.
	// Das Seitenverhältnis des Spiels ist B:H = 16:11
	app := apps.NewBPApp(1024)
	//app := apps.NewBeispielApp(1024)
	apps.RunApp(app)
}
