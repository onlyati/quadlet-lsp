package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

func ParseX509Certificates(bytes []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate

	for len(bytes) > 0 {
		var block *pem.Block
		block, bytes = pem.Decode(bytes)
		if block != nil {
			if certificate, err := x509.ParseCertificate(block.Bytes); err == nil {
				certificates = append(certificates, certificate)
			} else {
				return nil, err
			}
		} else {
			bytes = nil
		}
	}

	return certificates, nil
}

func ParseX509CertificatePool(bytes []byte) (*x509.CertPool, error) {
	if certificates, err := ParseX509Certificates(bytes); err == nil {
		if len(certificates) > 0 {
			certPool := x509.NewCertPool()
			for _, certificate := range certificates {
				certPool.AddCert(certificate)
			}
			return certPool, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

func RandomSerialNumber() (*big.Int, error) {
	return rand.Int(rand.Reader, serialNumberLimit)
}

// Initializes a X.509 certificate with a random serial number.
//
// If duration is 0 it will default to one year.
func NewX509Certificate(organization string, host string, duration time.Duration, rsa bool, ca bool) (*x509.Certificate, error) {
	// See: https://golang.org/src/crypto/tls/generate_cert.go

	if serialNumber, err := RandomSerialNumber(); err == nil {
		now := time.Now()

		if duration == 0 {
			duration = 365 * 24 * time.Hour // 1 year
		}

		certificate := x509.Certificate{
			Subject: pkix.Name{
				Organization: []string{organization},
			},
			DNSNames:     []string{host},
			SerialNumber: serialNumber,
			KeyUsage:     x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
			},
			//SignatureAlgorithm:    x509.SHA256WithRSA,
			BasicConstraintsValid: true,
			NotBefore:             now,
			NotAfter:              now.Add(duration),
		}

		if rsa {
			certificate.KeyUsage |= x509.KeyUsageKeyEncipherment
		}

		if ca {
			certificate.IsCA = true
			certificate.KeyUsage |= x509.KeyUsageCertSign
		}

		return &certificate, nil
	} else {
		return nil, err
	}
}

// Signs a X.509 certificate.
//
// For RSA, privateKey should be [*rsa.PrivateKey] and publicKey should be [*rsa.PublicKey].
func SignX509Certificate(certificate *x509.Certificate, privateKey any, publicKey any) (*x509.Certificate, error) {
	if certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, publicKey, privateKey); err == nil {
		return x509.ParseCertificate(certificateBytes)
	} else {
		return nil, err
	}
}

// Generates a random RSA key pair and uses it to sign a X.509 certificate
// with a random serial number.
//
// If host is empy it will default to "localhost". If rsaBits is 0 it will default to 2048.
// If duration is 0 it will default to one year.
func CreateRSAX509Certificate(organization string, host string, rsaBits int, duration time.Duration) (*rsa.PrivateKey, *x509.Certificate, error) {
	if host == "" {
		host = "localhost"
	}

	if rsaBits == 0 {
		rsaBits = 2048
	}

	if privateKey, err := rsa.GenerateKey(rand.Reader, rsaBits); err == nil {
		if certificate, err := NewX509Certificate(organization, host, duration, true, true); err == nil {
			if certificate, err = SignX509Certificate(certificate, privateKey, &privateKey.PublicKey); err == nil {
				return privateKey, certificate, nil
			} else {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}
}
