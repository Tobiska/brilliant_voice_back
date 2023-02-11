package conn

import (
	"brillian_voice_back/internal/domain/entity/room"
	"brillian_voice_back/internal/domain/entity/user"
	"strconv"
)

type StateInf struct {
	Code       string        `json:"code"`
	Users      []UserInf     `json:"users"`
	IsFully    bool          `json:"is_fully"`
	Properties PropertiesInf `json:"properties"`
	OwnerName  string        `json:"owner_name"`
	//todo status
	//todo question
	//todo time
	//todo number of round
}

func ToInfState(state room.GameState) (StateInf, error) {
	o, err := state.GetOwner()
	if err != nil {
		return StateInf{}, err
	}
	return StateInf{
		Code:    state.Descriptor.Code,
		Users:   toInfUsers(state.Users),
		IsFully: state.IsFully,
		Properties: PropertiesInf{
			CountPlayers: state.Properties.CountPlayers,
			Time:         strconv.Itoa(state.Properties.TimerDuration),
		},
		OwnerName: o.Username,
	}, nil
}

type UserInf struct {
	Name string `json:"name"`
	//todo answer
}

func toInfUsers(users map[string]*user.User) []UserInf {
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
