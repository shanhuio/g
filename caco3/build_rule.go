package caco3

type buildRuleMeta struct {
	name string
	deps []string
	outs []string

	dockerOut bool // Output is a docker container image.

	// digest captures all non-dependency input such as action type, binded
	// input, external input, etc.  returns empty string if this always needs
	// re-execution.
	digest string
}

type buildRule interface {
	// meta returns meta information of a build rule.
	meta(env *env) (*buildRuleMeta, error)

	// build executes the build action.
	build(env *env, opts *buildOpts) error
}
