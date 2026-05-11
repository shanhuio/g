package aries

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"shanhu.io/g/strutil"
)

type staticFileSystem struct {
	fs http.FileSystem
}

func (s *staticFileSystem) Open(name string) (http.File, error) {
	lastDot := strings.LastIndex(filepath.Base(name), ".")
	if lastDot >= 0 {
		return s.fs.Open(name)
	}

	f, err := s.fs.Open(name)
	if errors.Is(err, fs.ErrNotExist) {
		html := name + ".html"
		f2, err2 := s.fs.Open(html)
		if err2 == nil {
			return f2, nil
		}
		if !errors.Is(err2, fs.ErrNotExist) {
			log.Printf("try to open %q: %s", html, err2)
		}
	}
	return f, err
}

// StaticFiles is a module that serves static files.
type StaticFiles struct {
	cacheControl string
	h            http.Handler
}

// DefaultStaticPath is the default path for static files.
const DefaultStaticPath = "lib/site"

func cacheControl(ageSecs int) string {
	return fmt.Sprintf("max-age=%d; must-revalidate", ageSecs)
}

// NewStaticFiles creates a module that serves static files.
func NewStaticFiles(p string) *StaticFiles {
	p = strutil.Default(p, DefaultStaticPath)
	fs := &staticFileSystem{fs: http.Dir(p)}
	return &StaticFiles{
		cacheControl: cacheControl(10),
		h:            http.FileServer(fs),
	}
}

// CacheAge sets the maximum age for the cache.
func (s *StaticFiles) CacheAge(ageSecs int) {
	if ageSecs < 0 {
		s.cacheControl = ""
	} else {
		s.cacheControl = cacheControl(ageSecs)
	}
}

var contentTypeSuffix = []struct {
	suffix      string
	contentType string
}{
	{suffix: ".js", contentType: "application/javascript;charset=UTF-8"},
	{suffix: ".css", contentType: "text/css;charset=UTF-8"},
}

// Serve serves incoming HTTP requests.
func (s *StaticFiles) Serve(c *C) error {
	c.Req.URL.Path = c.Path
	if s.cacheControl != "" {
		c.Resp.Header().Add("Cache-Control", s.cacheControl)
	}
	for _, suf := range contentTypeSuffix {
		if strings.HasSuffix(c.Req.URL.Path, suf.suffix) {
			c.Resp.Header().Set("Content-Type", suf.contentType)
			break
		}
	}
	s.h.ServeHTTP(c.Resp, c.Req)
	return nil
}
