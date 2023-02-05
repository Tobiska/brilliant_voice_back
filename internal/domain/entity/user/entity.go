package user

import "fmt"

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
	fmt.Println("connection close: ", u.ID, u.Username)
	return u.conn.Close()
}

func (u *User) DeleteAndClose() error {
	fmt.Println("user delete and connection close: ", u.ID, u.Username)
	u.conn.RequestToLeave(u.ID)
	return u.conn.Close()
}
