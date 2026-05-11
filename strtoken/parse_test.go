package strtoken

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	o := func(line string, args ...string) {
		strs, errs := Parse(line)
		if errs != nil {
			t.Errorf("Parse(%q): unexpected errors", line)
			for _, err := range errs {
				t.Log(err)
			}
			return
		}

		if len(strs) != len(args) {
			t.Errorf("Parse(%q): expect %d args, got %d",
				line, len(args), len(strs),
			)
			return
		}

		for i, s := range strs {
			if s != args[i] {
				t.Errorf("Parse(%q), arg %d: expect %q, got %q",
					line, i, args[i], s,
				)
			}
		}
	}
	o("")
	o("a", "a")
	o(`"a"`, "a")
	o(`/something`, "/something")
	o(`ls /x_file`, "ls", "/x_file")
	o("       ls \t\t a", "ls", "a")
	o(`"a-b" something`, "a-b", "something")
	o(`a-b`, "a-b")
	o(`?`, "?")
	o(`!x`, "!x")
	o(`mkdir -p /root/.ssh`, "mkdir", "-p", "/root/.ssh")

	e := func(line string) {
		_, errs := Parse(line)
		if errs == nil {
			t.Errorf("Parse(%q): expect error but passed", line)
		}
	}

	e(`"`)
	e(`"asdf" asdf "xx`)
}
