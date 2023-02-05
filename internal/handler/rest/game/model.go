package game

import validation "github.com/go-ozzo/ozzo-validation"

type CreateInput struct {
	Username        string `json:"username" form:"username" binding:"required"`
	ID              string `json:"id" form:"id" binding:"required"`
	CountPlayers    int    `json:"count_players" form:"count_players" binding:"required"`
	TimeDurationMin int    `json:"time_duration" form:"time_duration" binding:"required"`
}

func (in *CreateInput) Validate() error {
	return validation.ValidateStruct(in,
		validation.Field(in.CountPlayers, validation.Required, validation.Length(0, 6)),
		validation.Field(in.TimeDurationMin, validation.Required, validation.Length(1, 3)),
	)
}

type JoinInput struct {
	Code     string `json:"code" form:"code" schema:"code" binding:"required"`
	Username string `json:"username" form:"username" schema:"username" binding:"required"`
	ID       string `json:"id" form:"id" schema:"id" binding:"required"`
}
