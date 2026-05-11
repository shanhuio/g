package pisces

import (
	"testing"
)

func TestMemKV(t *testing.T) {
	for _, test := range kvTestSuite {
		t.Log(test.name)
		var kv *KV
		if !test.ordered {
			kv = NewMemKV()
		} else {
			kv = NewOrderedMemKV()
		}
		test.f(t, kv)
	}
}
