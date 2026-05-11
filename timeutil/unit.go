package timeutil

import (
	"time"
)

// Approximate time units.
const (
	Day    = 24 * time.Hour
	Week   = 7 * Day
	Mounth = 30 * Day
	Year   = 364 * Day
)
