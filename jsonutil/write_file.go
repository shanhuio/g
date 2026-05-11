package jsonutil

import (
	"bytes"
	"encoding/json"
	"os"
)

// WriteFile marshals a JSON object and writes it into a file.
func WriteFile(file string, obj interface{}) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return os.WriteFile(file, bs, 0644)
}

// WriteFileReadable marshals a JSON object with indents and writes it into a
// file.
func WriteFileReadable(f string, v interface{}) error {
	buf := new(bytes.Buffer)
	bs, err := json.MarshalIndent(v, "", formatIndent)
	if err != nil {
		return err
	}
	buf.Write(bs)
	buf.Write([]byte("\n"))

	return os.WriteFile(f, buf.Bytes(), 0644)
}
