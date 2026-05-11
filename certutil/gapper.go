package certutil

import (
	"sync"
	"time"
)

// gapper provides a check function that returns true where the first true is
// after the first specified time point, and every two other adjacent true has
// a gap of no less than period.
type gapper struct {
	mu     sync.Mutex
	next   time.Time
	period time.Duration
}

func newGapper(period time.Duration, first time.Time) *gapper {
	return &gapper{
		next:   first,
		period: period,
	}
}

func newGapperNow(period time.Duration, now time.Time) *gapper {
	return newGapper(period, now.Add(period))
}

// check returns true if now is after the first specified time point if this is
// the first true, or if now is not before period after the last time check()
// returns true.
func (t *gapper) check(now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !now.Before(t.next) {
		t.next = now.Add(t.period)
		return true
	}
	return false
}
