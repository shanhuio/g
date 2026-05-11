package states

import (
	"net/url"
)

// States is a generic interface for saving states.
type States interface {
	Get(ctx C, key string) ([]byte, error)
	Put(ctx C, key string, data []byte) error
	Del(ctx C, key string) error
	URL() *url.URL
}
