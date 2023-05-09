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

func (r *Round) CommitAnswer(u *User, text string) ResultAnswer { //todo refactor
	r.Answers[u.Username] = ResultAnswer{
		TextAnswer: text,
		Result:     r.CheckAnswer(text),
	}
	return r.Answers[u.Username]
}

func (r *Round) CheckAnswer(user string) bool {
	return r.Question.CorrectAnswer == user
}
