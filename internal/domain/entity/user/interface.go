package user

type IConn interface {
	Send() error
	Close() error
	RequestToLeave(id string)
}
