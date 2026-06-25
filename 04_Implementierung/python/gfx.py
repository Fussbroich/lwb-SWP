"""
gfx — Einfache 2D-Grafik-Bibliothek für den Informatik-Unterricht.
Basiert auf Pygame. Installation: pip install pygame

Beispiel — ein roter Kreis auf schwarzem Hintergrund:

    import gfx

    gfx.Fenster(800, 600, "Hallo Welt")

    while gfx.fenster_offen():
        gfx.Cls()                       # Bildschirm löschen
        gfx.Stiftfarbe(255, 0, 0)       # Rot
        gfx.Vollkreis(400, 300, 50)     # Kreis in der Mitte
        gfx.Aktualisiere()              # Bild anzeigen

    gfx.FensterAus()
"""

import pygame
import math
import os

# =====================================================================
#  Interner Zustand — nicht direkt benutzen, nur über die Funktionen!
# =====================================================================

_screen = None          # Pygame-Anzeigefläche
_canvas = None          # Zeichenfläche (unterstützt Transparenz)
_clock = None           # Steuert die Bildwiederholrate
_offen = False          # Ist das Fenster gerade offen?
_breite = 0             # Fensterbreite in Pixeln
_hoehe = 0              # Fensterhöhe in Pixeln

_farbe = (0, 0, 0)      # Aktuelle Stiftfarbe (R, G, B)
_alpha = 255             # 255 = deckend, 0 = unsichtbar

_font = None             # Aktueller Pygame-Font
_font_cache = {}         # Pfad+Größe → Font (vermeidet Neuladen)

_sound_cache = {}        # Pfad → Sound (vermeidet Neuladen)

# Eingabe-Zustand (wird jeden Frame in Aktualisiere() erneuert)
_tasten_neu = []         # Tasten, die in DIESEM Frame gedrückt wurden
_maus_klick = []         # Maustasten, die in DIESEM Frame geklickt wurden
_maus_los = []           # Maustasten, die in DIESEM Frame losgelassen wurden
_scroll = 0              # Scrollrad-Richtung in diesem Frame


# =====================================================================
#  Fensterverwaltung
# =====================================================================

def Fenster(breite, hoehe, titel="gfx"):
    """Öffnet ein Grafikfenster.
    breite, hoehe: Größe in Pixeln.
    titel: Text in der Titelleiste."""
    global _screen, _canvas, _clock, _offen, _breite, _hoehe

    pygame.init()
    pygame.mixer.init(frequency=44100, size=-16, channels=2, buffer=1024)
    _screen = pygame.display.set_mode((breite, hoehe))
    _canvas = pygame.Surface((breite, hoehe), pygame.SRCALPHA)
    _clock = pygame.time.Clock()
    _breite, _hoehe = breite, hoehe
    _offen = True
    pygame.display.set_caption(titel)


def fenster_offen():
    """Gibt True zurück, solange das Fenster offen ist."""
    return _offen


def FensterAus():
    """Schließt das Fenster und beendet Pygame. Darf mehrfach aufgerufen werden."""
    global _offen, _font, _screen, _canvas
    if not _offen:
        return
    _offen = False
    _font = None
    _font_cache.clear()
    _screen = None
    _canvas = None
    pygame.quit()


def Fenstertitel(titel):
    """Ändert den Text in der Titelleiste."""
    pygame.display.set_caption(titel)


# =====================================================================
#  Zeichensteuerung
# =====================================================================

def Aktualisiere(fps=60):
    """Zeigt alles Gezeichnete an und wartet bis zum nächsten Frame.
    fps: Bildwiederholrate (Standard: 60 Bilder pro Sekunde)."""
    global _tasten_neu, _maus_klick, _maus_los, _scroll, _offen

    # Eingabe-Ereignisse verarbeiten
    _tasten_neu = []
    _maus_klick = []
    _maus_los = []
    _scroll = 0

    for ereignis in pygame.event.get():
        if ereignis.type == pygame.QUIT:
            _offen = False
        elif ereignis.type == pygame.KEYDOWN:
            _tasten_neu.append(ereignis.key)
        elif ereignis.type == pygame.MOUSEBUTTONDOWN:
            if ereignis.button <= 3:
                _maus_klick.append(ereignis.button)
        elif ereignis.type == pygame.MOUSEBUTTONUP:
            if ereignis.button <= 3:
                _maus_los.append(ereignis.button)
        elif ereignis.type == pygame.MOUSEWHEEL:
            _scroll = ereignis.y    # positiv = hoch, negativ = runter

    # Zeichenfläche auf den Bildschirm kopieren
    _screen.blit(_canvas, (0, 0))
    pygame.display.flip()
    _clock.tick(fps)


