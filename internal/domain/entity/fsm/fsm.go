package fsm

type IState interface {
	Exec(a Action) (IState, error)

	BeforeTransition() error //todo add options
}

type Fsm struct {
	cur IState
}

func (f *Fsm) Send(a Action) error {
	if err := f.cur.BeforeTransition(); err != nil {
		return nil
	}

	nextState, err := f.cur.Exec(a)
	if err != nil {
		return err
	}

	f.cur = nextState
	return nil
}
