package model

type Status uint8

const (
	StatusUndefined = iota
	StatusOK
	StatusFail
)

const (
	statusLabelUndefined = "undefined"
	statusLabelOK        = "ok"
	statusLabelFail      = "fail"
)

func (s Status) String() string {
	if s == StatusOK {
		return statusLabelOK
	}

	if s == StatusFail {
		return statusLabelFail
	}

	return statusLabelUndefined
}
