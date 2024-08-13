package clock

import "github.com/jonboulle/clockwork"

var clock clockwork.Clock

func InitMock(clock_ clockwork.Clock) {
	clock = clock_
}

func Get() clockwork.Clock {
	if clock == nil {
		clock = clockwork.NewRealClock()
	}
	return clock
}
