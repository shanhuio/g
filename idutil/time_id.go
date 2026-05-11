package idutil

import (
	"time"

	"shanhu.io/g/rand"
)

// TimeRandID generates a timestamp based unique ID.
func TimeRandID(t time.Time, nbytes int) string {
	return t.Format("20060102-150405") + "-" + rand.HexBytes(nbytes)
}
