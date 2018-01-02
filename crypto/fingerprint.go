package crypto

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"reflect"
)

const (
	hashSeparator             = "\000\000\000"
	hashFunction              = crypto.SHA256
	fingerprintGroupSeparator = ":"
)

func FingerprintPublicKey(
	key *rsa.PublicKey,
	encryptionType string,
	challengeLevel int,
	fingerprintLengthInGroups int,
) (string, error) {
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)

	keyBytes, err := marshalPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("marshalling public key failed: %s", err)
	}

	var hashSum []byte
	for modifier := 0; !isPassChallenge(hashSum, challengeLevel); modifier++ {
		hashSum, err = generateHashSum(modifier, keyBytes, encryptionType, hashFunction)
		if err != nil {
			return "", fmt.Errorf("calculating hash sum failed: %s", err)
		}
	}

	hashSumString := encoding.EncodeToString(hashSum[challengeLevel:])
	fingerprint := generateGroupedFingerprint(hashSumString,
		fingerprintLengthInGroups,
		fingerprintGroupSeparator)

	return fingerprint, nil
}

func isPassChallenge(hashInput []byte, challengeLevel int) bool {
	if hashInput == nil {
		return false
	}

	for _, b := range hashInput[:challengeLevel] {
		if b != 0 {
			return false
		}
	}

	return true
}

func generateHashSum(modifier int,
	keyBytes []byte,
	encryptionType string,
	hashFunction crypto.Hash,
) ([]byte, error) {
	hashInput, err := generateHashInput(modifier, keyBytes, encryptionType)
	if err != nil {
		return nil, fmt.Errorf("generating hash input failed: %s", err)
	}

	hash := hashFunction.New()
	_, err = hash.Write(hashInput)
	if err != nil {
		hashType := reflect.TypeOf(hash)
		return nil, fmt.Errorf("%s: %s", hashType, err)
	}

	return hash.Sum([]byte{}), nil
}

func generateHashInput(modifier int, keyBytes []byte, encryptionType string) ([]byte, error) {
	var hashInputBuffer bytes.Buffer

	hashInputBuffer.Write(keyBytes)
	hashInputBuffer.WriteString(hashSeparator)

	hashInputBuffer.WriteString(encryptionType)
	hashInputBuffer.WriteString(hashSeparator)

	err := binary.Write(&hashInputBuffer, binary.BigEndian, int64(modifier))
	if err != nil {
		return nil, fmt.Errorf("writing modifier failed: %s", err)
	}

	// @todo #27 add validity time

	return hashInputBuffer.Bytes(), nil
}

func generateGroupedFingerprint(hashSumString string, numberOfGroups int, groupSeparator string) string {
	var groupedFingerprintBuffer bytes.Buffer
	for groupIdx := 0; groupIdx < numberOfGroups; groupIdx++ {
		groupedFingerprintBuffer.WriteString(hashSumString[2*groupIdx : 2*groupIdx+2])

		if groupIdx < numberOfGroups-1 {
			groupedFingerprintBuffer.WriteString(groupSeparator)
		}
	}
	fingerprint := groupedFingerprintBuffer.String()
	return fingerprint
}
