package counting

import (
	"sync"
	"testing"
)

func TestCounterGet(t *testing.T) {
	c := Counter{0}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				c.Add(1)
			}
		}()
	}
	wg.Wait()
	if c.Count() != 10000 {
		t.Errorf("got %d as counter value, want 10000", c.Count())
	}
}
