"""
tisch — Der Billardtisch mit Spielregeln und Physiksimulation.

Aufgabe für Schüler:
    Implementiere die Klasse BillardSpiel.
    Das Spiel verwaltet die Kugeln, Taschen, Punkte und den Countdown.
    Die Methode aktualisiere() treibt die Simulation einen Tick voran.
"""

import random
import time
from vec2 import Vec2
from kugel import Kugel


class BillardSpiel:
    """Ein 9-Ball-Billardspiel.

    Öffentliche Attribute:
        stoss_richtung: Vec2 — aktuelle Zielrichtung
        stoss_kraft:    float — Stärke des nächsten Stoßes (0–14)
        stillstand:     bool — stehen alle Kugeln still?
        strafpunkte:    int — Anzahl der Fouls
        restzeit:       float — verbleibende Spielzeit in Sekunden
        rk:             float — Kugelradius
    """

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
        self._laeuft = False
        self.spielzeit = 4 * 60.0
        self.restzeit = self.spielzeit
        self._letzte_zeit = time.time()

        # 6 Taschen: Ecken + Mitte
        rt, rtm = 1.9 * self.rk, 1.5 * self.rk
        b, h = self.breite, self.hoehe
        self.taschen = [Vec2(0,0), Vec2(0,h), Vec2(b/2,h),
                        Vec2(b,h), Vec2(b,0), Vec2(b/2,0)]
        self.taschen_radius = [rt, rt, rtm, rt, rt, rtm]

        self._aufstellen_9ball()

    # ---------- Aufstellung ----------

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
            b - rk - random.random() * (b/4 - 2*rk),
            random.random() * (h - 2*rk) + rk)

    # ---------- Steuerung ----------

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

    def stosse(self, sound_abspielen):
        """Stößt die Spielkugel. sound_abspielen: Funktion für den Queue-Sound."""
        if not self.stillstand or self.stoss_kraft <= 0.15:
            return
        if self.stoss_richtung.ist_null():
            return
        sound_abspielen()
        self._vorige_kugeln = [k.kopie() for k in self.kugeln]
        self._angestossen, self._angespielte = True, None
        self.spielkugel.v = self.stoss_richtung.mal(self.stoss_kraft)
        self.stoss_kraft = 5.0
        self.stillstand = False

    # ---------- Simulation ----------

    def aktualisiere(self, sound_ball, sound_bande, sound_tasche):
        """Ein Tick der Simulation.
        sound_ball, sound_bande, sound_tasche: Sound-Funktionen."""
        if not self._laeuft:
            return
        jetzt = time.time()
        self.restzeit = max(0, self.restzeit - (jetzt - self._letzte_zeit))
        self._letzte_zeit = jetzt

        aktive = [k for k in self.kugeln if not k.eingelocht]
        for i, k1 in enumerate(aktive):
            for k2 in aktive[i+1:]:
                k1.pruefe_kollision(k2, lambda: (
                    sound_ball(), self._notiere_beruehrt(k1, k2)))
            k1._kollision_mit = None
            k1.pruefe_bande(self.breite, self.hoehe, sound_bande)
            k1.bewegen()
            for tp, tr in zip(self.taschen, self.taschen_radius):
                if tp.minus(k1.pos).betrag() < tr:
                    sound_tasche()
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

    # ---------- Regeln ----------

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

    def treffer(self):
        return len(self.eingelochte)

    def aktive_kugeln(self):
        return [k for k in self.kugeln if not k.eingelocht]

    def alle_eingelocht(self):
        return all(k.eingelocht or k is self.spielkugel for k in self.kugeln)
