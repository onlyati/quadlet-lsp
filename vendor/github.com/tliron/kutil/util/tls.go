package util

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"time"
)

func CreateTLSConfig(certificate []byte, key []byte) (*tls.Config, error) {
	if certificate, err := tls.X509KeyPair(certificate, key); err == nil {
		return &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}, nil
	} else {
		return nil, err
	}
}

// Creates a TLS config with a single X.509 certificate signed with a new RSA key pair.
//
// If rsaBits is 0 it will default to 2048. If duration is 0 it will default to one year.
func CreateSelfSignedTLSConfig(organization string, host string, rsaBits int, duration time.Duration) (*tls.Config, error) {
	if privateKey, certificate, err := CreateRSAX509Certificate(organization, host, rsaBits, duration); err == nil {
		return &tls.Config{
			Certificates: []tls.Certificate{
				{
					Certificate: [][]byte{certificate.Raw},
					PrivateKey:  privateKey,
				},
			},
		}, nil
	} else {
		return nil, err
	}
}

func WriteTLSCertificatePEM(writer io.Writer, certificate *tls.Certificate) error {
	if len(certificate.Certificate) > 0 {
		return pem.Encode(writer, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certificate.Certificate[0],
		})
	} else {
		return errors.New("no certificate data")
	}
}

func WriteTLSRSAKeyPEM(writer io.Writer, certificate *tls.Certificate) error {
	if privateKey, ok := certificate.PrivateKey.(*rsa.PrivateKey); ok {
		return pem.Encode(writer, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
	} else {
		return errors.New("no RSA private key")
	}
}
