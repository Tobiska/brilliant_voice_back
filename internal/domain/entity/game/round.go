package game

type Round struct {
	Question Question
	Answers  map[string]string
}

type Question struct {
	Text          string
	CorrectAnswer string
}
