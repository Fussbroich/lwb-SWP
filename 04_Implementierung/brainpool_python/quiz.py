"""
quiz — Quizfragen aus CSV-Dateien laden und verwalten.

Aufgabe für Schüler:
    Implementiere QuizFrage und die Funktion lade_quiz().
    Eine CSV-Datei hat das Format:
    Frage;Antwort1;Antwort2;Antwort3;Antwort4;RichtigerIndex
"""

import csv
import random


class QuizFrage:
    """Eine Quizfrage mit vier Antwortmöglichkeiten.

    Attribute:
        frage:      str — der Fragetext
        antworten:  list[str] — vier Antwortmöglichkeiten
    """

    def __init__(self, frage, antworten, richtige_index):
        self.frage = frage
        self.antworten = antworten
        self.richtige_index = richtige_index

    def ist_richtig(self, index):
        """Prüft, ob die gewählte Antwort (0–3) korrekt ist."""
        return index == self.richtige_index


def lade_quiz(csv_pfad):
    """Lädt Quizfragen aus einer CSV-Datei und mischt sie.

    Gibt eine Liste von QuizFrage-Objekten zurück.
    Bei Fehler wird eine leere Liste zurückgegeben.
    """
    fragen = []
    try:
        with open(csv_pfad, encoding="utf-8") as f:
            for zeile in csv.reader(f, delimiter=";"):
                if len(zeile) == 6:
                    fragen.append(QuizFrage(
                        zeile[0], zeile[1:5], int(zeile[5])))
    except FileNotFoundError:
        return []
    random.shuffle(fragen)
    return fragen
