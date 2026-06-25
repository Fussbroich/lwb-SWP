"""
kugel — Eine Billardkugel mit Physik.

    Eine Kugel hat eine Position, eine Geschwindigkeit, einen Radius
    und eine Nummer (Wert). Sie kann sich bewegen, an Banden abprallen
    und mit anderen Kugeln kollidieren (elastischer Stoß).
"""

from vec2 import Vec2


class Kugel:
    """Eine Billardkugel.

    Attribute:
        pos:        Vec2 — Position auf dem Tisch
        v:          Vec2 — Geschwindigkeit (Pixel pro Tick)
        radius:     float — Kugelradius in Pixeln
        wert:       int — Nummer der Kugel (0 = weiße Spielkugel)
        eingelocht: bool — Wurde die Kugel versenkt?
    """

    def __init__(self, pos, radius, wert):
        self.pos = pos
        self.v = Vec2()
        self.radius = radius
        self.wert = wert
        self.eingelocht = False
        self._kollision_mit = None

    def bewegen(self):
        """Bewegt die Kugel einen Tick weiter und bremst sie durch Reibung."""
        self.pos = self.pos.plus(self.v)
        geschwindigkeit = self.v.betrag()
        if geschwindigkeit > 0.15:
            self.v = self.v.mal(1 - 0.02 / geschwindigkeit)
        else:
            self.v = Vec2()

    def pruefe_bande(self, breite, hoehe, bei_treffer):
        """Prüft, ob die Kugel eine Bande berührt, und reflektiert sie.
        bei_treffer: Funktion, die bei Bandenkontakt aufgerufen wird."""
        if self.eingelocht:
            return
        vx, vy = self.v.x, self.v.y
        x, y, r = self.pos.x, self.pos.y, self.radius
        trifft = False

        am_rand = not (r <= x <= breite - r and r <= y <= hoehe - r)
        if not am_rand and x + vx < r:       vx = -vx; trifft = True
        if not am_rand and x + vx > breite-r: vx = -vx; trifft = True
        if not am_rand and y + vy < r:        vy = -vy; trifft = True
        if not am_rand and y + vy > hoehe-r:  vy = -vy; trifft = True

        if trifft:
            bei_treffer()
            self.v = Vec2(vx, vy)

    def pruefe_kollision(self, andere, bei_treffer):
        """Elastischer Stoß mit einer anderen Kugel (gleiche Masse).
        bei_treffer: Funktion, die bei Kontakt aufgerufen wird."""
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
            bei_treffer()

        self.v, andere.v = u1, u2
        self._kollision_mit = andere
        andere._kollision_mit = self

    def kopie(self):
        """Gibt eine unabhängige Kopie dieser Kugel zurück."""
        k = Kugel(Vec2(self.pos.x, self.pos.y), self.radius, self.wert)
        k.eingelocht = self.eingelocht
        return k
