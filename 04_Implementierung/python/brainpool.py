"""
BrainPool — Das Mini-Billard für Schlaue.
Ein 9-Ball-Billardspiel mit Quiz-Modus.

Steuerung:
    Maus bewegen   → Zielrichtung ändern
    Links klicken  → Stoßen
    Scrollrad      → Stoßkraft ändern
    H              → Hilfe
    N              → Neues Spiel
    M              → Musik an
    D              → Dunkel/Hell
    S oder Q       → Beenden

Benötigt: gfx.py (im selben Ordner) und pip install pygame
"""

import gfx
import math
import random
import csv
import os
import time

# =====================================================================
#  Pfade zu den Spiel-Dateien
# =====================================================================

ASSET_DIR = os.path.join(os.path.dirname(__file__), "..", "assets")
FONT_BOLD = os.path.join(ASSET_DIR, "fontfiles", "LiberationMono-Bold.ttf")
FONT_REGULAR = os.path.join(ASSET_DIR, "fontfiles", "LiberationMono-Regular.ttf")
SOUND_BALL_BALL = os.path.join(ASSET_DIR, "soundfiles", "ballHitsBall.wav")
SOUND_BALL_RAIL = os.path.join(ASSET_DIR, "soundfiles", "ballHitsRail.wav")
SOUND_BALL_POCKET = os.path.join(ASSET_DIR, "soundfiles", "ballIntoPocket.wav")
SOUND_CUE = os.path.join(ASSET_DIR, "soundfiles", "cueHitsBall.wav")
SOUND_AMBIENCE = os.path.join(ASSET_DIR, "soundfiles", "billardPubAmbience.wav")
SOUND_MUSIK = os.path.join(ASSET_DIR, "soundfiles", "coolJazzLoop2641.wav")
QUIZ_CSV = os.path.join(ASSET_DIR, "quizfragen", "InformatiksystemQuiz.csv")


# =====================================================================
#  Farbschema — alle Farben an einem Ort, umschaltbar Hell/Dunkel
# =====================================================================

HELL = {
    "hintergrund":  (225, 232, 236),
    "text":         (1, 88, 122),
    "bande":        (1, 88, 122),
    "anzeige":      (255, 255, 255),
    "tuch":         (92, 179, 193),
    "diamanten":    (180, 230, 255),
    "treffer":      (243, 186, 0),
    "fouls":        (219, 80, 0),
    "quiz":         (225, 232, 236),
    "quiz_a0":      (243, 186, 0),
    "quiz_a1":      (92, 179, 193),
    "quiz_a2":      (92, 179, 193),
    "quiz_a3":      (243, 186, 0),
}

DUNKEL = {
    "hintergrund":  (25, 32, 36),
    "text":         (125, 170, 195),
    "bande":        (1, 88, 122),
    "anzeige":      (25, 32, 36),
    "tuch":         (92, 179, 193),
    "diamanten":    (180, 230, 255),
    "treffer":      (143, 86, 0),
    "fouls":        (119, 40, 0),
    "quiz":         (25, 32, 36),
    "quiz_a0":      (80, 64, 19),
    "quiz_a1":      (40, 61, 65),
    "quiz_a2":      (40, 61, 65),
    "quiz_a3":      (80, 64, 19),
}

farben = HELL   # aktives Farbschema


# =====================================================================
#  Vec2 — Ein einfacher 2D-Vektor
# =====================================================================

class Vec2:
    """Beschreibt eine Position oder Richtung in 2D (x, y).

    Verwendung:
        position = Vec2(100, 200)
        richtung = Vec2(1, 0)           # nach rechts
        neue_pos = position.plus(richtung.mal(5))
    """

    def __init__(self, x=0.0, y=0.0):
        self.x = float(x)
        self.y = float(y)

    def plus(self, other):
        return Vec2(self.x + other.x, self.y + other.y)

    def minus(self, other):
        return Vec2(self.x - other.x, self.y - other.y)

    def mal(self, skalar):
        return Vec2(self.x * skalar, self.y * skalar)

    def betrag(self):
        return math.sqrt(self.x ** 2 + self.y ** 2)

    def normiert(self):
        b = self.betrag()
        return Vec2(self.x / b, self.y / b) if b > 0 else Vec2()

    def ist_null(self):
        return self.x == 0 and self.y == 0

    def skalarprodukt(self, other):
        return self.x * other.x + self.y * other.y

    def projiziert_auf(self, other):
        b2 = other.skalarprodukt(other)
        if b2 == 0:
            return Vec2()
        return other.mal(self.skalarprodukt(other) / b2)


