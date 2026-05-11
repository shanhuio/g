package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

const formatIndent = "  "

// Fprint pretty prints a JSON data blob into a writer.
func Fprint(w io.Writer, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", formatIndent)
	if err != nil {
		return err
	}
	if _, err := w.Write(bs); err != nil {
		return err
	}
	_, err = fmt.Fprintln(w)
	return err
}

// Print pretty prints a JSON data blob into stdout.
func Print(v interface{}) {
	if err := Fprint(os.Stdout, v); err != nil {
		log.Println(err)
	}
}

// Format pretty formats JSON data bytes.
func Format(bs []byte) ([]byte, error) {
	out := new(bytes.Buffer)
	if err := json.Indent(out, bs, "", formatIndent); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
