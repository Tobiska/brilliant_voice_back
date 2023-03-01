package game

type IConn interface {
	UpdateGame(Game)
	SendError(error)
	Close()
}
