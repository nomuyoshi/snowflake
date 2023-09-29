package clock

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func NewRealClock() realClock {
	return realClock{}
}

func (r realClock) Now() time.Time {
	return time.Now()
}

type fakeClock struct {
	t time.Time
}

func NewFakeClock(t time.Time) fakeClock {
	return fakeClock{t: t}
}

func (f fakeClock) Now() time.Time {
	return f.t
}
