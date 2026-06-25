"""
farben — Farbschemata für das Spiel (Hell und Dunkel).

Jede Farbe ist ein (R, G, B)-Tupel mit Werten von 0 bis 255.
Das aktive Schema kann zur Laufzeit gewechselt werden.
"""

# Farben der Billardkugeln (Index = Kugelwert)
KUGEL_FARBEN = [
    (252, 253, 242),   # 0: weiß
    (255, 201, 78),    # 1: gelb
    (34, 88, 175),     # 2: blau
    (249, 73, 68),     # 3: rot
    (84, 73, 149),     # 4: violett
    (255, 139, 33),    # 5: orange
    (47, 159, 52),     # 6: grün
    (155, 53, 30),     # 7: dunkelrot
    (48, 49, 54),      # 8: schwarz
    (255, 201, 78),    # 9: gelb (Streifen)
]

# Helles Farbschema (Standard)
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

# Dunkles Farbschema
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