# =====================================================================
#  Kugel — Eine Billardkugel
# =====================================================================

class Kugel:
    """Eine Billardkugel mit Position, Geschwindigkeit und Nummer.
    wert 0 = weiße Spielkugel, 1-8 = Vollfarben, 9 = Streifen."""

    def __init__(self, pos, radius, wert):
        self.pos = pos
        self.v = Vec2()
        self.radius = radius
        self.wert = wert
        self.eingelocht = False
        self._kollision_mit = None

    def bewegen(self):
        self.pos = self.pos.plus(self.v)
        geschwindigkeit = self.v.betrag()
        if geschwindigkeit > 0.15:
            self.v = self.v.mal(1 - 0.02 / geschwindigkeit)
        else:
            self.v = Vec2()

    def pruefe_bande(self, breite, hoehe, bei_kollision):
        if self.eingelocht:
            return
        vx, vy = self.v.x, self.v.y
        x, y, r = self.pos.x, self.pos.y, self.radius
        trifft = False
        am_rand = not (r <= x <= breite - r and r <= y <= hoehe - r)
        if not am_rand and x + vx < r:
            vx = -vx; trifft = True
        if not am_rand and x + vx > breite - r:
            vx = -vx; trifft = True
        if not am_rand and y + vy < r:
            vy = -vy; trifft = True
        if not am_rand and y + vy > hoehe - r:
            vy = -vy; trifft = True
        if trifft:
            bei_kollision()
            self.v = Vec2(vx, vy)

    def pruefe_kollision(self, andere, bei_kollision):
        if self._kollision_mit is andere:
            return
        v1, v2 = self.v, andere.v
        dist_vor = andere.pos.plus(v2).minus(self.pos.plus(v1))
        dist_jetzt = andere.pos.minus(self.pos)
        if dist_vor.betrag() > self.radius + andere.radius:
            return
        ueberlappen = dist_jetzt.betrag() < self.radius + andere.radius
        n = dist_vor.normiert()
        v1p, v1o = v1.projiziert_auf(n), v1.minus(v1.projiziert_auf(n))
        v2p, v2o = v2.projiziert_auf(n), v2.minus(v2.projiziert_auf(n))
        if ueberlappen:
            if dist_vor.betrag() < dist_jetzt.betrag():
                u1, u2 = v2p.plus(v1o), v1p.plus(v2o)
            else:
                u1, u2 = v1, v2
        else:
            u1, u2 = v2p.plus(v1o), v1p.plus(v2o)
            bei_kollision()
        self.v, andere.v = u1, u2
        self._kollision_mit = andere
        andere._kollision_mit = self

    def kopie(self):
        k = Kugel(Vec2(self.pos.x, self.pos.y), self.radius, self.wert)
        k.eingelocht = self.eingelocht
        return k


# =====================================================================
#  BillardSpiel — Physik und Spielregeln
# =====================================================================

KUGEL_FARBEN = [
    (252, 253, 242), (255, 201, 78), (34, 88, 175), (249, 73, 68),
    (84, 73, 149), (255, 139, 33), (47, 159, 52), (155, 53, 30),
    (48, 49, 54), (255, 201, 78),
]


