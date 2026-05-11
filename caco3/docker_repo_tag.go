package caco3

import (
	"fmt"

	"shanhu.io/g/dock"
)

func parseRepoTag(repoTag string) (string, string) {
	repo, tag := dock.ParseImageTag(repoTag)
	if tag == "" {
		tag = "latest"
	}
	return repo, tag
}

func repoTag(repo, tag string) string {
	if tag == "" {
		return repo
	}
	return fmt.Sprintf("%s:%s", repo, tag)
}
