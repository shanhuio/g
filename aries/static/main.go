package static

import (
	"flag"
	"log"

	"shanhu.io/g/aries"
)

func makeService(root string) aries.Service {
	return aries.NewStaticFiles(root)
}

// Main is the main entrance for smlstatic binary
func Main() {
	root := flag.String("root", "lib/site", "static directory to serve")
	addr := aries.DeclareAddrFlag("localhost:8000")
	flag.Parse()

	s := makeService(*root)
	if err := aries.ListenAndServe(*addr, s); err != nil {
		log.Fatal(err)
	}
}
