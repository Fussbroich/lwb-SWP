"""
BrainPool — Das Mini-Billard für Schlaue.

Dies ist das Hauptprogramm. Es verbindet alle Module:
    gfx       — Grafik-Bibliothek
    tisch     — Spiellogik und Physik
    quiz      — Quizfragen
    zeichnen  — Darstellung
    farben    — Farbschemata

Steuerung:
    Maus bewegen / klicken / scrollen — Zielen, Stoßen, Kraft
    H = Hilfe, N = Neu, M = Musik, D = Dunkel/Hell, S = Schließen
"""

import os
import gfx
from vec2 import Vec2
from tisch import BillardSpiel
from quiz import lade_quiz
from farben import HELL, DUNKEL
from zeichnen import zeichne_spielfeld, zeichne_anzeige, Button

# =====================================================================
#  Pfade
# =====================================================================

ASSET = os.path.join(os.path.dirname(__file__), "..", "assets")
FONT_BOLD = os.path.join(ASSET, "fontfiles", "LiberationMono-Bold.ttf")
FONT_REGULAR = os.path.join(ASSET, "fontfiles", "LiberationMono-Regular.ttf")
SOUND_CUE = os.path.join(ASSET, "soundfiles", "cueHitsBall.wav")
SOUND_BALL = os.path.join(ASSET, "soundfiles", "ballHitsBall.wav")
SOUND_RAIL = os.path.join(ASSET, "soundfiles", "ballHitsRail.wav")
SOUND_POCKET = os.path.join(ASSET, "soundfiles", "ballIntoPocket.wav")
SOUND_AMBIENCE = os.path.join(ASSET, "soundfiles", "billardPubAmbience.wav")
SOUND_MUSIK = os.path.join(ASSET, "soundfiles", "coolJazzLoop2641.wav")
QUIZ_CSV = os.path.join(ASSET, "quizfragen", "InformatiksystemQuiz.csv")


# =====================================================================
#  Hauptprogramm
# =====================================================================

