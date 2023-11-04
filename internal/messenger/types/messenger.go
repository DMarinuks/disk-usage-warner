package types

type Messenger interface {
	Send(hostname string, warnings []*WarningInfo) error
}

type WarningInfo struct {
	Device  string
	Percent string
}
