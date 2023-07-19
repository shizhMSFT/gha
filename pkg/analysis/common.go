package analysis

import "time"

type TimeFrame struct {
	Start time.Time
	End   time.Time
}

func (tf *TimeFrame) Union(other TimeFrame) {
	if !tf.Start.IsZero() {
		if other.Start.IsZero() || tf.Start.After(other.Start) {
			tf.Start = other.Start
		}
	}
	if !tf.End.IsZero() {
		if other.End.IsZero() || tf.End.Before(other.End) {
			tf.End = other.End
		}
	}
}

func (tf *TimeFrame) Contains(t time.Time) bool {
	return (tf.Start.IsZero() || !t.Before(tf.Start)) && (tf.End.IsZero() || !t.After(tf.End))
}
