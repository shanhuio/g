package caco3

type subBuilds struct {
	dirs []string
	rule *SubBuilds
}

func newSubBuilds(env *env, p string, v *SubBuilds) *subBuilds {
	var dirs []string
	for _, d := range v.Dirs {
		dirs = append(dirs, makeRelPath(p, d))
	}

	return &subBuilds{dirs: dirs, rule: v}
}

func (b *subBuilds) Dirs() []string { return b.dirs }
