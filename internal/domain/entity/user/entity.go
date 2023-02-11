package user

import (
	"github.com/rs/zerolog/log"
)

type User struct {
	ID       string
	Username string

	conn IConn
}

func NewUser(ID, username string, conn IConn) *User {
	return &User{
		ID:       ID,
		Username: username,
		conn:     conn,
	}
}

func (u *User) Update() error {
	return u.conn.Send()
}

func (u *User) Close() error {
	log.Info().Str("ID", u.ID).Msg("connection close")
	return u.conn.Close()
}

func (u *User) DeleteAndClose() error {
	log.Info().Str("ID", u.ID).Msg("user delete and connection close")
	u.conn.RequestToLeave(u.ID)
	return u.conn.Close()
}
