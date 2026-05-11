package subcmd

import (
	"strings"
)

func parseCmd(arg string) (string, string) {
	cmd, host, found := strings.Cut(arg, "@")
	if found {
		return cmd, host
	}
	return arg, ""
}
