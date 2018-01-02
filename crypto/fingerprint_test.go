package crypto

import (
	"crypto/rsa"
	"math/big"
	"testing"
)

func TestFingerprintPublicKey(t *testing.T) {

	key := &rsa.PublicKey{
		N: big.NewInt(42),
		E: 13,
	}
	encryptionType := encryptionTypeAESPlusRSA

	assertFingerprint(t, "QX:HG:W6:YO:R2:AC:N4:R3", 0, 8, key, encryptionType)
	assertFingerprint(t, "BC:5X:HY:BO:VU:IB:IW:QC", 4, 8, key, encryptionType)
	assertFingerprint(t, "GA:PB:TP:LR:WY:YU:TG:C7", 8, 8, key, encryptionType)
	assertFingerprint(t, "4B:NN:B2:63:ZY:IK:UF:XG", 16, 8, key, encryptionType)
	assertFingerprint(t, "YL:IK:5C:KX:6B:LO:Z7:TN", 21, 8, key, encryptionType)

}

func TestIsPassChallenge(t *testing.T) {

	assertPassChallenge([]byte{0}, 1, t)
	assertPassChallenge([]byte{0}, 0, t)

	assertDoesntPassChallenge([]byte{1}, 8, t)
	assertPassChallenge([]byte{1}, 7, t)
	assertPassChallenge([]byte{1}, 1, t)
	assertPassChallenge([]byte{1}, 0, t)

	assertDoesntPassChallenge([]byte{255}, 8, t)
	assertDoesntPassChallenge([]byte{255}, 1, t)
	assertPassChallenge([]byte{255}, 0, t)

	assertPassChallenge([]byte{1, 0}, 7, t)
	assertDoesntPassChallenge([]byte{1, 0}, 8, t)
	assertDoesntPassChallenge([]byte{1, 0}, 9, t)
	assertDoesntPassChallenge([]byte{1, 1}, 8, t)

	assertPassChallenge([]byte{0, 1}, 8, t)
	assertPassChallenge([]byte{0, 1}, 15, t)
	assertDoesntPassChallenge([]byte{0, 1}, 16, t)

	assertPassChallenge([]byte{0, 255}, 8, t)
	assertDoesntPassChallenge([]byte{0, 255}, 9, t)

	assertPassChallenge([]byte{0, 0, 255}, 16, t)
	assertDoesntPassChallenge([]byte{0, 0, 255}, 17, t)

}

func TestGenerateGroupedFingerprint(t *testing.T) {

	assertGroupedFingerprint("", "", 0, t)
	assertGroupedFingerprint("", "", 1, t)

	assertGroupedFingerprint("", "ab", 0, t)
	assertGroupedFingerprint("ab", "ab", 1, t)
	assertGroupedFingerprint("ab:", "ab", 2, t)

	assertGroupedFingerprint("", "abcd", 0, t)
	assertGroupedFingerprint("ab", "abcd", 1, t)
	assertGroupedFingerprint("ab:cd", "abcd", 2, t)
	assertGroupedFingerprint("ab:cd:", "abcd", 3, t)

}

func assertFingerprint(
	t *testing.T,
	expected string,
	challengeLevel uint,
	fingerprintLengthInGroups int,
	key *rsa.PublicKey,
	encryptionType string,
) {
	t.Helper()

	fingerprint, err := FingerprintPublicKey(key, encryptionType, challengeLevel, fingerprintLengthInGroups)
	if err != nil {
		t.Fatalf("generating fingerprint failed: %s", err)
	}

	if fingerprint != expected {
		t.Fatalf("unexpected fingerprint: %s", fingerprint)
	}
}

func assertPassChallenge(hashInput []byte, challengeLevel uint, t *testing.T) {
	t.Helper()
	zeroHashSum := make([]byte, len(hashInput))
	if !isPassChallenge(zeroHashSum, hashInput, challengeLevel) {
		t.Fatalf("expected true, but was false")
	}
}

func assertDoesntPassChallenge(hashInput []byte, challengeLevel uint, t *testing.T) {
	t.Helper()
	zeroHashSum := make([]byte, len(hashInput))
	if isPassChallenge(zeroHashSum, hashInput, challengeLevel) {
		t.Fatalf("expected false, but was true")
	}
}

func assertGroupedFingerprint(expected string, input string, numberOfGroups int, t *testing.T) {
	t.Helper()
	actual := generateGroupedFingerprint(input, numberOfGroups, ":")
	if actual != expected {
		t.Fatalf("expected '%s', but was '%s'", expected, actual)
	}
}
