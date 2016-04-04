package main

import "time"

func newEventTracker() *eventTracker {
	failure := make([]time.Time, 0, 56)
	success := make([]time.Time, 0, 56)
	ev := new(eventTracker)
	ev.Failure = failure
	ev.Success = success

	return ev
}

func (e *eventTracker) addEvent(t time.Time, ev bool) {
	if ev {
		e.Success = append(e.Failure[0:], t)
	} else {
		e.Failure = append(e.Failure[0:], t)
	}
}

func (e *eventTracker) merge(t eventTracker) {
	e.Failure = append(e.Failure[0:], t.Failure[0:]...)
	e.Success = append(e.Success[0:], t.Success[0:]...)
}
