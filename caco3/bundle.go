package caco3

import (
	"shanhu.io/g/errcode"
)

type bundle struct {
	rule *Bundle
	name string
	deps []string
}

func newBundle(env *env, p string, r *Bundle) *bundle {
	name := makeRelPath(p, r.Name)
	var deps []string
	for _, dep := range r.Deps {
		deps = append(deps, makePath(p, dep))
	}

	return &bundle{
		name: name,
		deps: deps,
		rule: r,
	}
}

func (b *bundle) build(env *env, opts *buildOpts) error {
	return nil
}

func (b *bundle) meta(env *env) (*buildRuleMeta, error) {
	d, err := makeDigest(ruleBundle, b.name, struct{}{})
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	return &buildRuleMeta{
		name:   b.name,
		deps:   b.deps,
		digest: d,
	}, nil
}
