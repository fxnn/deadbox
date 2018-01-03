package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	certificateOrganization = "github.com/fxnn/deadbox self-signed certificate"
	pemBlockTypeCertificate = "CERTIFICATE"
)

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

func GenerateCertificateBytes(privateKey *rsa.PrivateKey, validFor time.Duration, hosts []string) ([]byte, error) {
	certificate, err := createCertificate(certificateOrganization, validFor, hosts, false)
	if err != nil {
		return nil, err
	}

	return marshalSelfSignedCertificateToPEMBytes(privateKey, certificate)
}

func marshalSelfSignedCertificateToPEMBytes(privateKey *rsa.PrivateKey, certificate *x509.Certificate) ([]byte, error) {
	return marshalCertificateToPEMBytes(privateKey, certificate, &privateKey.PublicKey, certificate)
}

func marshalCertificateToPEMBytes(signingPrivateKey *rsa.PrivateKey, signingCertificate *x509.Certificate, signeePublicKey *rsa.PublicKey, signeeCertificate *x509.Certificate) ([]byte, error) {
	derBytes, err := x509.CreateCertificate(rand.Reader, signeeCertificate, signingCertificate, signeePublicKey, signingPrivateKey)
	if err != nil {
		return nil, err
	}

	pemBlock := &pem.Block{
		Type:  pemBlockTypeCertificate,
		Bytes: derBytes,
	}

	return pem.EncodeToMemory(pemBlock), nil
}

func createCertificate(organization string, validFor time.Duration, hosts []string, isCA bool) (*x509.Certificate, error) {
	// based on https://golang.org/src/crypto/tls/generate_cert.go

	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	certificate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{certificateOrganization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			certificate.IPAddresses = append(certificate.IPAddresses, ip)
		} else {
			certificate.DNSNames = append(certificate.DNSNames, host)
		}
	}

	if isCA {
		certificate.IsCA = true
		certificate.KeyUsage |= x509.KeyUsageCertSign
	}

	return &certificate, nil
}