class BillardSpiel:
    """Ein 9-Ball-Billardspiel mit Physiksimulation."""

    def __init__(self, breite, hoehe):
        self.breite, self.hoehe = float(breite), float(hoehe)
        self.rk = self.breite * 57.2 / 2540
        self.kugeln = []
        self.spielkugel = None
        self.eingelochte = []
        self.strafpunkte = 0
        self.stoss_richtung = Vec2(1, 0)
        self.stoss_kraft = 5.0
        self.stillstand = True
        self._angestossen = False
        self._angespielte = None
        self._vorige_kugeln = []
        self._orig_kugeln = []
        self.spielzeit = 4 * 60.0
        self.restzeit = self.spielzeit
        self._letzte_zeit = time.time()
        self._laeuft = False

        rt, rtm = 1.9 * self.rk, 1.5 * self.rk
        b, h = self.breite, self.hoehe
        self.taschen = [Vec2(0,0), Vec2(0,h), Vec2(b/2,h),
                        Vec2(b,h), Vec2(b,0), Vec2(b/2,0)]
        self.taschen_radius = [rt, rt, rtm, rt, rt, rtm]
        self._aufstellen_9ball()

    def _aufstellen_9ball(self):
        rk, b, h = self.rk, self.breite, self.hoehe
        dx, dy = 0.866 * (2*rk+1), 0.5 * (2*rk+1)
        p1 = Vec2(b/4, h/2)
        pos = {1: p1,
               2: p1.plus(Vec2(-dx,-dy)), 3: p1.plus(Vec2(-dx,dy)),
               4: p1.plus(Vec2(-2*dx,-2*dy)), 5: p1.plus(Vec2(-2*dx,2*dy)),
               6: p1.plus(Vec2(-3*dx,-dy)), 7: p1.plus(Vec2(-3*dx,dy)),
               8: p1.plus(Vec2(-4*dx,0)), 9: p1.plus(Vec2(-2*dx,0))}
        self.kugeln = [Kugel(Vec2(3*b/4, h/2), rk, 0)]
        for w, p in pos.items():
            self.kugeln.append(Kugel(p, rk, w))
        self.spielkugel = self.kugeln[0]
        self._spielkugel_neusetzen()
        self._orig_kugeln = [k.kopie() for k in self.kugeln]

    def _spielkugel_neusetzen(self):
        b, h, rk = self.breite, self.hoehe, self.rk
        self.spielkugel.pos = Vec2(
            b - rk - random.random()*(b/4 - 2*rk),
            random.random()*(h - 2*rk) + rk)

    def starte(self):
        self._laeuft = True
        self._letzte_zeit = time.time()

    def stoppe(self):
        self._laeuft = False

    def laeuft(self):
        return self._laeuft

    def reset(self):
        self.kugeln = [k.kopie() for k in self._orig_kugeln]
        self.spielkugel = self.kugeln[0]
        self._spielkugel_neusetzen()
        self.eingelochte, self.strafpunkte = [], 0
        self.stillstand, self._angestossen = True, False
        self.restzeit = self.spielzeit
        self._letzte_zeit = time.time()

    def stosse(self):
        if not self.stillstand or self.stoss_kraft <= 0.15:
            return
        if self.stoss_richtung.ist_null():
            return
        gfx.SpieleSound(SOUND_CUE)
        self._vorige_kugeln = [k.kopie() for k in self.kugeln]
        self._angestossen, self._angespielte = True, None
        self.spielkugel.v = self.stoss_richtung.mal(self.stoss_kraft)
        self.stoss_kraft = 5.0
        self.stillstand = False

    def aktualisiere(self):
        if not self._laeuft:
            return
        jetzt = time.time()
        self.restzeit = max(0, self.restzeit - (jetzt - self._letzte_zeit))
        self._letzte_zeit = jetzt
        aktive = [k for k in self.kugeln if not k.eingelocht]
        for i, k1 in enumerate(aktive):
            for k2 in aktive[i+1:]:
                k1.pruefe_kollision(k2, lambda: (
                    gfx.SpieleSound(SOUND_BALL_BALL),
                    self._notiere_beruehrt(k1, k2)))
            k1._kollision_mit = None
            k1.pruefe_bande(self.breite, self.hoehe,
                            lambda: gfx.SpieleSound(SOUND_BALL_RAIL))
            k1.bewegen()
            for tp, tr in zip(self.taschen, self.taschen_radius):
                if tp.minus(k1.pos).betrag() < tr:
                    gfx.SpieleSound(SOUND_BALL_POCKET)
                    k1.eingelocht, k1.v = True, Vec2()
                    if k1 is not self.spielkugel and k1 not in self.eingelochte:
                        self.eingelochte.append(k1)
                    break
        if all(k.v.ist_null() for k in self.kugeln):
            self.stillstand = True
        if self._angestossen and self.stillstand:
            self._angestossen = False
            if self._ist_foul():
                self.strafpunkte += 1
            if self.spielkugel.eingelocht:
                self._stoss_wiederholen()

    def _notiere_beruehrt(self, k1, k2):
        if self._angespielte is None:
            if k1 is self.spielkugel:   self._angespielte = k2
            elif k2 is self.spielkugel: self._angespielte = k1

    def _ist_foul(self):
        if self.spielkugel.eingelocht or self._angespielte is None:
            return True
        return len(self.eingelochte) <= sum(1 for k in self._vorige_kugeln if k.eingelocht)

    def _stoss_wiederholen(self):
        self.kugeln = [k.kopie() for k in self._vorige_kugeln]
        self.spielkugel, self.stillstand = self.kugeln[0], True

    def aktive_kugeln(self):
        return [k for k in self.kugeln if not k.eingelocht]

    def treffer(self):
        return len(self.eingelochte)

    def alle_eingelocht(self):
        return all(k.eingelocht or k is self.spielkugel for k in self.kugeln)


