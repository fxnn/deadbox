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

var fingerprintEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

func FingerprintPublicKey(
	key *rsa.PublicKey,
	challengeLevel uint,
	fingerprintLengthInGroups int,
) (string, error) {
	keyBytes, err := marshalPublicKey(key)
	if err != nil {
		return "", fmt.Errorf("marshalling public key failed: %s", err)
	}

	hashInputPrefix, err := generateHashInputPrefix(challengeLevel, keyBytes)
	if err != nil {
		return "", fmt.Errorf("generating hash input failed: %s", err)
	}

	hashSum, _, err := findChallengeSolution(hashInputPrefix, hashFunction, challengeLevel)
	if err != nil {
		return "", fmt.Errorf("generating hash sum failed: %s", err)
	}

	numberOfOmittedBytes := (challengeLevel + 7) / 8
	hashSumString := fingerprintEncoding.EncodeToString(hashSum[numberOfOmittedBytes:])
	fingerprint := generateGroupedFingerprint(hashSumString,
		fingerprintLengthInGroups,
		fingerprintGroupSeparator)

	return fingerprint, nil
}

func findChallengeSolution(
	hashInputPrefix []byte,
	hashFunction crypto.Hash,
	challengeLevel uint,
) (hashSum []byte, challengeSolution int, err error) {
	var hashInputSuffix = make([]byte, 8)               // HINT: will be filled each round
	var zeroHashSum = make([]byte, hashFunction.Size()) // HINT: used for comparison later

	for challengeSolution = 0; !isPassChallenge(zeroHashSum, hashSum, challengeLevel); challengeSolution++ {
		binary.BigEndian.PutUint64(hashInputSuffix, uint64(challengeSolution))
		hashInput := append(hashInputPrefix, hashInputSuffix...)

		hashSum, err = generateHashSum(hashInput, hashFunction)
		if err != nil {
			return
		}
	}

	return
}

func isPassChallenge(zeroHashSum []byte, hashInput []byte, challengeLevel uint) bool {
	if hashInput == nil {
		return false
	}
	if challengeLevel == 0 {
		return true
	}

	idxOfFirstNonZeroByte := challengeLevel / 8 // note, that '/' is always floor'd
	if !bytes.Equal(zeroHashSum[:idxOfFirstNonZeroByte], hashInput[:idxOfFirstNonZeroByte]) {
		return false
	}

	if uint(len(hashInput)) > idxOfFirstNonZeroByte {
		lastByteRequiredToContainZeroBit := hashInput[idxOfFirstNonZeroByte]
		numberOfRequiredZeroBits := challengeLevel % 8
		shouldBeZero := lastByteRequiredToContainZeroBit >> (8 - numberOfRequiredZeroBits)
		return shouldBeZero == 0
	}

	return true
}

func generateHashSum(
	hashInput []byte,
	hashFunction crypto.Hash,
) ([]byte, error) {
	hash := hashFunction.New()
	_, err := hash.Write(hashInput)
	if err != nil {
		hashType := reflect.TypeOf(hash)
		return nil, fmt.Errorf("%s: %s", hashType, err)
	}

	return hash.Sum([]byte{}), nil
}

func generateHashInputPrefix(
	challengeLevel uint,
	keyBytes []byte,
) ([]byte, error) {
	var hashInputBuffer bytes.Buffer

	hashInputBuffer.Write(keyBytes)
	hashInputBuffer.WriteString(hashSeparator)

	if err := binary.Write(&hashInputBuffer, binary.BigEndian, int64(challengeLevel)); err != nil {
		return nil, fmt.Errorf("writing challengeLevel failed: %s", err)
	}
	hashInputBuffer.WriteString(hashSeparator)

	// @todo #27 add validity time

	// NOTE: challengeSolution is appended here in each round
	return hashInputBuffer.Bytes(), nil
}

func generateGroupedFingerprint(hashSumString string, numberOfGroups int, groupSeparator string) string {
	var groupedFingerprintBuffer bytes.Buffer
	for groupIdx := 0; groupIdx < numberOfGroups && 2*groupIdx+2 <= len(hashSumString); groupIdx++ {
		groupedFingerprintBuffer.WriteString(hashSumString[2*groupIdx : 2*groupIdx+2])

		if groupIdx < numberOfGroups-1 {
			groupedFingerprintBuffer.WriteString(groupSeparator)
		}
	}
	fingerprint := groupedFingerprintBuffer.String()
	return fingerprint
}