def Cls(r=0, g=0, b=0):
    """Löscht den Bildschirm. Standard: Schwarz.
    Oder: Cls(100, 100, 100) für Grau."""
    _canvas.fill((r, g, b, 255))


def Stiftfarbe(r, g, b):
    """Setzt die Zeichenfarbe. Werte von 0 bis 255.
    Beispiel: Stiftfarbe(255, 0, 0) = Rot."""
    global _farbe
    _farbe = (r, g, b)


def Transparenz(t):
    """Setzt die Durchsichtigkeit. 0 = deckend, 255 = unsichtbar."""
    global _alpha
    _alpha = 255 - t


def _rgba():
    """Gibt die aktuelle Farbe als (R, G, B, A)-Tupel zurück."""
    return (*_farbe, _alpha)


# =====================================================================
#  Zeichenprimitiven
# =====================================================================

def Linie(x1, y1, x2, y2, breite=1):
    """Zeichnet eine Linie von (x1,y1) nach (x2,y2).
    breite: Strichbreite in Pixeln (Standard: 1)."""
    pygame.draw.line(_canvas, _rgba(),
                     (int(x1), int(y1)), (int(x2), int(y2)), max(1, int(breite)))


def Kreis(x, y, r):
    """Zeichnet einen Kreisumriss."""
    pygame.draw.circle(_canvas, _rgba(), (int(x), int(y)), int(r), 1)


def Vollkreis(x, y, r):
    """Zeichnet einen ausgefüllten Kreis."""
    pygame.draw.circle(_canvas, _rgba(), (int(x), int(y)), int(r))


def Rechteck(x, y, b, h, eckradius=0):
    """Zeichnet einen Rechteckumriss bei (x,y) mit Breite b und Höhe h.
    eckradius: Abrundung der Ecken in Pixeln (0 = eckig)."""
    pygame.draw.rect(_canvas, _rgba(), (int(x), int(y), int(b), int(h)), 1,
                     border_radius=int(eckradius))


def Vollrechteck(x, y, b, h, eckradius=0):
    """Zeichnet ein ausgefülltes Rechteck.
    eckradius: Abrundung der Ecken in Pixeln (0 = eckig)."""
    pygame.draw.rect(_canvas, _rgba(), (int(x), int(y), int(b), int(h)),
                     border_radius=int(eckradius))


def Volldreieck(x1, y1, x2, y2, x3, y3):
    """Zeichnet ein ausgefülltes Dreieck."""
    punkte = [(int(x1), int(y1)), (int(x2), int(y2)), (int(x3), int(y3))]
    pygame.draw.polygon(_canvas, _rgba(), punkte)


def Kreissektor(x, y, r, w1, w2):
    """Zeichnet einen Kreisbogen (Umriss).
    Winkel in Grad: 0° = rechts (Osten), gegen den Uhrzeigersinn."""
    punkte = _sektor_punkte(x, y, r, w1, w2)
    if len(punkte) >= 2:
        pygame.draw.lines(_canvas, _rgba(), True, punkte)


def Vollkreissektor(x, y, r, w1, w2):
    """Zeichnet einen ausgefüllten Kreissektor (Tortenstück).
    Winkel in Grad: 0° = rechts (Osten), gegen den Uhrzeigersinn."""
    punkte = _sektor_punkte(x, y, r, w1, w2)
    if len(punkte) >= 3:
        pygame.draw.polygon(_canvas, _rgba(), punkte)


def _sektor_punkte(cx, cy, r, w1, w2):
    """Berechnet die Eckpunkte eines Kreissektors als Polygon."""
    start, end = float(w1), float(w2)
    if end <= start:
        end += 360
    segmente = max(4, int((end - start) / 5) + 1)

    punkte = [(int(cx), int(cy))]
    for i in range(segmente + 1):
        winkel = math.radians(start + (end - start) * i / segmente)
        px = cx + r * math.cos(winkel)
        py = cy - r * math.sin(winkel)      # Y-Achse ist gespiegelt
        punkte.append((int(px), int(py)))
    return punkte


# =====================================================================
#  Textausgabe
# =====================================================================

def SetzeFont(pfad, groesse):
    """Lädt eine Schriftart (TTF-Datei) in der angegebenen Größe.
    Beispiel: SetzeFont("assets/fontfiles/LiberationMono-Bold.ttf", 16)"""
    global _font, _font_pfad
    if not pygame.font.get_init():
        pygame.font.init()
    schluessel = (pfad, groesse)
    if schluessel not in _font_cache:
        _font_cache[schluessel] = pygame.font.Font(pfad, groesse)
    _font = _font_cache[schluessel]
    _font_pfad = pfad


