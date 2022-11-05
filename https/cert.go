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

package https

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"shanhu.io/pub/errcode"
)

// Cert contains a certificate in memory.
type Cert struct {
	Cert []byte // Marshalled PEM block for the certificate.
	Key  []byte // Marshalled PEM block for the private key.
}

// X509KeyPair converts the PEM blocks into a X509 key pair
// for use in an HTTPS server.
func (c *Cert) X509KeyPair() (tls.Certificate, error) {
	return tls.X509KeyPair(c.Cert, c.Key)
}

// CertConfig is the configuration for creating a RSA-based HTTPS
// certificate.
type CertConfig struct {
	Hosts    []string
	IsCA     bool
	Start    *time.Time
	Duration time.Duration
}

// NewCACert creates a CA cert for the given domain.
func NewCACert(domain string) (*Cert, error) {
	c := &CertConfig{
		Hosts: []string{domain},
		IsCA:  true,
	}
	return MakeRSACert(c, 0)
}

func (c *CertConfig) start() time.Time {
	if c.Start != nil {
		return *c.Start
	}
	return time.Now()
}

func (c *CertConfig) duration() time.Duration {
	if c.Duration <= 0 {
		return time.Hour * 24 * 30
	}
	return c.Duration
}

func makeTemplate(c *CertConfig) (*x509.Certificate, error) {
	start := c.start()
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errcode.Annotate(err, "generate serial number")
	}

	const org = "Acme Co"
	const keyUsage = x509.KeyUsageKeyEncipherment |
		x509.KeyUsageDigitalSignature

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{Organization: []string{org}},
		NotBefore:    start,
		NotAfter:     start.Add(c.duration()),

		KeyUsage:    keyUsage,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		BasicConstraintsValid: true,
	}

	for _, h := range c.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if c.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	return &template, nil
}

// MakeRSACert creates RSA-based TLS certificates.
func MakeRSACert(c *CertConfig, bits int) (*Cert, error) {
	if len(c.Hosts) == 0 {
		return nil, errcode.InvalidArgf("no host specified")
	}

	if bits == 0 {
		bits = 2048
	}
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, errcode.Annotate(err, "generate private key")
	}
	template, err := makeTemplate(c)
	if err != nil {
		return nil, errcode.Annotate(err, "make certificate template")
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader, template, template, &priv.PublicKey, priv,
	)
	if err != nil {
		return nil, errcode.Annotate(err, "create certificate")
	}

	certOut := new(bytes.Buffer)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut := new(bytes.Buffer)
	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	return &Cert{
		Cert: certOut.Bytes(),
		Key:  keyOut.Bytes(),
	}, nil
}

// MakeEC256Cert creates a ECDSA-256 TLS certificate.
func MakeEC256Cert(c *CertConfig) (*Cert, error) {
	return MakeECCert(c, elliptic.P256())
}

// MakeECCert creates a ECDSA-based HTTPs certificate.
func MakeECCert(c *CertConfig, curve elliptic.Curve) (*Cert, error) {
	if len(c.Hosts) == 0 {
		return nil, errcode.InvalidArgf("no host specified")
	}

	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, errcode.Annotate(err, "generate private key")
	}
	template, err := makeTemplate(c)
	if err != nil {
		return nil, errcode.Annotate(err, "make certificate template")
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader, template, template, &priv.PublicKey, priv,
	)
	if err != nil {
		return nil, errcode.Annotate(err, "create certificate")
	}

	certOut := new(bytes.Buffer)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	pemBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, errcode.Annotate(err, "marshal key bytes")
	}
	keyOut := new(bytes.Buffer)
	pem.Encode(keyOut, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: pemBytes,
	})

	return &Cert{
		Cert: certOut.Bytes(),
		Key:  keyOut.Bytes(),
	}, nil
}
