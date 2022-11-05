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
	"net/http"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/trie"
)

// Router is a path router. Similar to mux, but routing base on
// a filesystem-like syntax.
type Router struct {
	index Service
	miss  Service

	trie  *trie.Trie
	nodes map[string]*routerNode
}

type routerNode struct {
	s      Service
	isDir  bool
	method string
}

// NewRouter creates a new router for filesystem like path routing.
func NewRouter() *Router {
	return &Router{
		trie:  trie.New(),
		nodes: make(map[string]*routerNode),
	}
}

// Index sets the handler function for handling the index page when hitting
// this router, that is when hitting the root of it. One can only hit this
// route node when the path is ended with a slash '/'.
func (r *Router) Index(f Func) { r.index = f }

// Default sets a default handler for handling routes that does
// not hit anything in the routing tree.
func (r *Router) Default(f Func) { r.miss = f }

// MethodFile adds a routing file node into the routing tree that accepts
// only the given method.
func (r *Router) MethodFile(m, p string, f Func) error {
	return r.add(p, &routerNode{s: f, method: m})
}

// File adds a routing file node into the routing tree.
func (r *Router) File(p string, f Func) error {
	return r.MethodFile("", p, f)
}

// Get adds a routing file node into the routing tree that handles GET
// requests.
func (r *Router) Get(p string, f Func) error {
	return r.MethodFile(http.MethodGet, p, f)
}

// Post adds a routing file node into the routing tree that handles POST
// requests.
func (r *Router) Post(p string, f Func) error {
	return r.MethodFile(http.MethodPost, p, f)
}

// JSONCall adds a JSON marshalled POST based RPC call node into the routing
// tree. The function must be in the form of
// `func(c *aries.C, req *RequestType) (resp *ResponseType, error)`,
// where RequestType
// and ResponseType are both JSON marshallable.
func (r *Router) JSONCall(p string, f interface{}) error {
	return r.Post(p, JSONCall(f))
}

// JSONCallMust is the same as JSONCall, but panics if there is an error.
func (r *Router) JSONCallMust(p string, f interface{}) {
	if err := r.JSONCall(p, f); err != nil {
		panic(err)
	}
}

// Call is an alias of JSONCallMust
func (r *Router) Call(p string, f interface{}) { r.JSONCallMust(p, f) }

// Dir adds a routing directory node into the routing tree.
func (r *Router) Dir(p string, f Func) error { return r.DirService(p, f) }

// DirService adds a service into the router tree under a directory node.
func (r *Router) DirService(p string, s Service) error {
	return r.add(p, &routerNode{s: s, isDir: true})
}

func (r *Router) add(p string, n *routerNode) error {
	if n.s == nil {
		panic("function is nil")
	}

	route := newRoute(p)
	if route.p == "" {
		panic("trying to add empty route, use Index() instead")
	}
	if r.nodes[route.p] != nil {
		return errcode.InvalidArgf("path %s already assigned", route.p)
	}

	r.nodes[route.p] = n
	ok := r.trie.Add(route.routes, route.p)
	if !ok {
		panic("adding to trie failed")
	}

	return nil
}

func (r *Router) notFound(c *C) error {
	if r.miss == nil {
		return Miss
	}
	return r.miss.Serve(c)
}

// Serve serves the incoming context. It returns Miss if the path hits
// nothing and Default() is not set.
func (r *Router) Serve(c *C) error {
	rel := c.Rel()
	if rel == "" {
		if r.index == nil {
			return r.notFound(c)
		}
		return r.index.Serve(c)
	}

	route := c.RelRoute()
	hitRoute, p := r.trie.Find(route)
	if p == "" {
		return r.notFound(c)
	}
	n := r.nodes[p]
	if n == nil {
		panic(errcode.InvalidArgf("route function not found for %q", p))
	}

	c.ShiftRoute(len(hitRoute))
	if n.isDir || (c.Rel() == "" && !c.PathIsDir()) {
		m := c.Req.Method
		if n.method != "" && m != n.method {
			return errcode.InvalidArgf("unsupported method: %q", m)
		}
		return n.s.Serve(c)
	}
	return r.notFound(c)
}
