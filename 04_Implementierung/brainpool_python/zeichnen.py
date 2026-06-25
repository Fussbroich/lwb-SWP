"""
zeichnen — Zeichenroutinen für Kugeln, Spielfeld und UI-Elemente.

Aufgabe für Schüler:
    Implementiere die Funktionen, die das Spiel auf dem Bildschirm darstellen.
    Jede Funktion bekommt die nötigen Daten als Parameter — kein globaler Zustand.
"""

import gfx
from farben import KUGEL_FARBEN


# =====================================================================
#  Kugel zeichnen
# =====================================================================

def zeichne_kugel(kugel, ox, oy, font_bold):
    """Zeichnet eine einzelne Billardkugel.
    ox, oy: Offset des Spielfelds auf dem Bildschirm.
    font_bold: Pfad zur Schriftart für die Kugelnummer."""
    x, y, r = int(ox + kugel.pos.x), int(oy + kugel.pos.y), int(kugel.radius)
    farbe = KUGEL_FARBEN[min(kugel.wert, len(KUGEL_FARBEN) - 1)]

    if kugel.wert <= 8:
        gfx.Stiftfarbe(*farbe)
        gfx.Vollkreis(x, y, r)
    else:
        # Streifenkugel: weißer Grund + farbiger Streifen
        gfx.Stiftfarbe(252, 253, 242)
        gfx.Vollkreis(x, y, r)
        gfx.Stiftfarbe(*farbe)
        gfx.Vollrechteck(x-int(r*0.75), y-int(r*0.6), int(1.5*r), int(1.2*r))
        gfx.Vollkreissektor(x, y, r, 325, 35)
        gfx.Vollkreissektor(x, y, r, 145, 215)

    # Nummer (nicht bei der weißen Kugel)
    if kugel.wert > 0:
        gfx.Stiftfarbe(252, 253, 242)
        gfx.Vollkreis(x, y, int(r / 2))
        gfx.Stiftfarbe(0, 0, 0)
        gfx.SetzeFont(font_bold, r - 3)
        text = str(kugel.wert)
        tw = gfx.TextBreite("0") * len(text)
        gfx.SchreibeFont(x - tw // 2, y - (r - 3) // 2, text)


# =====================================================================
#  Spielfeld zeichnen
# =====================================================================

def zeichne_spielfeld(spiel, sx, sy, breite, hoehe, schema, font_bold, font_regular):
    """Zeichnet Tuch, Taschen, Kugeln und die Stoßanzeige."""
    rk = spiel.rk

    # Tuch
    gfx.Stiftfarbe(*schema["tuch"])
    gfx.Vollrechteck(sx, sy, breite, hoehe)

    # Diamanten am Rand
    gfx.Stiftfarbe(*schema["diamanten"])
    for i in [1, 2, 3, 5, 6, 7]:
        _diamant(sx + int(i*breite/8), sy - int(rk), int(rk/3))
        _diamant(sx + int(i*breite/8), sy + hoehe + int(rk), int(rk/3))
    for i in [1, 2, 3]:
        _diamant(sx - int(rk), sy + int(i*hoehe/4), int(rk/3))
        _diamant(sx + breite + int(rk), sy + int(i*hoehe/4), int(rk/3))

    # Taschen
    gfx.Stiftfarbe(0, 0, 0)
    for tp in spiel.taschen:
        gfx.Vollkreis(sx + int(tp.x), sy + int(tp.y), int(rk * 1.3))

    # Kugeln
    for kugel in spiel.aktive_kugeln():
        zeichne_kugel(kugel, sx, sy, font_bold)

    # Queue (Stoßanzeige) — nur bei Stillstand
    sk = spiel.spielkugel
    if spiel.stillstand and not sk.eingelocht:
        kraft = spiel.stoss_kraft
        if kraft < 7:     gfx.Stiftfarbe(47, 159, 52)
        elif kraft > 9:   gfx.Stiftfarbe(249, 73, 68)
        else:             gfx.Stiftfarbe(250, 175, 50)
        stoss = spiel.stoss_richtung.mal(kraft * rk)
        x1, y1 = sx + int(sk.pos.x), sy + int(sk.pos.y)
        x2, y2 = x1 + int(stoss.x), y1 + int(stoss.y)
        gfx.Linie(x1, y1, x2, y2, max(2, int(rk / 4)))

        gfx.Stiftfarbe(100, 100, 100)
        gfx.SetzeFont(font_regular, max(8, int(rk * 0.67)))
        gfx.SchreibeFont(x2, y2 - int(2 * rk), f"Stärke: {int(kraft + 0.5)}")


def _diamant(x, y, d):
    hd = d // 2
    gfx.Volldreieck(x-hd, y, x+hd, y, x, y-hd)
    gfx.Volldreieck(x-hd, y, x+hd, y, x, y+hd)


# =====================================================================
#  Anzeigen (Treffer, Fouls, Countdown)
# =====================================================================

def zeichne_anzeige(spiel, sx, sy, breite, g, schema, font_bold):
    """Zeichnet Treffer/Foul-Balken und den Countdown."""
    bx, by = sx - int(g * 1.3), g
    bb, bh = 14 * g, 2 * g

    gfx.Stiftfarbe(*schema["anzeige"])
    gfx.Vollrechteck(bx, by, bb, bh)

    gfx.Stiftfarbe(*schema["text"])
    gfx.SetzeFont(font_bold, g * 3 // 5)
    d = (g - g * 3 // 5) // 2
    gfx.SchreibeFont(bx + d, by + d, "Treffer")
    gfx.SchreibeFont(bx + d, by + g + d, "Fouls")

    anz = max(1, len(spiel.kugeln) - 1)
    label_b = gfx.TextBreite("Treffer") + 2 * d
    balken_x, balken_b = bx + label_b, bb - label_b - d

    gfx.Stiftfarbe(*schema["treffer"])
    gfx.Vollrechteck(balken_x, by+1, min(balken_b, balken_b*spiel.treffer()//anz), g-2)
    gfx.Stiftfarbe(*schema["fouls"])
    gfx.Vollrechteck(balken_x, by+g+1, min(balken_b, balken_b*spiel.strafpunkte//anz), g-2)

    if spiel.treffer() > 0:
        gfx.Stiftfarbe(*schema["fouls"])
        gfx.SchreibeFont(balken_x + d, by + d, str(spiel.treffer()))
    if spiel.strafpunkte > 0:
        gfx.Stiftfarbe(*schema["treffer"])
        gfx.SchreibeFont(balken_x + d, by + g + d, str(spiel.strafpunkte))

    # Countdown rechts
    minuten = int(spiel.restzeit) // 60
    sekunden = int(spiel.restzeit) % 60
    cx = sx + breite - 5 * g
    gfx.Stiftfarbe(*schema["anzeige"])
    gfx.Vollrechteck(cx, by, 6 * g, bh)
    gfx.Stiftfarbe(*schema["text"])
    gfx.SetzeFont(font_bold, int(bh * 0.8))
    gfx.SchreibeFont(cx + g // 2, by + d, f"{minuten:02d}:{sekunden:02d}")


# =====================================================================
#  Button
# =====================================================================

class Button:
    """Ein klickbarer Button mit abgerundeten Ecken."""

    def __init__(self, text, x, y, b, h, aktion, eckradius=0):
        self.text = text
        self.x, self.y, self.b, self.h = x, y, b, h
        self.aktion = aktion
        self.eckradius = eckradius

    def zeichne(self, schema, font_bold):
        gfx.Stiftfarbe(*schema["hintergrund"])
        gfx.Vollrechteck(self.x, self.y, self.b, self.h, self.eckradius)
        gfx.Stiftfarbe(*schema["text"])
        gfx.Rechteck(self.x, self.y, self.b, self.h, self.eckradius)
        gfx.SetzeFont(font_bold, self.h * 2 // 3)
        tw = gfx.TextBreite(self.text)
        gfx.SchreibeFont(self.x + (self.b - tw) // 2, self.y + self.h // 6, self.text)

    def geklickt(self, mx, my):
        return (self.x <= mx <= self.x + self.b and
                self.y <= my <= self.y + self.h)
