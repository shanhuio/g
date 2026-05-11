package dock

import (
	"sort"
)

func unmapEnv(m map[string]string) []string {
	var ret []string
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ret = append(ret, k+"="+m[k])
	}
	return ret
}
