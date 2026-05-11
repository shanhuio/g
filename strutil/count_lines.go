package strutil

import (
	"bufio"
	"bytes"
)

// CountLines counts the number of lines in the given byte slice.
func CountLines(bs []byte) int {
	s := bufio.NewScanner(bytes.NewReader(bs))
	n := 0
	for s.Scan() {
		n++
	}
	return n
}
