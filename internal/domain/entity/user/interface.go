package user

import "io"

type IConn interface {
	SendState() error
	Close() error
	RequestToLeave(id string)

	io.Writer
}
