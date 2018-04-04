// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Modifications Copyright 2014 Skycoin authors.

package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/utc"
)

var logger = logging.MustGetLogger("util")

// GenerateCert generates a self-signed X.509 certificate for a TLS server. Outputs to
// certFile and keyFile and will overwrite existing files.
func GenerateCert(certFile, keyFile, host, organization string, rsaBits int,
	isCA bool, validFrom time.Time, validFor time.Duration) error {
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return fmt.Errorf("Failed to generate private key: %v", err)
	}

	notBefore := validFrom
	notAfter := notBefore.Add(validFor)

	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	if notAfter.After(endOfTime) {
		notAfter = endOfTime
	}

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{organization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("Failed to create certificate: %v", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("Failed to open %s for writing: %v", certFile, err)
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open %s for writing:%v", keyFile, err)
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return nil
}

func certKeyXor(certFile, keyFile string) (bool, error) {
	certInfo, err := os.Stat(certFile)
	certExists := !os.IsNotExist(err)
	certIsFile := certExists && certInfo.Mode().IsRegular()

	keyInfo, err := os.Stat(keyFile)
	keyExists := !os.IsNotExist(err)
	keyIsFile := keyExists && keyInfo.Mode().IsRegular()

	switch {
	case certExists && certIsFile && keyExists && keyIsFile:
		return true, nil
	case !certExists && !keyExists:
		return false, nil
	case !certExists:
		return false, fmt.Errorf("Cert %s does not exist", certFile)
	case !certIsFile:
		return false, fmt.Errorf("Cert %s is not a file", certFile)
	case !keyExists:
		return false, fmt.Errorf("Key %s does not exist", keyFile)
	case !keyIsFile:
		return false, fmt.Errorf("Key %s is not a file", keyFile)
	default:
		panic("unreachable code")
	}
}

// CreateCertIfNotExists verifies that certFile and keyFile exist and are files, and
// if neither certFile nor keyFile exist, they are automatically created for host
func CreateCertIfNotExists(host, certFile, keyFile string, appName string) error {
	// check that cert/key both exist, or dont
	exist, err := certKeyXor(certFile, keyFile)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	// Automatically create a new cert if neither files exist
	logger.Infof("Creating certificate %s", certFile)
	logger.Infof("Creating key %s", keyFile)
	lifetime := time.Hour * 365 * 24 // 1 year
	if err := GenerateCert(certFile, keyFile, host, appName, 2048, false, utc.Now(), lifetime); err != nil {
		return err
	}

	logger.Infof("Created certificate %s for host %s", certFile, host)
	logger.Infof("Created key %s for host %s", keyFile, host)
	return nil
}
