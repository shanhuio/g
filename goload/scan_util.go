package goload

import (
	"go/build"
	"sort"
)

func isNoGoError(e error) bool {
	if e == nil {
		return false
	}
	_, hit := e.(*build.NoGoError)
	return hit
}

func inSet(s map[string]bool, k string) bool {
	if s == nil {
		return false
	}
	return s[k]
}

func findInSorted(strs []string, target string) bool {
	index := sort.SearchStrings(strs, target)
	return index < len(strs) && strs[index] == target
}
