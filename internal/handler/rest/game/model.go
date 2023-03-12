package game

import validation "github.com/go-ozzo/ozzo-validation"

type CreateInput struct {
	Username        string `json:"username" form:"username" binding:"required"`
	ID              string `json:"id" form:"id" binding:"required"`
	CountPlayers    int    `json:"count_players" form:"count_players" binding:"required"`
	TimeDurationSec int    `json:"time_duration" form:"time_duration" binding:"required"`
}

func (in CreateInput) Validate() error {
	return validation.ValidateStruct(&in,
		validation.Field(&in.CountPlayers, validation.Required, validation.Min(1), validation.Max(6)),
		validation.Field(&in.TimeDurationSec, validation.Required, validation.Min(10), validation.Max(360)),
	)
}

type JoinInput struct {
	Code     string `params:"code" binding:"required"`
	Username string `query:"username" binding:"required"`
	ID       string `query:"id" binding:"required"`
}
