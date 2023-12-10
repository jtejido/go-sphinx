package util

import (
	"fmt"
	"math"
	"time"
)

var timeFormatter = "%.4f"

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type Timer struct {
	name                                             string
	sum, count, startTime, curTime, minTime, maxTime int64
	notReliable                                      bool
}

func NewTimer(name string) *Timer {
	tm := new(Timer)
	tm.name = name
	tm.Reset()
	return tm
}

func (tm Timer) Name() string {
	return tm.name
}

func (tm *Timer) Reset() {
	tm.startTime = 0
	tm.count = 0
	tm.sum = 0
	tm.minTime = math.MaxInt64
	tm.maxTime = 0
	tm.notReliable = false
}

func (tm Timer) IsStarted() bool {
	return (tm.startTime > 0)
}

func (tm *Timer) Start() {
	if tm.startTime != 0 {
		tm.notReliable = true // start called while timer already running
		fmt.Printf("%s timer.start() called without a stop()\n", tm.Name())
	}
	tm.startTime = makeTimestamp()
}

func (tm *Timer) StartFrom(time int64) {
	if tm.startTime != 0 {
		tm.notReliable = true // start called while timer already running
		fmt.Printf("%s timer.start() called without a stop()\n", tm.Name())
	}
	if time > makeTimestamp() {
		panic("Start time is later than current time")
	}

	tm.startTime = time
}

func (tm *Timer) Stop() int64 {
	if tm.startTime == 0 {
		tm.notReliable = true // stop called, but start never called
		panic("Timer.stop() called without a start()")
	}
	tm.curTime = makeTimestamp() - tm.startTime
	tm.startTime = 0
	if tm.curTime > tm.maxTime {
		tm.maxTime = tm.curTime
	}
	if tm.curTime < tm.minTime {
		tm.minTime = tm.curTime
	}
	tm.count++
	tm.sum += tm.curTime
	return tm.curTime
}

func (tm Timer) Count() int64 {
	return tm.count
}

func (tm Timer) CurTime() int64 {
	return tm.curTime
}

func (tm Timer) AverageTime() int64 {
	if tm.count == 0 {
		return 0.0
	}
	return tm.sum / tm.count
}

func (tm Timer) MinTime() int64 {
	return tm.minTime
}

func (tm Timer) MaxTime() int64 {
	return tm.maxTime
}

func (tm Timer) FmtTime(time int64) string {
	return PadWithMinLength(fmt.Sprintf(timeFormatter+"s", time/1000.0), 10)
}
