package conn

import (
	"brillian_voice_back/internal/domain/entity/game"
)

type StateInf struct {
	Code          string    `json:"code"`
	Users         []UserInf `json:"users"`
	OwnerName     string    `json:"owner_name"`
	Status        string    `json:"status"`
	RestTime      string    `json:"rest_time"`
	Rounds        RoundInf  `json:"round"`
	NumberOfRound int       `json:"number_of_Round"`
}

type RoundInf struct {
	Question string `json:"question"`
	Answers  map[string]string
}

func ToInfState(state game.Game) (StateInf, error) {
	o, err := state.GetOwner()
	if err != nil {
		return StateInf{}, err
	}
	as := make(map[string]string)

	for u, a := range state.Rounds[state.NumberOfRound].Answers {
		as[u] = a.TextAnswer
	}

	return StateInf{
		NumberOfRound: state.NumberOfRound,
		Status:        state.Status,
		Code:          state.Descriptor.Code,
		Users:         toInfUsers(state.Users.ToSlice()),
		RestTime:      state.RestTime.String(),
		Rounds: RoundInf{
			Question: state.Rounds[state.NumberOfRound].Question.Text,
			Answers:  as,
		},
		OwnerName: o.Username,
	}, nil
}

type UserInf struct {
	Username string        `json:"name"`
	Answer   *ResultAnswer `json:"answer,omitempty"`
	Ready    bool          `json:"ready,omitempty"` //ready wait users
}

type ResultAnswer struct {
	TextAnswer string `json:"text_answer"`
	Result     bool   `json:"is_correct"`
}

func toInfUsers(users []*game.User) []UserInf {
	usersInf := make([]UserInf, 0, len(users))
	for _, u := range users {
		usersInf = append(usersInf, toInfUser(u))
	}
	return usersInf
}

func toInfUser(u *game.User) UserInf {
	ui := UserInf{
		Username: u.Username,
		Ready:    u.Ready,
	}

	if u.Answer != nil {
		ui.Answer = &ResultAnswer{
			TextAnswer: u.Answer.TextAnswer,
			Result:     u.Answer.Result,
		}
	}
	return ui
}

type PropertiesInf struct {
	CountPlayers int    `json:"count_players"`
	Time         string `json:"time"`
}
