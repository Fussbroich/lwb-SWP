package welt

func NewBeispielQuiz() *quiz {
	return &quiz{
		fragen: []QuizFrage{
			&quizfrage{frage: "Wie lautet der erste Buchstabe des Alphabets?",
				antworten: [4]string{
					"Der erste Buchstabe des Alphabets lautet A.",
					"Die Frage ist sinnlos.",
					"Was ist ein Alphabet?",
					"D"},
				richtig: 0},
			&quizfrage{frage: "Wie lautet der zweite Buchstabe des Alphabets?",
				antworten: [4]string{
					"Der Buchstabe, der als zweites kommt, lautet A, nicht wahr?",
					"B",
					"Das ist C.",
					"D"},
				richtig: 1},
			&quizfrage{frage: "Wie lautet der dritte Buchstabe des Alphabets?",
				antworten: [4]string{
					"A",
					"Der Buchstabe, der als zweites kommt, lautet B oder C, nicht wahr?",
					"C",
					"D"},
				richtig: 2},
			&quizfrage{frage: "Wie lautet der vierte Buchstabe des Alphabets?",
				antworten: [4]string{
					"A",
					"B",
					"Äh vier. Der Buchstabe, der als viertes kommt, lautet C, nicht wahr?",
					"D"},
				richtig: 3},
			&quizfrage{frage: "Wie lautet der fünfte Buchstabe des Alphabets?",
				antworten: [4]string{
					"D",
					"E",
					"F",
					"G"},
				richtig: 0},
			&quizfrage{frage: "Wie lautet der sechste Buchstabe des Alphabets?",
				antworten: [4]string{
					"D",
					"E",
					"F",
					"G"},
				richtig: 1},
			&quizfrage{frage: "Wie lautet der siebte Buchstabe des Alphabets?",
				antworten: [4]string{
					"D",
					"E",
					"F",
					"G"},
				richtig: 2}}}
}
