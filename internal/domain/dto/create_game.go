package dto

type InputCreateGameDto struct {
	Username          string
	ID                string
	CountPlayers      int
	TimeDurationRound int
}

type InputJoinGameDto struct {
	Username string
	ID       string
	Code     string
}
