package identity

import (
	"testing"

	"time"
)

func TestMemCore(t *testing.T) {
	now := time.Now()
	core := NewMemCore(func() time.Time { return now })

	coreConfig := SingleKeyCoreConfig(now.Add(time.Hour))
	initID, err := core.Init(coreConfig)
	if err != nil {
		t.Fatal("init core: ", err)
	}

	if len(initID.PublicKeys) != 1 {
		t.Fatalf("got %d keys, want one", len(initID.PublicKeys))
	}

	k := initID.PublicKeys[0]
	t.Logf("key id: %s", k.ID)
}