def SchreibeFont(x, y, text):
    """Schreibt Text an Position (x,y) mit der aktuell gesetzten Schriftart."""
    if _font is None:
        return
    bild = _font.render(text, True, _farbe)
    if _alpha < 255:
        bild.set_alpha(_alpha)
    _canvas.blit(bild, (int(x), int(y)))


def SchreibeTextbox(x, y, breite, hoehe, text, groesse=16):
    """Schreibt Text mit automatischem Zeilenumbruch in ein Rechteck.
    Wörter, die nicht mehr in eine Zeile passen, wandern in die nächste.
    Zeilen, die unter das Rechteck fallen würden, werden abgeschnitten.

    Beispiel:
        gfx.Stiftfarbe(255, 255, 255)
        gfx.SchreibeTextbox(50, 50, 300, 200,
                            "Dies ist ein langer Text, der automatisch umbrochen wird.")
    """
    SetzeFont(_font_pfad if _font_pfad else None, groesse)
    zeilenhoehe = groesse + 4
    cx, cy = x, y

    for absatz in text.split("\n"):
        if absatz == "":
            cy += zeilenhoehe
            continue
        woerter = absatz.split()
        for wort in woerter:
            wort_breite = TextBreite(wort)
            # Passt das Wort noch in die aktuelle Zeile?
            if cx > x and cx + wort_breite > x + breite:
                cx = x
                cy += zeilenhoehe
            # Unterhalb des Rechtecks? Aufhören.
            if cy + groesse > y + hoehe:
                return
            SchreibeFont(cx, cy, wort)
            cx += wort_breite + TextBreite(" ")
        # Absatzende → neue Zeile
        cx = x
        cy += zeilenhoehe


# Intern: letzter Font-Pfad für SchreibeTextbox
_font_pfad = None

def TextBreite(text):
    """Gibt die Breite des Textes in Pixeln zurück."""
    if _font is None:
        return len(text) * 8
    try:
        return _font.size(text)[0]
    except Exception:
        return len(text) * 8


# =====================================================================
#  Eingabe — Tastatur
# =====================================================================

def TasteGedrueckt(taste):
    """Wurde diese Taste in DIESEM Frame neu gedrückt?
    taste: ein Zeichen ('h', 'q', '1') oder ein pygame-Keycode.

    Beispiel: if gfx.TasteGedrueckt('h'): hilfe_anzeigen()"""
    code = ord(taste.lower()) if isinstance(taste, str) else taste
    return code in _tasten_neu


def TasteGehalten(taste):
    """Wird diese Taste gerade gehalten (dauerhaft gedrückt)?
    Nützlich für Bewegungssteuerung (z.B. Pfeiltasten)."""
    code = ord(taste.lower()) if isinstance(taste, str) else taste
    return pygame.key.get_pressed()[code]


# =====================================================================
#  Eingabe — Maus
# =====================================================================

def MausPosition():
    """Gibt die aktuelle Mausposition als (x, y) zurück."""
    return pygame.mouse.get_pos()


def MausGeklickt(taste=1):
    """Wurde die Maustaste in DIESEM Frame geklickt?
    taste: 1 = links, 2 = mitte, 3 = rechts."""
    return taste in _maus_klick


def MausLosgelassen(taste=1):
    """Wurde die Maustaste in DIESEM Frame losgelassen?"""
    return taste in _maus_los


def MausScroll():
    """Gibt die Scrollrichtung zurück: +1 = hoch, -1 = runter, 0 = kein Scroll."""
    return _scroll


# =====================================================================
#  Sound
# =====================================================================

def LadeSound(pfad):
    """Lädt eine WAV-Datei und gibt ein Sound-Objekt zurück.
    Das Objekt kann mit sound.play() abgespielt werden."""
    if pfad not in _sound_cache:
        _sound_cache[pfad] = pygame.mixer.Sound(pfad)
    return _sound_cache[pfad]


def SpieleSound(pfad):
    """Spielt eine WAV-Datei einmal ab. Lädt sie beim ersten Aufruf."""
    LadeSound(pfad).play()


def SpieleMusik(pfad, loop=True):
    """Spielt eine Musikdatei im Hintergrund (streamt, spart RAM).
    loop: True = Endlosschleife, False = einmal abspielen."""
    pygame.mixer.music.load(pfad)
    pygame.mixer.music.play(-1 if loop else 0)


def StoppeMusik():
    """Stoppt die Hintergrundmusik."""
    pygame.mixer.music.stop()
