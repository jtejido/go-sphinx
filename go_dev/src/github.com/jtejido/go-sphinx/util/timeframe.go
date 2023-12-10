package util

import "fmt"

type TimeFrame struct {
	start, end int64
}

var (
	NULL     = NewTimeFrameFromDuration(0)
	INFINITE = NewTimeFrame(0, 1<<63-1)
)

func NewTimeFrameFromDuration(duration int64) *TimeFrame {
	return NewTimeFrame(0, duration)
}

func NewTimeFrame(start, end int64) *TimeFrame {
	return &TimeFrame{start, end}
}

func (tf *TimeFrame) GetStart() int64 {
	return tf.start
}

func (tf *TimeFrame) GetEnd() int64 {
	return tf.end
}

func (tf *TimeFrame) Length() int64 {
	return tf.end - tf.start
}

func (tf *TimeFrame) String() string {
	return fmt.Sprintf("%d:%d", tf.start, tf.end)
}
