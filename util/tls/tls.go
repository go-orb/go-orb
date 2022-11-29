/*
 * Copyright 2017 Farsight Security, Inc.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Source: https://github.com/farsightsec/go-config/blob/master/tls.go
 */

// Package tls provides TLS utilities.
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// ClientAuth provides a convenience wrapper for tls.ClientAuthType and
// conversion to and from string format.
//
// Supported string values are:
//
//	"none":           tls.NoClientCert  (default)
//	"request":        tls.RequestClientCert
//	"require":        tls.RequireAnyClientCert
//	"verify":         tls.VerifyClientCertIfGiven
//	"require+verify": tls.RequireAndVerifyClientCert
type ClientAuth struct{ tls.ClientAuthType }

var clientAuthTypes = map[string]tls.ClientAuthType{ //nolint:gochecknoglobals
	"none":           tls.NoClientCert,
	"request":        tls.RequestClientCert,
	"require":        tls.RequireAnyClientCert,
	"verify":         tls.VerifyClientCertIfGiven,
	"require+verify": tls.RequireAndVerifyClientCert,
}

// String satisfies the flag.Value interface.
func (auth *ClientAuth) String() string {
	for name, typ := range clientAuthTypes {
		if typ == auth.ClientAuthType {
			return name
		}
	}

	return ""
}

type invalidClientAuthTypeError string

func (i invalidClientAuthTypeError) Error() string {
	return fmt.Sprintf(`Invalid ClientAuthType "%s".`, string(i))
}

type invalidClientAuthTypeValueError tls.ClientAuthType

func (i invalidClientAuthTypeValueError) Error() string {
	return fmt.Sprintf("Invalid ClientAuthType value %d", i)
}

// Set satisfies the flag.Value interface.
func (auth *ClientAuth) Set(s string) error {
	if a, ok := clientAuthTypes[s]; ok {
		auth.ClientAuthType = a
		return nil
	}

	return invalidClientAuthTypeError(s)
}

// MarshalJSON satisfies the json.Marshaler interface.
func (auth ClientAuth) MarshalJSON() ([]byte, error) {
	s := auth.String()
	if s == "" {
		return nil, invalidClientAuthTypeValueError(auth.ClientAuthType)
	}

	return json.Marshal(s)
}

// MarshalYAML satisfies the yaml.Marshaler interface.
func (auth ClientAuth) MarshalYAML() (interface{}, error) {
	s := auth.String()
	if s == "" {
		return nil, invalidClientAuthTypeValueError(auth.ClientAuthType)
	}

	return s, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (auth *ClientAuth) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	return auth.Set(strings.ToLower(s))
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface.
func (auth *ClientAuth) UnmarshalYAML(u func(interface{}) error) error {
	var s string
	if err := u(&s); err != nil {
		return err
	}

	return auth.Set(strings.ToLower(s))
}

// ConfigFiles contains the configuration for TLS as it appears on the JSON
// or YAML config. Values parsed from the config are translated and loaded
// into corresponding fields in tls.Config.
type ConfigFiles struct {
	RootCAFiles   []string   `json:"rootCAFiles,omitempty" yaml:"rootCAFiles,omitempty"`     //nolint:tagliatelle
	ClientCAFiles []string   `json:"clientCAFiles,omitempty" yaml:"clientCAFiles,omitempty"` //nolint:tagliatelle
	ClientAuth    ClientAuth `json:"clientAuth,omitempty" yaml:"clientAuth,omitempty"`
	Certificates  []struct {
		CertFile string `json:"certFile" yaml:"certFile"`
		KeyFile  string `json:"keyFile" yaml:"keyFile"`
	} `json:"certificates,omitempty" yaml:"certificates,omitempty"`
}

// Config provides JSON and YAML Marshalers and Unmarshalers for loading
// values into tls.Config.
//
// The JSON and YAML configuration format is provided by the embedded
// type TLSConfig.
type Config struct {
	ConfigFiles
	*tls.Config
}

// MarshalJSON satisfies the json.Marshaler interface.
func (t Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.ConfigFiles)
}

// MarshalYAML satisfies the yaml.Marshaler interface.
func (t Config) MarshalYAML() (interface{}, error) {
	return t.ConfigFiles, nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (t *Config) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &t.ConfigFiles); err != nil {
		return
	}

	t.Config, err = loadTLSConfig(t.ConfigFiles)

	return
}

// UnmarshalYAML satisfies the yaml.Unmarshaler interface.
func (t *Config) UnmarshalYAML(u func(interface{}) error) (err error) {
	if err = u(&t.ConfigFiles); err != nil {
		return
	}

	t.Config, err = loadTLSConfig(t.ConfigFiles)

	return
}

func loadTLSConfig(jc ConfigFiles) (*tls.Config, error) {
	var err error

	tlsConfig := new(tls.Config)

	tlsConfig.ClientAuth = jc.ClientAuth.ClientAuthType

	if len(jc.RootCAFiles) > 0 {
		tlsConfig.RootCAs, err = loadCertPool(jc.RootCAFiles)
		if err != nil {
			return nil, err
		}
	}

	if len(jc.ClientCAFiles) > 0 {
		tlsConfig.ClientCAs, err = loadCertPool(jc.ClientCAFiles)
		if err != nil {
			return nil, err
		}
	}

	for _, kp := range jc.Certificates {
		cert, err := tls.LoadX509KeyPair(kp.CertFile, kp.KeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	return tlsConfig, nil
}

func loadCertPool(files []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	for _, f := range files {
		if pem, err := os.ReadFile(path.Clean(f)); err == nil {
			pool.AppendCertsFromPEM(pem)
		} else {
			return nil, err
		}
	}

	return pool, nil
}
