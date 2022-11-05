// Copyright (C) 2022  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package aries

import (
	"strings"
)

type trieNode struct {
	branch string
	child  map[byte]*trieNode
	prefix string
	hit    bool // false if the node is not a leaf
}

func newTrieRoot() *trieNode {
	return newTrieNode("", "")
}

func newTrieNode(branch, prefix string) *trieNode {
	return &trieNode{
		branch: branch,
		prefix: prefix,
		child:  make(map[byte]*trieNode),
		hit:    true,
	}
}

// trieFind finds string from a trie node, return the longest matched prefix
func trieFind(root *trieNode, s string) (string, bool) {
	return root.find(s, "")
}

// add adds string s to a trie node, return false if s already exists
func (t *trieNode) add(s string) bool {
	i := 0
	m := len(s)
	if m == 0 {
		return false
	}
	key := s[0]
	if t.child[key] == nil {
		t.addChild(newTrieNode(s, t.prefix+s))
		return true
	}
	cnode := t.child[key]
	n := len(cnode.branch)
	for i < n && i < m && cnode.branch[i] == s[i] {
		i++
	}
	if i == m && i == n {
		if !cnode.hit {
			cnode.hit = true
			return true
		}
		return false
	}
	if i == m {
		newNode := newTrieNode(s, t.prefix+s)
		cnode.branch = cnode.branch[m:n]
		newNode.addChild(cnode)
		t.child[key] = newNode
		return true
	}
	if i == n {
		return cnode.add(s[n:m])
	}
	newNode := newTrieNode(s[:i], t.prefix+s[:i])
	newNode.hit = false
	cnode.branch = cnode.branch[i:]
	newNode.addChild(cnode)
	newNode.addChild(newTrieNode(s[i:], t.prefix+s))
	t.child[key] = newNode
	return true

}

func (t *trieNode) find(s, res string) (string, bool) {
	if len(s) == 0 {
		return res, t.hit
	}
	key := s[0]
	if t.child[key] == nil {
		return res, false
	}
	child := t.child[key]
	if strings.HasPrefix(s, child.branch) {
		if child.hit {
			res = child.prefix
		}
		return child.find(strings.TrimPrefix(s, child.branch), res)
	}
	return res, false
}

func (t *trieNode) addChild(c *trieNode) {
	key := c.branch[0]
	if t.child[key] != nil {
		panic("illegal trieNode append, same key")
	}
	t.child[key] = c
}
