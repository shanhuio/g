package aries

import (
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"shanhu.io/g/unixhttp"
)

// ListenAndServe serves on the address. If the address ends
// with .sock, it ListenAndServe's on the unix domain socket.
func ListenAndServe(addr string, s Service) error {
	log.Printf("serve on %q", addr)
	if strings.HasSuffix(addr, ".sock") {
		return unixhttp.ListenAndServe(addr, Serve(s))
	}
	return http.ListenAndServe(addr, Serve(s))
}

// Listen listens on the address. If the address ends with
// .sock, it Listen's on the unix domain socket.
func Listen(addr string) (net.Listener, error) {
	if strings.HasSuffix(addr, ".sock") {
		return unixhttp.Listen(addr)
	}
	return net.Listen("tcp", addr)
}

// DefaultAddr gets the default address for an application
func DefaultAddr(app string) string {
	if addr := os.Getenv("ADDR"); addr != "" {
		return addr
	}
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	h := fnv.New32()
	io.WriteString(h, app)
	const offset = 8000
	return fmt.Sprintf("localhost:%d", offset+h.Sum32()%10000)
}

// DeclareAddrFlag declares the -addr flag.
func DeclareAddrFlag(def string) *string {
	if def == "" {
		def = DefaultAddr(os.Args[0])
	}
	return flag.String("addr", def, "address to listen on")
}
