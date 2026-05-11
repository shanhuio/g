// Package settings provides a generic, simple key-value pair like interface
// for saving application settings.
package settings

// Settings is an interface for saving simple JSON object based settings.
type Settings interface {
	// Get gets a setting. Returns errcode.NotFound error when the
	// key is missing.
	Get(key string, v interface{}) error

	// Set sets a setting.
	Set(key string, v interface{}) error

	// Has checks if a key exists. It does not have to read the key.
	Has(key string) (bool, error)
}

// String gets a string-type value from the settings.
func String(b Settings, key string) (string, error) {
	var s string
	if err := b.Get(key, &s); err != nil {
		return "", err
	}
	return s, nil
}

// Bytes gets a []byte type value from the settings.
func Bytes(b Settings, key string) ([]byte, error) {
	var bs []byte
	if err := b.Get(key, &bs); err != nil {
		return nil, err
	}
	return bs, nil
}
