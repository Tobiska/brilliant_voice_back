package game

type TimerAdapter struct {
	Start chan TimerInfo
}

type TimerInfo struct {
	TimeOutPeriod int
	TickerPeriod  int
}
