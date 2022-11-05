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

package httpstest

import (
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"shanhu.io/pub/https"
)

// TLSConfigs creates the certificate setup required for a set of domains.
type TLSConfigs struct {
	Domains []string
	Server  *tls.Config
	Client  *tls.Config
}

// NewTLSConfigs creates a new setup with proper TLS config and HTTP
func NewTLSConfigs(domains []string) (*TLSConfigs, error) {
	hosts := []string{"127.0.0.1", "::1"}
	hosts = append(hosts, domains...)
	c := &https.CertConfig{
		Hosts: hosts,
		IsCA:  true,
	}
	cert, err := https.MakeECCert(c, elliptic.P256())
	if err != nil {
		return nil, fmt.Errorf("make RSA cert: %s", err)
	}

	tlsCert, err := cert.X509KeyPair()
	if err != nil {
		return nil, fmt.Errorf("unmarshal TLS cert: %s", err)
	}

	serverConfig := &tls.Config{
		NextProtos:   []string{"http/1.1", "h2"},
		Certificates: []tls.Certificate{tlsCert},
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("parse x509 cert error: %s", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(x509Cert)

	return &TLSConfigs{
		Domains: domains,
		Server:  serverConfig,
		Client:  &tls.Config{RootCAs: certPool},
	}, nil
}

// Sink returns a transport where every outgoing connection dials the same
// sinkAddr but assumes the address is certified as the domains in the
// TLSConfigs.
func (c *TLSConfigs) Sink(sinkAddr string) *http.Transport {
	return &http.Transport{
		DialContext:     sink(sinkAddr),
		TLSClientConfig: c.Client,
	}
}

// SinkHTTPS returns a transport where outgoing https connections dial httpsAddr
// and all other outgoing connections dial httpAddr. When it is HTTPS, it
// assumes the address is certified as the domains in TLSConfigs.
func (c *TLSConfigs) SinkHTTPS(httpAddr, httpsAddr string) *http.Transport {
	return &http.Transport{
		DialContext:     sinkHTTPS(httpAddr, httpsAddr),
		TLSClientConfig: c.Client,
	}
}

// InsecureSink returns a transport that always dials to the specified address,
// and skips certificate verification.
func InsecureSink(sinkAddr string) *http.Transport {
	return &http.Transport{
		DialContext:     sink(sinkAddr),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
