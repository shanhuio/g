package caco3

import (
	"log"
	"strings"
)

type gitRemote struct {
	name  string
	git   string
	fetch bool
	push  bool
}

func listRemotes(dir string) (map[string]*gitRemote, error) {
	output, err := runCmdOutput(dir, "git", "remote", "-v")
	if err != nil {
		return nil, err
	}

	remotes := make(map[string]*gitRemote)

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 3 {
			name := fields[0]
			git := fields[1]
			method := fields[2]
			remote, ok := remotes[name]
			if ok {
				if git != remote.git {
					log.Printf(
						"inconsistent remote url for %q: %q",
						name, line,
					)
					continue
				}
			} else {
				remote = &gitRemote{
					name: name,
					git:  git,
				}
				remotes[name] = remote
			}
			if method == "(fetch)" {
				remote.fetch = true
			} else if method == "(push)" {
				remote.push = true
			} else {
				log.Printf("unknown git remote method: %q", line)
			}
		} else {
			log.Printf("weird git remote line: %q", line)
		}
	}

	return remotes, nil
}
