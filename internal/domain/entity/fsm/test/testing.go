package test

import (
	"brillian_voice_back/internal/domain/entity/game"
	"brillian_voice_back/internal/domain/entity/properties"
)

var (
	Game = &game.Game{
		Users:   game.NewUsersContainer(),
		OwnerId: "admin_code",
		Descriptor: game.Descriptor{
			Code:    "code",
			IsFully: false,
			Properties: properties.Properties{
				CountPlayers:  2,
				TimerDuration: 5,
			},
		},
	}
)
