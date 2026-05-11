package termutil

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// ReadPassword reads a password from stdin.
func ReadPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	line, err := term.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		return nil, err
	}
	return []byte(line), nil
}
