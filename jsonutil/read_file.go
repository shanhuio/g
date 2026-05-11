package jsonutil

import (
	"encoding/json"
	"os"
)

// ReadFile reads and unmarshals a JSON file.
func ReadFile(file string, obj interface{}) error {
	bs, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, obj)
}
