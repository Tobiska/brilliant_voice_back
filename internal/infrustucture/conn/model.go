package conn

import (
	"brillian_voice_back/internal/domain/entity/game"
	"time"
)

type StateInf struct {
	Code          string        `json:"code"`
	Users         []UserInf     `json:"users"`
	OwnerName     string        `json:"owner_name"`
	Status        string        `json:"status"`
	RestTime      string        `json:"rest_time"`
	Rounds        RoundInf      `json:"round"`
	Properties    PropertiesInf `json:"properties"`
	NumberOfRound int           `json:"number_of_round"`
}

type RoundInf struct {
	Question string                  `json:"question"`
	Answers  map[string]ResultAnswer `json:"answers"`
}

func ToInfState(state game.Game) (StateInf, error) {
	o, err := state.GetOwner()
	if err != nil {
		return StateInf{}, err
	}
	as := make(map[string]ResultAnswer)

	for u, a := range state.Rounds[state.NumberOfRound].Answers {
		as[u] = ResultAnswer{
			TextAnswer: a.TextAnswer,
			Result:     a.Result,
		}
	}

	return StateInf{
		NumberOfRound: state.NumberOfRound,
		Status:        state.Status,
		Code:          state.Descriptor.Code,
		Users:         toInfUsers(state.Users.ToSlice()),
		RestTime:      state.RestTime.Truncate(time.Second).String(),
		Rounds: RoundInf{
			Question: state.Rounds[state.NumberOfRound].Question.Text,
			Answers:  as,
		},
		Properties: PropertiesInf{
			CountPlayers: state.Properties.CountPlayers,
			Time:         state.Properties.TimerDuration.String(),
		},
		OwnerName: o.Username,
	}, nil
}

type UserInf struct {
	Username string `json:"name"`
	Ready    bool   `json:"ready,omitempty"`
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
	return ui
}

type PropertiesInf struct {
	CountPlayers int    `json:"count_players"`
	Time         string `json:"time"`
}
