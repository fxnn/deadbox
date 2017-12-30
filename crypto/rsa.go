package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
)

const (
	pemBlockTypePrivateKey = "RSA PRIVATE KEY"
	rsaKeySize             = 2048
)

func GeneratePublicKeyBytes(privateKey *rsa.PrivateKey) ([]byte, error) {
	return marshalPublicKey(generatePublicKey(privateKey))
}

func UnmarshalPublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	var publicKey rsa.PublicKey

	_, err := asn1.Unmarshal(publicKeyBytes, &publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %s", err)
	}

	return &publicKey, nil
}

func marshalPublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	return asn1.Marshal(*publicKey)
}

func generatePublicKey(privateKey *rsa.PrivateKey) *rsa.PublicKey {
	return &privateKey.PublicKey
}

func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

func MarshalPrivateKeyToPEMBytes(key *rsa.PrivateKey) []byte {
	keyPKCS1Bytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := &pem.Block{
		Type:  pemBlockTypePrivateKey,
		Bytes: keyPKCS1Bytes,
	}
	return pem.EncodeToMemory(keyPEM)
}

func UnmarshalPrivateKeyFromPEMBytes(bytes []byte) (*rsa.PrivateKey, error) {
	if keyBlock, _ := pem.Decode(bytes); keyBlock == nil {
		return nil, fmt.Errorf("could not parse PEM data")
	} else if keyBlock.Type != pemBlockTypePrivateKey {
		return nil, fmt.Errorf("unsupported private key type: %s", keyBlock.Type)
	} else {
		return unmarshalPrivateKeyFromPKCS1Bytes(keyBlock.Bytes)
	}
}

func unmarshalPrivateKeyFromPKCS1Bytes(keyBytes []byte) (*rsa.PrivateKey, error) {
	if privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes); err != nil {
		return nil, fmt.Errorf("could not parse PKCS1 data: %s", err)
	} else {
		return privateKey, nil
	}
}
