package caco3

import (
	"shanhu.io/g/lexing"
)

const (
	nodeSrc  = "src"
	nodeRule = "rule"
	nodeOut  = "out"
	nodeRun  = "run"
	nodeSub  = "sub"
)

type buildNode struct {
	name string
	typ  string
	deps []string
	pos  *lexing.Pos

	ruleType string
	rule     buildRule
	ruleMeta *buildRuleMeta

	sub *subBuilds
}

func (n *buildNode) mainOut() string {
	if m := n.ruleMeta; m != nil && len(m.outs) > 0 {
		return m.outs[0]
	}
	return ""
}
