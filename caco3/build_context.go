package caco3

type buildContext struct {
	nodes map[string]*buildNode

	built map[string]string // mapping to digests

	cache *buildCache
}

func (c *buildContext) nodeType(n string) string {
	node, ok := c.nodes[n]
	if !ok {
		return ""
	}
	return node.typ
}

func (c *buildContext) ruleType(n string) string {
	node, ok := c.nodes[n]
	if !ok {
		return ""
	}
	return node.ruleType
}
