package welt

func NewBeispielQuiz() *quiz {
	quizfragen := []QuizFrage{}
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der erste Buchstabe des Alphabets?",
		"Der erste Buchstabe des Alphabets lautet A.",
		"Die Frage ist sinnlos.",
		"Was ist ein Alphabet?",
		"D", 0))
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der zweite Buchstabe des Alphabets?",
		"Der Buchstabe, der als zweites kommt, lautet A, nicht wahr?",
		"B",
		"Das ist C.",
		"D", 1))
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der dritte Buchstabe des Alphabets?",
		"A",
		"Der Buchstabe, der als zweites kommt, lautet B oder C, nicht wahr?",
		"C",
		"D", 2))
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der vierte Buchstabe des Alphabets?",
		"A",
		"B",
		"Äh vier. Der Buchstabe, der als viertes kommt, lautet C, nicht wahr?",
		"D", 3))
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der fünfte Buchstabe des Alphabets?",
		"D",
		"E",
		"F",
		"G", 1))
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der sechste Buchstabe des Alphabets?",
		"Meistens D.",
		"Manchmal E.",
		"F",
		"G", 2))
	quizfragen = append(quizfragen, NewQuizFrage(
		"Wie lautet der siebte Buchstabe des Alphabets?",
		"Der lautet D.",
		"Das ist so ganz eindeutig nicht entscheidbar.",
		"Das F.",
		"G", 3))

	return &quiz{fragen: quizfragen}
}