# =====================================================================
#  Quiz
# =====================================================================

class QuizFrage:
    def __init__(self, frage, antworten, richtige_index):
        self.frage = frage
        self.antworten = antworten
        self.richtige_index = richtige_index

    def ist_richtig(self, index):
        return index == self.richtige_index


def lade_quiz(csv_pfad):
    fragen = []
    with open(csv_pfad, encoding="utf-8") as f:
        for z in csv.reader(f, delimiter=";"):
            if len(z) == 6:
                fragen.append(QuizFrage(z[0], z[1:5], int(z[5])))
    random.shuffle(fragen)
    return fragen


# =====================================================================
#  Zeichenfunktionen
# =====================================================================

def zeichne_kugel(kugel, ox, oy):
    x, y, r = int(ox + kugel.pos.x), int(oy + kugel.pos.y), int(kugel.radius)
    farbe = KUGEL_FARBEN[min(kugel.wert, len(KUGEL_FARBEN)-1)]

    if kugel.wert <= 8:
        gfx.Stiftfarbe(*farbe)
        gfx.Vollkreis(x, y, r)
    else:
        gfx.Stiftfarbe(252, 253, 242)
        gfx.Vollkreis(x, y, r)
        gfx.Stiftfarbe(*farbe)
        gfx.Vollrechteck(x-int(r*0.75), y-int(r*0.6), int(1.5*r), int(1.2*r))
        gfx.Vollkreissektor(x, y, r, 325, 35)
        gfx.Vollkreissektor(x, y, r, 145, 215)

    if kugel.wert > 0:
        gfx.Stiftfarbe(252, 253, 242)
        gfx.Vollkreis(x, y, int(r/2))
        gfx.Stiftfarbe(0, 0, 0)
        gfx.SetzeFont(FONT_BOLD, r-3)
        text = str(kugel.wert)
        tw = gfx.TextBreite("0") * len(text)
        gfx.SchreibeFont(x - tw//2, y - (r-3)//2, text)


def zeichne_spielfeld(spiel, sx, sy, breite, hoehe):
    # Tuch
    gfx.Stiftfarbe(*farben["tuch"])
    gfx.Vollrechteck(sx, sy, breite, hoehe)

    rk = spiel.rk

    # Diamanten
    gfx.Stiftfarbe(*farben["diamanten"])
    for i in [1, 2, 3, 5, 6, 7]:
        _diamant(sx + int(i*breite/8), sy - int(rk), int(rk/3))
        _diamant(sx + int(i*breite/8), sy + hoehe + int(rk), int(rk/3))
    for i in [1, 2, 3]:
        _diamant(sx - int(rk), sy + int(i*hoehe/4), int(rk/3))
        _diamant(sx + breite + int(rk), sy + int(i*hoehe/4), int(rk/3))

    # Taschen
    gfx.Stiftfarbe(0, 0, 0)
    for tp in spiel.taschen:
        gfx.Vollkreis(sx + int(tp.x), sy + int(tp.y), int(rk*1.3))

    # Kugeln
    for kugel in spiel.aktive_kugeln():
        zeichne_kugel(kugel, sx, sy)

    # Queue (Stoßanzeige) — nur bei Stillstand
    sk = spiel.spielkugel
    if spiel.stillstand and not sk.eingelocht:
        kraft = spiel.stoss_kraft
        if kraft < 7:     qfarbe = (47, 159, 52)
        elif kraft > 9:   qfarbe = (249, 73, 68)
        else:             qfarbe = (250, 175, 50)
        gfx.Stiftfarbe(*qfarbe)
        stoss = spiel.stoss_richtung.mal(kraft * rk)
        x1, y1 = sx + int(sk.pos.x), sy + int(sk.pos.y)
        x2, y2 = x1 + int(stoss.x), y1 + int(stoss.y)
        gfx.Linie(x1, y1, x2, y2, max(2, int(rk / 4)))

        gfx.Stiftfarbe(100, 100, 100)
        gfx.SetzeFont(FONT_REGULAR, max(8, int(rk*0.67)))
        gfx.SchreibeFont(x2, y2 - int(2*rk), f"Stärke: {int(kraft+0.5)}")


def zeichne_anzeige(spiel, sx, sy, breite, hoehe, g):
    """Zeichnet Treffer/Foul-Balken und Countdown."""
    # Hintergrundbox für Treffer/Fouls
    bx, by = sx - int(g*1.3), g
    bb, bh = 14*g, 2*g
    gfx.Stiftfarbe(*farben["anzeige"])
    gfx.Vollrechteck(bx, by, bb, bh)

    # Labels
    gfx.Stiftfarbe(*farben["text"])
    gfx.SetzeFont(FONT_BOLD, g*3//5)
    d = (g - g*3//5) // 2
    gfx.SchreibeFont(bx + d, by + d, "Treffer")
    gfx.SchreibeFont(bx + d, by + g + d, "Fouls")

    # Fortschrittsbalken
    anz = max(1, len(spiel.kugeln) - 1)
    label_b = gfx.TextBreite("Treffer") + 2 * d
    balken_x = bx + label_b
    balken_b = bb - label_b - d

    treffer_b = min(balken_b, balken_b * spiel.treffer() // anz)
    gfx.Stiftfarbe(*farben["treffer"])
    gfx.Vollrechteck(balken_x, by + 1, treffer_b, g - 2)

    foul_b = min(balken_b, balken_b * spiel.strafpunkte // anz)
    gfx.Stiftfarbe(*farben["fouls"])
    gfx.Vollrechteck(balken_x, by + g + 1, foul_b, g - 2)

    # Zahlenwerte auf den Balken
    if spiel.treffer() > 0:
        gfx.Stiftfarbe(*farben["fouls"])
        gfx.SchreibeFont(balken_x + d, by + d, str(spiel.treffer()))
    if spiel.strafpunkte > 0:
        gfx.Stiftfarbe(*farben["treffer"])
        gfx.SchreibeFont(balken_x + d, by + g + d, str(spiel.strafpunkte))

    # Countdown (rechts)
    minuten = int(spiel.restzeit) // 60
    sekunden = int(spiel.restzeit) % 60
    gfx.Stiftfarbe(*farben["anzeige"])
    cx = sx + breite - 5*g
    gfx.Vollrechteck(cx, by, 6*g, bh)
    gfx.Stiftfarbe(*farben["text"])
    gfx.SetzeFont(FONT_BOLD, int(bh * 0.8))
    gfx.SchreibeFont(cx + g//2, by + d, f"{minuten:02d}:{sekunden:02d}")


def _diamant(x, y, d):
    hd = d // 2
    gfx.Volldreieck(x-hd, y, x+hd, y, x, y-hd)
    gfx.Volldreieck(x-hd, y, x+hd, y, x, y+hd)


# =====================================================================
#  Buttons — klickbar mit Maus und Tastatur
# =====================================================================

class Button:
    """Ein einfacher klickbarer Button mit abgerundeten Ecken."""
    def __init__(self, text, x, y, b, h, aktion, eckradius=0):
        self.text = text
        self.x, self.y, self.b, self.h = x, y, b, h
        self.aktion = aktion
        self.eckradius = eckradius

    def zeichne(self):
        gfx.Stiftfarbe(*farben["hintergrund"])
        gfx.Vollrechteck(self.x, self.y, self.b, self.h, self.eckradius)
        gfx.Stiftfarbe(*farben["text"])
        gfx.Rechteck(self.x, self.y, self.b, self.h, self.eckradius)
        gfx.SetzeFont(FONT_BOLD, self.h * 2 // 3)
        tw = gfx.TextBreite(self.text)
        tx = self.x + (self.b - tw) // 2
        gfx.SchreibeFont(tx, self.y + self.h // 6, self.text)

    def geklickt(self, mx, my):
        return (self.x <= mx <= self.x + self.b and
                self.y <= my <= self.y + self.h)


# =====================================================================
#  Hauptprogramm
# =====================================================================

def main():
    global farben

    # ---------- Layout ----------
    FENSTER_BREITE = 1024
    g = FENSTER_BREITE // 32
    FENSTER_HOEHE = 22 * g
    SF_BREITE = 3 * FENSTER_BREITE // 4
    SF_HOEHE = SF_BREITE // 2
    SX, SY = 4 * g, 6 * g
    BANDE = int(g + g / 3)

    # ---------- Init ----------
    gfx.Fenster(FENSTER_BREITE, FENSTER_HOEHE, "BrainPool — Das Mini-Billard")
    spiel = BillardSpiel(SF_BREITE, SF_HOEHE)
    spiel.starte()

    for pfad in [SOUND_BALL_BALL, SOUND_BALL_RAIL, SOUND_BALL_POCKET, SOUND_CUE]:
        gfx.LadeSound(pfad)

    # Ambience starten
    if os.path.exists(SOUND_AMBIENCE):
        gfx.SpieleMusik(SOUND_AMBIENCE)

    fragen = lade_quiz(QUIZ_CSV) if os.path.exists(QUIZ_CSV) else []
    aktuelle_frage = fragen[0] if fragen else None
    frage_index = 0
    modus = "spiel"
    darkmode = False

    # ---------- Aktionen ----------
    def hilfe_toggle():
        nonlocal modus
        if modus == "hilfe":
            modus = "spiel"; spiel.starte()
        else:
            modus = "hilfe"; spiel.stoppe()

    def neues_spiel():
        nonlocal modus
        spiel.reset(); spiel.starte(); modus = "spiel"

    def musik_toggle():
        if os.path.exists(SOUND_MUSIK):
            gfx.SpieleMusik(SOUND_MUSIK)

    def darkmode_toggle():
        global farben
        nonlocal darkmode
        darkmode = not darkmode
        farben = DUNKEL if darkmode else HELL

    def beenden():
        nonlocal modus
        gfx.StoppeMusik()
        modus = "ende"

    # ---------- Buttons (gleichmäßig verteilt, abgerundet) ----------
    btn_y = SY + SF_HOEHE + 5 * g // 2
    btn_h = g * 3 // 4
    eck = g // 3
    n_btn = 5
    gesamt_b = SF_BREITE + 2 * BANDE - 2 * g
    btn_b = gesamt_b // n_btn - g // 2
    btn_start = SX - BANDE + g
    btn_abstand = gesamt_b // n_btn
    buttons = [
        Button("(H)ilfe",     btn_start,                btn_y, btn_b, btn_h, hilfe_toggle, eck),
        Button("(N)eu",       btn_start + btn_abstand,  btn_y, btn_b, btn_h, neues_spiel, eck),
        Button("(M)usik",     btn_start + 2*btn_abstand,btn_y, btn_b, btn_h, musik_toggle, eck),
        Button("(D)unkel",    btn_start + 3*btn_abstand,btn_y, btn_b, btn_h, darkmode_toggle, eck),
        Button("(S)chließen", btn_start + 4*btn_abstand,btn_y, btn_b, btn_h, beenden, eck),
    ]

    # ---------- Game-Loop ----------
    while gfx.fenster_offen():

        # ====== Eingabe ======
        mx, my = gfx.MausPosition()

        # Tasten
        if gfx.TasteGedrueckt('s') or gfx.TasteGedrueckt('q'): beenden()
        if modus == "ende": break
        if gfx.TasteGedrueckt('h'): hilfe_toggle()
        if gfx.TasteGedrueckt('n'): neues_spiel()
        if gfx.TasteGedrueckt('m'): musik_toggle()
        if gfx.TasteGedrueckt('d'): darkmode_toggle()

        # Button-Klicks
        if gfx.MausLosgelassen():
            for btn in buttons:
                if btn.geklickt(mx, my):
                    btn.aktion()
                    break
        if modus == "ende": break

        # Spielsteuerung — Schuss nur innerhalb des Spielfelds
        im_spielfeld = (SX <= mx <= SX + SF_BREITE and SY <= my <= SY + SF_HOEHE)

        if modus == "spiel" and spiel.laeuft() and spiel.stillstand:
            if not spiel.spielkugel.eingelocht:
                kx, ky = SX + spiel.spielkugel.pos.x, SY + spiel.spielkugel.pos.y
                r = Vec2(mx - kx, my - ky)
                if r.betrag() > 0:
                    spiel.stoss_richtung = r.normiert()
            if gfx.MausGeklickt() and im_spielfeld:
                spiel.stosse()
            scroll = gfx.MausScroll()
            if scroll:
                spiel.stoss_kraft = max(0, min(14, spiel.stoss_kraft + scroll))

        # Quiz-Klick (gleiche Geometrie wie beim Zeichnen)
        if modus == "quiz" and aktuelle_frage and gfx.MausLosgelassen():
            ix, iy = SX - BANDE + 2 + oeck, SY - BANDE + 2 + oeck
            ib = SF_BREITE + 2*BANDE - 4 - 2*oeck
            ih = SF_HOEHE + 2*BANDE - 4 - 2*oeck
            a_top = iy + ih * 2 // 8
            a_rest = iy + ih - a_top
            d = 4
            for i in range(4):
                ax = ix + (i%2) * (ib//2) + d
                ay = a_top + (i//2) * (a_rest//2) + d
                ab2 = ib//2 - 2*d
                ah2 = a_rest//2 - 2*d
                if ax <= mx <= ax + ab2 and ay <= my <= ay + ah2:
                    if aktuelle_frage.ist_richtig(i):
                        spiel.strafpunkte = max(0, spiel.strafpunkte - 1)
                    frage_index = (frage_index + 1) % len(fragen)
                    aktuelle_frage = fragen[frage_index]
                    break

        # ====== Logik ======
        if modus == "spiel":
            spiel.aktualisiere()
            if spiel.strafpunkte > spiel.treffer() + 2 and fragen:
                spiel.stoppe(); modus = "quiz"
            elif spiel.restzeit <= 0:
                spiel.stoppe(); modus = "gameover"

        if modus == "quiz":
            if spiel.strafpunkte == 0 or spiel.strafpunkte < spiel.treffer():
                modus = "spiel"; spiel.starte()

        # ====== Zeichnen ======
        gfx.Cls(*farben["hintergrund"])

        # Bande (abgerundete Ecken)
        gfx.Stiftfarbe(*farben["bande"])
        gfx.Vollrechteck(SX - BANDE, SY - BANDE,
                         SF_BREITE + 2*BANDE, SF_HOEHE + 2*BANDE, BANDE)

        # Spielfeld
        zeichne_spielfeld(spiel, SX, SY, SF_BREITE, SF_HOEHE)

        # Anzeige (Treffer/Fouls/Zeit)
        zeichne_anzeige(spiel, SX, SY, SF_BREITE, SF_HOEHE, g)

        # Buttons
        for btn in buttons:
            btn.zeichne()

        # Overlays
        ox, oy = SX - BANDE + 2, SY - BANDE + 2
        ob, oh = SF_BREITE + 2*BANDE - 4, SF_HOEHE + 2*BANDE - 4
        oeck = BANDE - 2

        if modus == "quiz" and aktuelle_frage:
            gfx.Stiftfarbe(*farben["quiz"])
            gfx.Vollrechteck(ox, oy, ob, oh, oeck)
            # Inhaltsbereich: um Eckradius eingerückt
            ix, iy = ox + oeck, oy + oeck
            ib, ih = ob - 2 * oeck, oh - 2 * oeck
            # Frage oben (oberes Viertel)
            gfx.Stiftfarbe(*farben["text"])
            gfx.SetzeFont(FONT_REGULAR, 24)
            gfx.SchreibeTextbox(ix + 4, iy + 4, ib - 8, ih * 2 // 8 - 8,
                                aktuelle_frage.frage, 24)
            # 4 Antwortkärtchen (2×2 Raster)
            a_farben = [farben["quiz_a0"], farben["quiz_a1"],
                        farben["quiz_a2"], farben["quiz_a3"]]
            a_top = iy + ih * 2 // 8
            a_rest = iy + ih - a_top
            d = 4
            pad = 8
            for i, antwort in enumerate(aktuelle_frage.antworten):
                ax = ix + (i % 2) * (ib // 2) + d
                ay = a_top + (i // 2) * (a_rest // 2) + d
                ab2 = ib // 2 - 2 * d
                ah2 = a_rest // 2 - 2 * d
                gfx.Stiftfarbe(*a_farben[i])
                gfx.Vollrechteck(ax, ay, ab2, ah2)
                gfx.Stiftfarbe(*farben["text"])
                gfx.SchreibeTextbox(ax + pad, ay + pad,
                                    ab2 - 2*pad, ah2 - 2*pad, antwort, 22)

        elif modus == "hilfe":
            gfx.Stiftfarbe(*farben["quiz"])
            gfx.Vollrechteck(ox, oy, ob, oh, oeck)
            gfx.Stiftfarbe(*farben["text"])
            gfx.SchreibeTextbox(
                ox + oeck + 8, oy + oeck + 8,
                ob - 2*oeck - 16, oh - 2*oeck - 16,
                "Hilfe\n\n"
                "Im Spielmodus (und nur, wenn alle Kugeln still stehen): "
                "Maus bewegen ändert die Zielrichtung. "
                "Stoß durch Klicken mit der linken Maustaste. "
                "Die Stoßkraft wird durch Scrollen der Maus verändert.\n\n"
                "Du spielst gegen die Zeit. Alle neun Kugeln müssen "
                "versenkt werden. Es gibt ein Foul, wenn die weiße Kugel "
                "reingeht oder wenn bei einem Stoß gar keine Kugel "
                "versenkt wird.\n\n"
                "Im Quizmodus: Klicke die richtigen Antworten an, "
                "um Fouls abzuarbeiten.\n\n"
                "Die Bedienung erfolgt durch Anklicken der Buttons "
                "oder mit der angegebenen Taste auf der Tastatur.",
                18)

        elif modus == "gameover":
            gfx.Stiftfarbe(*farben["quiz"])
            gfx.Vollrechteck(ox, oy, ob, oh, oeck)
            gfx.Stiftfarbe(*farben["text"])
            gfx.SetzeFont(FONT_BOLD, FENSTER_BREITE // 12)
            gfx.SchreibeFont(ox + ob//4, oy + oh//3, "GAME OVER")

        # FPS
        gfx.Stiftfarbe(100, 100, 100)
        gfx.SetzeFont(FONT_BOLD, 12)
        gfx.SchreibeFont(4, 2, f"{int(gfx._clock.get_fps())} fps")

        gfx.Aktualisiere(60)

    gfx.FensterAus()


if __name__ == "__main__":
    main()
