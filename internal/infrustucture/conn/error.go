package conn

type ErrorMsg struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}

const (
	UndefinedCode = "UNDEFINED"
)

func (pc *PlayerConn) WriteError(err error) {
	m := &ErrorMsg{
		Msg:  err.Error(),
		Code: UndefinedCode,
	}
	_ = pc.ws.WriteJSON(m)
}
