package crypto

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/fxnn/deadbox/model"
)

const (
	encryptionTypePlain        = "encryptionType:github.com/fxnn/deadbox:plain:1.0"
	encryptionTypeAESPlusRSA   = "encryptionType:github.com/fxnn/deadbox:AESPlusRSA:1.0"
	keyFromCyphertextDelimiter = ":::"
)

func EncryptRequest(content []byte, rsaPublicKey *rsa.PublicKey) (contentEncrypted []byte, encryptionType string, err error) {
	encryptionType = encryptionTypeAESPlusRSA

	ciphertext, aesNonce, aesKey, err := encryptUsingAESPlusRSA(content, rsaPublicKey)
	if err != nil {
		err = fmt.Errorf("failed to encrypt content: %s", err)
		return
	}

	contentEncrypted = []byte(marshalToBase64(ciphertext) + keyFromCyphertextDelimiter +
		marshalToBase64(aesNonce) + keyFromCyphertextDelimiter +
		marshalToBase64(aesKey))

	return
}

func DecryptRequest(request model.WorkerRequest, rsaPrivateKey *rsa.PrivateKey) ([]byte, error) {
	if request.EncryptionType == encryptionTypePlain {
		return request.Content, nil
	}

	if request.EncryptionType != encryptionTypeAESPlusRSA {
		return nil, fmt.Errorf("unsupported encryption type: %s", request.EncryptionType)
	} else {
		ciphertextAndNonceAndKey := strings.Split(string(request.Content), keyFromCyphertextDelimiter)
		if len(ciphertextAndNonceAndKey) != 3 {
			return nil, fmt.Errorf("malformed content")
		}

		ciphertext, err := unmarshalBase64Encoded(ciphertextAndNonceAndKey[0])
		if err != nil {
			return nil, fmt.Errorf("unmarshalling ciphertext failed: %s", err)
		}

		aesNonce, err := unmarshalBase64Encoded(ciphertextAndNonceAndKey[1])
		if err != nil {
			return nil, fmt.Errorf("unmarshalling AES nonce failed: %s", err)
		}

		aesKey, err := unmarshalBase64Encoded(ciphertextAndNonceAndKey[2])
		if err != nil {
			return nil, fmt.Errorf("unmarshalling AES key failed: %s", err)
		}

		return decryptUsingAESPlusRSA(ciphertext, aesNonce, aesKey, rsaPrivateKey)
	}
}

func unmarshalBase64Encoded(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

func marshalToBase64(raw []byte) string {
	return base64.StdEncoding.EncodeToString(raw)
}
