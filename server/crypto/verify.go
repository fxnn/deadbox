package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strings"
)

type VerifyByFingerprint struct {
	Fingerprint               string
	FingerprintChallengeLevel uint
	FingerprintLength         uint
}

// VerifyPeerCertificate validates that one of the given certificates contains a public key matching the configured
// fingerprint.
func (v *VerifyByFingerprint) VerifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	errMessages := make([]string, len(rawCerts))
	for _, rawCert := range rawCerts {
		err := v.verifyRawCert(rawCert)
		if err == nil {
			fmt.Println("verify end")
			return nil
		}
		errMessages = append(errMessages, err.Error())
	}

	return fmt.Errorf("given certificates not valid: %s", strings.Join(errMessages, ","))
}

func (v *VerifyByFingerprint) verifyRawCert(rawCert []byte) error {
	if cert, err := x509.ParseCertificate(rawCert); err != nil {
		return err
	} else if rsaPublicKey, ok := cert.PublicKey.(*rsa.PublicKey); !ok {
		return fmt.Errorf("given certificate must be RSA")
	} else {
		// @todo #31: check validity, SAN/CommonName, extensions
		return v.verifyPublicKey(rsaPublicKey)
	}
}

func (v *VerifyByFingerprint) verifyPublicKey(publicKey *rsa.PublicKey) error {
	fingerprint, err := FingerprintPublicKey(publicKey, v.FingerprintChallengeLevel, v.FingerprintLength)
	if err != nil {
		return fmt.Errorf("fingerprint cannot be generated: %s", err)
	}

	if fingerprint != v.Fingerprint {
		return fmt.Errorf("actual fingerprint %s differs from configured one", fingerprint)
	}

	return nil
}
