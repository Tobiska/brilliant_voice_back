package game

type Round struct {
	Question Question
	Answers  map[string]ResultAnswer
}

type Question struct {
	Text          string
	CorrectAnswer string
}

type ResultAnswer struct {
	TextAnswer string
	Result     bool
}

func (r *Round) Answer(id string, answerText string) {
	r.Answers[id] = ResultAnswer{
		TextAnswer: answerText,
		Result:     r.compareAnswers(r.Question.CorrectAnswer, answerText),
	}
}

func (r *Round) compareAnswers(correct, user string) bool {
	return correct == user
}