def main():
    # ---------- Layout ----------
    BREITE = 1024
    g = BREITE // 32
    HOEHE = 22 * g
    SF_B, SF_H = 3 * BREITE // 4, 3 * BREITE // 8
    SX, SY = 4 * g, 6 * g
    BANDE = int(g + g / 3)
    OECK = BANDE - 2

    # ---------- Initialisierung ----------
    gfx.Fenster(BREITE, HOEHE, "BrainPool — Das Mini-Billard")
    spiel = BillardSpiel(SF_B, SF_H)
    spiel.starte()

    for s in [SOUND_CUE, SOUND_BALL, SOUND_RAIL, SOUND_POCKET]:
        gfx.LadeSound(s)
    if os.path.exists(SOUND_AMBIENCE):
        gfx.SpieleMusik(SOUND_AMBIENCE)

    fragen = lade_quiz(QUIZ_CSV)
    frage_idx = 0
    schema = HELL
    darkmode = False
    modus = "spiel"     # "spiel", "quiz", "hilfe", "gameover"

    # ---------- Buttons ----------
    eck = g // 3
    btn_y = SY + SF_H + 5 * g // 2
    btn_h = g * 3 // 4
    n_btn = 5
    total = SF_B + 2 * BANDE - 2 * g
    btn_b = total // n_btn - g // 2
    abst = total // n_btn
    x0 = SX - BANDE + g

    def hilfe():
        nonlocal modus
        if modus == "hilfe": modus = "spiel"; spiel.starte()
        else:                modus = "hilfe"; spiel.stoppe()

    def neu():
        nonlocal modus
        spiel.reset(); spiel.starte(); modus = "spiel"

    def musik():
        if os.path.exists(SOUND_MUSIK): gfx.SpieleMusik(SOUND_MUSIK)

    def dark():
        nonlocal schema, darkmode
        darkmode = not darkmode
        schema = DUNKEL if darkmode else HELL

    def ende():
        nonlocal modus
        gfx.StoppeMusik(); modus = "ende"

    buttons = [
        Button("(H)ilfe",     x0,            btn_y, btn_b, btn_h, hilfe, eck),
        Button("(N)eu",       x0 + abst,     btn_y, btn_b, btn_h, neu, eck),
        Button("(M)usik",     x0 + 2 * abst, btn_y, btn_b, btn_h, musik, eck),
        Button("(D)unkel",    x0 + 3 * abst, btn_y, btn_b, btn_h, dark, eck),
        Button("(S)chließen", x0 + 4 * abst, btn_y, btn_b, btn_h, ende, eck),
    ]

    # ---------- Game-Loop ----------
    while gfx.fenster_offen():
        mx, my = gfx.MausPosition()

        # ====== Eingabe ======
        if gfx.TasteGedrueckt('s') or gfx.TasteGedrueckt('q'): ende()
        if gfx.TasteGedrueckt('h'): hilfe()
        if gfx.TasteGedrueckt('n'): neu()
        if gfx.TasteGedrueckt('m'): musik()
        if gfx.TasteGedrueckt('d'): dark()
        if modus == "ende": break

        if gfx.MausLosgelassen():
            for btn in buttons:
                if btn.geklickt(mx, my):
                    btn.aktion(); break
        if modus == "ende": break

        # Spielsteuerung — nur im Spielfeld
        if modus == "spiel" and spiel.laeuft() and spiel.stillstand:
            if not spiel.spielkugel.eingelocht:
                kx, ky = SX + spiel.spielkugel.pos.x, SY + spiel.spielkugel.pos.y
                r = Vec2(mx - kx, my - ky)
                if r.betrag() > 0:
                    spiel.stoss_richtung = r.normiert()
            im_spielfeld = SX <= mx <= SX + SF_B and SY <= my <= SY + SF_H
            if gfx.MausGeklickt() and im_spielfeld:
                spiel.stosse(lambda: gfx.SpieleSound(SOUND_CUE))
            scroll = gfx.MausScroll()
            if scroll:
                spiel.stoss_kraft = max(0, min(14, spiel.stoss_kraft + scroll))

        # Quiz-Klick
        if modus == "quiz" and fragen and gfx.MausLosgelassen():
            ix, iy = SX - BANDE + 2 + OECK, SY - BANDE + 2 + OECK
            ib = SF_B + 2*BANDE - 4 - 2*OECK
            ih = SF_H + 2*BANDE - 4 - 2*OECK
            a_top = iy + ih * 2 // 8
            a_rest = iy + ih - a_top
            d = 4
            for i in range(4):
                ax = ix + (i%2)*(ib//2) + d
                ay = a_top + (i//2)*(a_rest//2) + d
                if ax <= mx <= ax + ib//2 - 2*d and ay <= my <= ay + a_rest//2 - 2*d:
                    if fragen[frage_idx].ist_richtig(i):
                        spiel.strafpunkte = max(0, spiel.strafpunkte - 1)
                    frage_idx = (frage_idx + 1) % len(fragen)
                    break

        # ====== Logik ======
        if modus == "spiel":
            spiel.aktualisiere(
                lambda: gfx.SpieleSound(SOUND_BALL),
                lambda: gfx.SpieleSound(SOUND_RAIL),
                lambda: gfx.SpieleSound(SOUND_POCKET))
            if spiel.strafpunkte > spiel.treffer() + 2 and fragen:
                spiel.stoppe(); modus = "quiz"
            elif spiel.restzeit <= 0:
                spiel.stoppe(); modus = "gameover"

        if modus == "quiz":
            if spiel.strafpunkte == 0 or spiel.strafpunkte < spiel.treffer():
                modus = "spiel"; spiel.starte()

        # ====== Zeichnen ======
        gfx.Cls(*schema["hintergrund"])

        gfx.Stiftfarbe(*schema["bande"])
        gfx.Vollrechteck(SX-BANDE, SY-BANDE, SF_B+2*BANDE, SF_H+2*BANDE, BANDE)

        zeichne_spielfeld(spiel, SX, SY, SF_B, SF_H, schema, FONT_BOLD, FONT_REGULAR)
        zeichne_anzeige(spiel, SX, SY, SF_B, g, schema, FONT_BOLD)
        for btn in buttons:
            btn.zeichne(schema, FONT_BOLD)

        # Overlays
        ox, oy = SX - BANDE + 2, SY - BANDE + 2
        ob, oh = SF_B + 2*BANDE - 4, SF_H + 2*BANDE - 4

        if modus == "quiz" and fragen:
            gfx.Stiftfarbe(*schema["quiz"])
            gfx.Vollrechteck(ox, oy, ob, oh, OECK)
            ix, iy = ox + OECK, oy + OECK
            ib, ih = ob - 2*OECK, oh - 2*OECK
            gfx.Stiftfarbe(*schema["text"])
            gfx.SchreibeTextbox(ix+4, iy+4, ib-8, ih*2//8-8,
                                fragen[frage_idx].frage, 24)
            a_top, a_rest = iy + ih*2//8, ih - ih*2//8
            a_farben = [schema["quiz_a0"], schema["quiz_a1"],
                        schema["quiz_a2"], schema["quiz_a3"]]
            d, pad = 4, 8
            for i, ant in enumerate(fragen[frage_idx].antworten):
                ax = ix + (i%2)*(ib//2) + d
                ay = a_top + (i//2)*(a_rest//2) + d
                ab2, ah2 = ib//2 - 2*d, a_rest//2 - 2*d
                gfx.Stiftfarbe(*a_farben[i])
                gfx.Vollrechteck(ax, ay, ab2, ah2)
                gfx.Stiftfarbe(*schema["text"])
                gfx.SchreibeTextbox(ax+pad, ay+pad, ab2-2*pad, ah2-2*pad, ant, 22)

        elif modus == "hilfe":
            gfx.Stiftfarbe(*schema["quiz"])
            gfx.Vollrechteck(ox, oy, ob, oh, OECK)
            gfx.Stiftfarbe(*schema["text"])
            gfx.SchreibeTextbox(
                ox+OECK+8, oy+OECK+8, ob-2*OECK-16, oh-2*OECK-16,
                "Hilfe\n\n"
                "Im Spielmodus (und nur, wenn alle Kugeln still stehen): "
                "Maus bewegen ändert die Zielrichtung. "
                "Stoß durch Klicken mit der linken Maustaste. "
                "Die Stoßkraft wird durch Scrollen verändert.\n\n"
                "Du spielst gegen die Zeit. Alle neun Kugeln müssen "
                "versenkt werden. Foul = weiße Kugel rein oder "
                "nichts versenkt. Zu viele Fouls → Quiz lösen.\n\n"
                "H=Hilfe N=Neu M=Musik D=Dunkel S=Schließen", 18)

        elif modus == "gameover":
            gfx.Stiftfarbe(*schema["quiz"])
            gfx.Vollrechteck(ox, oy, ob, oh, OECK)
            gfx.Stiftfarbe(*schema["text"])
            gfx.SetzeFont(FONT_BOLD, BREITE // 12)
            gfx.SchreibeFont(ox + ob//4, oy + oh//3, "GAME OVER")

        gfx.Stiftfarbe(100, 100, 100)
        gfx.SetzeFont(FONT_BOLD, 12)
        gfx.SchreibeFont(4, 2, f"{int(gfx._clock.get_fps())} fps")

        gfx.Aktualisiere(60)

    gfx.FensterAus()


if __name__ == "__main__":
    main()
