package fsm

import "fmt"

type Action string

const (
	StartType  Action = "start"
	CloseType  Action = "close"
	PingType   Action = "ping"
	ReadyType  Action = "ready"
	AnswerType Action = "answer"
)

type IAction interface {
	fmt.Stringer
	Type() Action
}

type MetaInfo struct {
	UserId string `json:"user_id"`
}

type Start struct {
	MetaInfo
}

func (s *Start) String() string {
	return fmt.Sprintf("action: %s, user_id: %s", s.Type(), s.UserId)
}

func (s *Start) Type() Action {
	return StartType
}

type Close struct {
	MetaInfo
}

func (c *Close) String() string {
	return fmt.Sprintf("action: %s, user_id: %s", c.Type(), c.UserId)
}

func (c *Close) Type() Action {
	return CloseType
}

type Ping struct {
	MetaInfo
}

func (p *Ping) String() string {
	return fmt.Sprintf("action: %s, user_id: %s", p.Type(), p.UserId)
}

func (p *Ping) Type() Action {
	return PingType
}

type Ready struct {
	MetaInfo
}

func (r *Ready) String() string {
	return fmt.Sprintf("action: %s, user_id: %s", r.Type(), r.UserId)
}

func (r *Ready) Type() Action {
	return ReadyType
}

type Answer struct {
	MetaInfo
	Text string `json:"text"`
}

func (a *Answer) String() string {
	return fmt.Sprintf(`action: %s, user_id: %s, text: %s`, a.Type(), a.UserId, a.Text)
}

func (a *Answer) Type() Action {
	return AnswerType
}
