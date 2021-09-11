package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

func NewServerCertificate(certfile string, keyfile string) (*tls.Config, error) {
	// generate a ec private key
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// write private key to pem file
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if nil != err {
		return nil, err
	}
	privPemBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}
	err = WritePemFile(keyfile, privPemBlock)
	if nil != err {
		return nil, err
	}

	// generate self signed certificate
	cert := &x509.Certificate{
		//		Version:      2,
		SerialNumber: RandomBigInt(),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		Issuer: pkix.Name{
			CommonName: "localhost",
		},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if nil != err {
		return nil, err
	}

	// write create to pem file
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	err = WritePemFile(certfile, certBlock)
	if nil != err {
		return nil, err
	}

	// create tls config
	tlsCfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	return tlsCfg, nil
}

func RandomBigInt() *big.Int {
	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil
	}

	return n
}

func ExpandFilePath(path string) string {
	exp, _ := homedir.Expand(path)
	return exp
}

func WritePemFile(filename string, block *pem.Block) error {
	// check if filename end is .pem
	if !strings.HasSuffix(filename, ".pem") {
		filename = fmt.Sprintf("%s.pem", filename)
	}

	// expand file path
	xfilename := ExpandFilePath(filename)

	// open the file
	pemfile, err := os.Create(xfilename)
	if nil != err {
		return fmt.Errorf("Unable to create file %s! Please check file permissions", xfilename)
	}

	// write pem file
	err = pem.Encode(pemfile, block)
	if nil != err {
		return err
	}

	// write a success message
	pemfile.Close()

	return nil
}
