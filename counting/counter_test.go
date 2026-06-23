package counting

import (
	"sync"
	"testing"
)

func TestCounterGet(t *testing.T) {
	c := Counter{0}
	var wg sync.WaitGroup

	for range 10 {
		wg.Go(func() {
			for range 1000 {
				c.Add(1)
			}
		})
	}
	wg.Wait()
	if c.Count() != 10000 {
		t.Errorf("got %d as counter value, want 10000", c.Count())
	}
}
