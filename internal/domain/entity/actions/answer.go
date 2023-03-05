package actions

import (
	"brillian_voice_back/internal/domain/entity/game"
	"fmt"
)

type Answer struct {
	Action
	Text string
}

func (a Answer) String() string {
	return fmt.Sprintf("answer text: %s game: %s", a.Text, a.U)
}

func AnswerAction(u *game.User, text string) Answer {
	return Answer{
		Text: text,
		Action: Action{
			U: u,
		},
	}
}
