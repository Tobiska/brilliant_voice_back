package conn

type ErrorMsg struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}

const (
	UndefinedCode = "UNDEFINED"
)

func (pc *PlayerConn) WriteError(err error) error {
	m := &ErrorMsg{
		Msg:  err.Error(),
		Code: UndefinedCode,
	}
	return pc.ws.WriteJSON(m)
}
