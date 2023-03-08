package conn

import (
	"brillian_voice_back/internal/domain/entity/game"
	"strconv"
)

type StateInf struct {
	Code            string        `json:"code"`
	Users           []UserInf     `json:"users"`
	IsFully         bool          `json:"is_fully"`
	Properties      PropertiesInf `json:"properties"`
	OwnerName       string        `json:"owner_name"`
	Status          string        `json:"status"`
	CurrentQuestion string        `json:"current_question"`
	//todo time
	NumberOfRound int `json:"number_of_Round"`
}

func ToInfState(state game.Game) (StateInf, error) {
	o, err := state.GetOwner()
	if err != nil {
		return StateInf{}, err
	}
	return StateInf{
		NumberOfRound: state.NumberOfRound,
		Status:        state.Status,
		Code:          state.Descriptor.Code,
		Users:         toInfUsers(state.Users.ToSlice()),
		IsFully:       state.IsFully,
		Properties: PropertiesInf{
			CountPlayers: state.Properties.CountPlayers,
			Time:         strconv.Itoa(state.Properties.TimerDuration),
		},
		OwnerName: o.Username,
	}, nil
}

type UserInf struct {
	Name   string `json:"name"`
	Answer string `json:"answer"`
}

func toInfUsers(users []*game.User) []UserInf {
	usersInf := make([]UserInf, 0, len(users))
	for _, u := range users {
		usersInf = append(usersInf, UserInf{
			Name: u.Username,
		})
	}
	return usersInf
}

type PropertiesInf struct {
	CountPlayers int    `json:"count_players"`
	Time         string `json:"time"`
}
