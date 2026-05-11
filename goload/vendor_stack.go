package goload

import (
	"path"
	"strings"
)

type vendorLayer struct {
	p      string
	prefix string
	pkgs   map[string]string
}

func newVendorLayer(p string) *vendorLayer {
	return &vendorLayer{
		p:      p,
		prefix: path.Join(p, "vendor") + "/",
		pkgs:   make(map[string]string),
	}
}

func (ly *vendorLayer) addPkg(p string) {
	if p == ly.p {
		return
	}
	if !strings.HasPrefix(p, ly.prefix) {
		panic("not inside the vendor directory")
	}
	ly.pkgs[strings.TrimPrefix(p, ly.prefix)] = p
}

type vendorStack struct {
	layers []*vendorLayer
}

func (s *vendorStack) push(ly *vendorLayer) {
	s.layers = append(s.layers, ly)
}

func (s *vendorStack) pop() {
	n := len(s.layers)
	if n == 0 {
		panic("nothing to pop")
	}
	s.layers = s.layers[:n-1]
}

func (s *vendorStack) mapImport(p string) (string, bool) {
	n := len(s.layers)
	for i := n - 1; i >= 0; i-- {
		ly := s.layers[i]
		if mapped, found := ly.pkgs[p]; found {
			return mapped, true
		}
	}
	return p, false
}
