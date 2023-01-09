// Copyright (C) 2023  Shanhu Tech Inc.
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

package httpstest

import (
	"net/http"
	"net/http/httptest"
)

// Server wraps a *httptest.Server with HTTP support.
type Server struct {
	*httptest.Server

	TLSConfigs *TLSConfigs
	Transport  *http.Transport
}

// Client creates an HTTP client which transport connects directly to the
// server.
func (s *Server) Client() *http.Client {
	return &http.Client{Transport: s.Transport}
}

// SinkTransport sinks the transport to the server
// and sets the TLS client config.
func (s *Server) SinkTransport(tr *http.Transport) {
	tr.DialContext = s.Transport.DialContext
	tr.TLSClientConfig = s.Transport.TLSClientConfig
}

// NewServer creates an HTTPS server for the given testing domains.
func NewServer(domains []string, h http.Handler) (*Server, error) {
	c, err := NewTLSConfigs(domains)
	if err != nil {
		return nil, err
	}

	server := httptest.NewUnstartedServer(h)
	server.TLS = c.Server
	server.StartTLS()

	serverHost := server.Listener.Addr().String()
	return &Server{
		Server:     server,
		TLSConfigs: c,
		Transport:  c.Sink(serverHost),
	}, nil
}

// DualServer wraps two *httptest.Server's with a transport that
// goes to one of them base on HTTP or HTTPS.
type DualServer struct {
	HTTP       *httptest.Server
	HTTPS      *httptest.Server
	TLSConfigs *TLSConfigs
	Transport  *http.Transport
}

// NewDualServer creates an HTTPS dual server for the given testing domains.
func NewDualServer(domains []string, h http.Handler) (*DualServer, error) {
	c, err := NewTLSConfigs(domains)
	if err != nil {
		return nil, err
	}

	httpsServer := httptest.NewUnstartedServer(h)
	httpsServer.TLS = c.Server
	httpsServer.StartTLS()

	httpServer := httptest.NewServer(h)

	httpAddr := httpServer.Listener.Addr().String()
	httpsAddr := httpsServer.Listener.Addr().String()
	return &DualServer{
		HTTP:       httpServer,
		HTTPS:      httpsServer,
		TLSConfigs: c,
		Transport:  c.SinkHTTPS(httpAddr, httpsAddr),
	}, nil
}

// Client creates an HTTP client which transport connects directly to one
// of the servers base on the protocol port.
func (s *DualServer) Client() *http.Client {
	return &http.Client{Transport: s.Transport}
}
