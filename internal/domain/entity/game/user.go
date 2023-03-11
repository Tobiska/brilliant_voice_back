package game

import (
	"fmt"
)

type User struct {
	ID       string
	Username string

	Answer *ResultAnswer

	Ready bool

	Conn IConn
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s, username: %s", u.ID, u.Username)
}

func NewUser(ID, username string, conn IConn) *User {
	return &User{
		ID:       ID,
		Username: username,
		Conn:     conn,
	}
}

func (u *User) Reset() {
	u.Answer = nil
	u.Ready = false
}
