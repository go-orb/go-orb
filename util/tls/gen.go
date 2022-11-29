package tls

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
	"errors"
	"math/big"
	"net"
	"time"
)

// LoadTLSConfig loads a TLS config from certificate files.
func LoadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	var err error

	certs := make([]tls.Certificate, 1)

	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config := tls.Config{
		Certificates: certs,
		MinVersion:   tls.VersionTLS13,
	}

	return &config, nil
}

// GenTLSConfig creates a self signed certificate and returns it in a TSL config.
func GenTLSConfig(addr ...string) (*tls.Config, error) {
	hosts := make([]string, len(addr))

	for _, addr := range addr {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	// Generate a certificate
	cert, certPool, err := Certificate(hosts...)
	// cert, err := CertificateQuic()
	if err != nil {
		return nil, err
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		NextProtos:   []string{"h1", "http/1.1"},
	}

	return &config, nil
}

// CertificateQuic generates a certificate for use with QUIC.
// NOTE: this is redundant, was just here to test, might remove later, still
// needs to check if its useful.
func CertificateQuic() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return tls.X509KeyPair(certPEM, keyPEM)
}

// Certificate generates a self signed certificate.
func Certificate(host ...string) (tls.Certificate, *x509.CertPool, error) {
	emptyCert := tls.Certificate{}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return emptyCert, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return emptyCert, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		Version:   3,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range host {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return emptyCert, nil, err
	}

	// Create public key
	certOut := bytes.NewBuffer(nil)
	if err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return emptyCert, nil, err
	}

	// Create private key
	keyOut := bytes.NewBuffer(nil)

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return emptyCert, nil, err
	}

	if err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
		return emptyCert, nil, err
	}

	cert, err := tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
	if err != nil {
		return emptyCert, nil, err
	}

	// This is needed for gRPC tests with self signed certs.
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certOut.Bytes()); !ok {
		return emptyCert, nil, errors.New("failed to append cert to root CA pool")
	}

	return cert, certPool, nil
}
