package signals

type Signal byte

const (
	SignalOn Signal = iota
	SignalOff
	SignalPublisherDisconnected
	SignalPublisherConnected
	SignalUpdateSubscribersStatistic
	SignalPing
	SignalPong
)
