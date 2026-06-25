"""
vec2 — Ein 2D-Vektor für Positionen und Geschwindigkeiten.

Aufgabe für Schüler:
    Implementiere die Methoden der Klasse Vec2.
    Ein Vektor hat zwei Komponenten (x, y) und beschreibt
    eine Position oder eine Richtung im 2D-Raum.
"""

import math


class Vec2:
    """Ein 2D-Vektor mit den Komponenten x und y.

    Beispiel:
        pos = Vec2(100, 200)
        richtung = Vec2(1, 0)         # zeigt nach rechts
        neue_pos = pos.plus(richtung.mal(5))
    """

    def __init__(self, x=0.0, y=0.0):
        self.x = float(x)
        self.y = float(y)

    def plus(self, other):
        """Addiert zwei Vektoren und gibt einen neuen zurück."""
        return Vec2(self.x + other.x, self.y + other.y)

    def minus(self, other):
        """Subtrahiert einen Vektor und gibt einen neuen zurück."""
        return Vec2(self.x - other.x, self.y - other.y)

    def mal(self, skalar):
        """Multipliziert den Vektor mit einer Zahl."""
        return Vec2(self.x * skalar, self.y * skalar)

    def betrag(self):
        """Die Länge des Vektors (Satz des Pythagoras)."""
        return math.sqrt(self.x ** 2 + self.y ** 2)

    def normiert(self):
        """Gibt einen Vektor gleicher Richtung mit Länge 1 zurück."""
        b = self.betrag()
        return Vec2(self.x / b, self.y / b) if b > 0 else Vec2()

    def ist_null(self):
        """Prüft ob der Vektor (0, 0) ist."""
        return self.x == 0 and self.y == 0

    def skalarprodukt(self, other):
        """Das Skalarprodukt (dot product) zweier Vektoren."""
        return self.x * other.x + self.y * other.y

    def projiziert_auf(self, other):
        """Projiziert diesen Vektor auf einen anderen (Vektorprojektion)."""
        b2 = other.skalarprodukt(other)
        if b2 == 0:
            return Vec2()
        return other.mal(self.skalarprodukt(other) / b2)

    def __repr__(self):
        return f"Vec2({self.x:.1f}, {self.y:.1f})"
