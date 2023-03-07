package inmemory

import (
	"brillian_voice_back/internal/domain/entity/game"
	"context"
)

var questions = []game.Question{
	{
		Text:          "Домодедово это дом дедов?",
		CorrectAnswer: "Скорее дед домов",
	},
	{
		Text:          "Если пёс то какой?",
		CorrectAnswer: "Корги королевский",
	},
	{
		Text:          "В попестамеско это фамилия?",
		CorrectAnswer: "Да",
	},
	{
		Text:          "Дочь папиного друга...",
		CorrectAnswer: "...",
	},
	{
		Text:          "Папа, мама, брат, ...",
		CorrectAnswer: "сестра",
	},
}

type Provider struct{}

func NewRoundProvider() *Provider {
	return &Provider{}
}

func (p *Provider) PrepareRounds(_ context.Context) ([]*game.Round, error) {
	rounds := make([]*game.Round, 5)
	roundsVal := make([]game.Round, 5)
	for i := 0; i < 5; i++ {
		roundsVal[i].Question = questions[i]
		roundsVal[i].Answers = make(map[string]string)
		rounds[i] = &roundsVal[i]
	}
	return rounds, nil
}
