package testutil

import (
	"time"
)

type Clock struct {
	now time.Time
}

func NewClock(now time.Time) *Clock {
	return &Clock{now}
}

func (p *Clock) Now() time.Time {
	return p.now
}
