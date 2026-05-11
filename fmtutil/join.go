package fmtutil

import (
	"bytes"
	"fmt"
	"reflect"
)

// Join joins a slice of stuff into a string with sep as the
// separator.
func Join(slice interface{}, sep string) string {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return fmt.Sprint(slice)
	}

	v := reflect.ValueOf(slice)
	n := v.Len()
	if n == 0 {
		return ""
	} else if n == 1 {
		return fmt.Sprint(v.Index(0).Interface())
	}

	buf := new(bytes.Buffer)
	for i := 0; i < n; i++ {
		x := v.Index(i).Interface()
		if i > 0 {
			fmt.Fprint(buf, sep)
		}
		fmt.Fprint(buf, x)
	}
	return buf.String()
}
