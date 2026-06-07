// Package dagvis provides clear visualization of a DAG.
package dagvis

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"shanhu.io/g/aries"
	"shanhu.io/g/dags"
	"shanhu.io/g/osutil"
	"shanhu.io/std/errcode"
)

type server struct {
	dag    *dags.M
	static *aries.StaticFiles
	tmpls  *aries.Templates
}

func (s *server) serveIndex(c *aries.C) error {
	pageData := struct {
		Graph *dags.M
	}{
		Graph: s.dag,
	}

	dat := struct {
		PageData any
	}{
		PageData: &pageData,
	}

	return s.tmpls.Serve(c, "dagview.html", &dat)
}

func makeService(home string) (aries.Service, error) {
	h, err := osutil.NewHome(home)
	if err != nil {
		return nil, errcode.Annotate(err, "make new home")
	}

	m := new(dags.M)
	dagBytes, err := os.ReadFile(h.Var("dagview.json"))
	if err != nil {
		return nil, errcode.Annotate(err, "read dagview.json")
	}
	if err := json.Unmarshal(dagBytes, m); err != nil {
		return nil, errcode.Annotate(err, "parse dagview.json")
	}

	s := &server{
		dag:    m,
		static: aries.NewStaticFiles(h.Lib("static")),
		tmpls:  aries.NewTemplates(h.Lib("tmpl"), nil),
	}

	serveStatic := s.static.Serve

	r := aries.NewRouter()
	r.Index(s.serveIndex)
	r.Get("style.css", serveStatic)
	r.Dir("js", serveStatic)
	r.Dir("jslib", serveStatic)

	return r, nil
}

// Main is main.
func Main() {
	addr := aries.DeclareAddrFlag("")
	home := flag.String("home", ".", "home dir")
	flag.Parse()

	s, err := makeService(*home)
	if err != nil {
		log.Fatal(err)
	}
	if err := aries.ListenAndServe(*addr, s); err != nil {
		log.Fatal(err)
	}
}
