package dock

import (
	"fmt"
	"log"

	"shanhu.io/g/errcode"
)

// RunTask runs a command line.
func RunTask(c *Cont, line string) error {
	log.Println("#", line)

	exit, err := c.Exec(line)
	if err != nil {
		return errcode.Annotatef(err, "exec %q", line)
	}
	if exit != 0 {
		return fmt.Errorf("exit value: %d", exit)
	}
	return nil
}

// RunTasks runs a series of command lines. All commands must succeed and
// return 0 exit value.
func RunTasks(c *Cont, lines []string) error {
	for _, line := range lines {
		if err := RunTask(c, line); err != nil {
			return err
		}
	}
	return nil
}
