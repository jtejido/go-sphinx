package util

import (
	"sync"
)

var (
	timerPoolInstance *TimerPool
)

type TimerPool struct {
	sync.Mutex
	weakRefTimerPool map[any][]*Timer
}

func initTimer() {
	timerPoolInstance = &TimerPool{
		weakRefTimerPool: make(map[any][]*Timer),
	}

}

// GetLogMath returns the singleton instance of LogMath
func GetTimerPool() *TimerPool {
	once.Do(initTimer)
	return timerPoolInstance
}

func (tp *TimerPool) GetTimer(owner any, timerName string) *Timer {
	tp.Lock()
	defer tp.Unlock()

	if _, ok := tp.weakRefTimerPool[owner]; !ok {
		tp.weakRefTimerPool[owner] = make([]*Timer, 0)
	}

	ownerTimers := tp.weakRefTimerPool[owner]

	for _, timer := range ownerTimers {
		if timer.Name() == timerName {
			return timer
		}
	}

	// there is no timer named 'timerName' yet, so create it
	requestedTimer := NewTimer(timerName)
	ownerTimers = append(ownerTimers, NewTimer(timerName))

	return requestedTimer
}
