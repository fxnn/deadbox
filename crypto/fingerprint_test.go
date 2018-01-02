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

	assertTrue(isPassChallenge([]byte{0}, 1), t)
	assertTrue(isPassChallenge([]byte{0}, 0), t)

	assertFalse(isPassChallenge([]byte{1}, 8), t)
	assertTrue(isPassChallenge([]byte{1}, 7), t)
	assertTrue(isPassChallenge([]byte{1}, 1), t)
	assertTrue(isPassChallenge([]byte{1}, 0), t)

	assertFalse(isPassChallenge([]byte{255}, 8), t)
	assertFalse(isPassChallenge([]byte{255}, 1), t)
	assertTrue(isPassChallenge([]byte{255}, 0), t)

	assertTrue(isPassChallenge([]byte{1, 0}, 7), t)
	assertFalse(isPassChallenge([]byte{1, 0}, 8), t)
	assertFalse(isPassChallenge([]byte{1, 0}, 9), t)
	assertFalse(isPassChallenge([]byte{1, 1}, 8), t)

	assertTrue(isPassChallenge([]byte{0, 1}, 8), t)
	assertTrue(isPassChallenge([]byte{0, 1}, 15), t)
	assertFalse(isPassChallenge([]byte{0, 1}, 16), t)

	assertTrue(isPassChallenge([]byte{0, 255}, 8), t)
	assertFalse(isPassChallenge([]byte{0, 255}, 9), t)

	assertTrue(isPassChallenge([]byte{0, 0, 255}, 16), t)
	assertFalse(isPassChallenge([]byte{0, 0, 255}, 17), t)

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

func assertTrue(actual bool, t *testing.T) {
	t.Helper()
	if actual != true {
		t.Fatalf("expected true, but was false")
	}
}

func assertFalse(actual bool, t *testing.T) {
	t.Helper()
	if actual != false {
		t.Fatalf("expected false, but was true")
	}
}
