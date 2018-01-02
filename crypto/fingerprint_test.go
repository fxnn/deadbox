package crypto

import (
	"crypto/rsa"
	"math/big"
	"testing"
)

func TestSomething(t *testing.T) {

	key := &rsa.PublicKey{
		N: big.NewInt(42),
		E: 13,
	}
	encryptionType := encryptionTypeAESPlusRSA

	assertFingerprint(t, "6W:HJ:QI:MS:AF:VL:HD:LB", 0, 8, key, encryptionType)
	assertFingerprint(t, "UO:4H:KV:XF:SL:GP:OH:AN", 1, 8, key, encryptionType)
	assertFingerprint(t, "VY:NC:DT:LR:I5:OJ:5B:SC", 2, 8, key, encryptionType)
	assertFingerprint(t, "6S:GL:5D:TN:A4:CJ:GZ:P6", 3, 8, key, encryptionType)

}

func assertFingerprint(
	t *testing.T,
	expected string,
	challengeLevel int,
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

func TestIsPassChallenge(t *testing.T) {

	assertTrue(isPassChallenge([]byte{0}, 1), t)
	assertTrue(isPassChallenge([]byte{0}, 0), t)

	assertFalse(isPassChallenge([]byte{1}, 1), t)
	assertTrue(isPassChallenge([]byte{1}, 0), t)

	assertFalse(isPassChallenge([]byte{1, 0}, 2), t)
	assertFalse(isPassChallenge([]byte{1, 0}, 1), t)
	assertFalse(isPassChallenge([]byte{1, 1}, 1), t)
	assertTrue(isPassChallenge([]byte{1, 0}, 0), t)
	assertTrue(isPassChallenge([]byte{0, 1}, 1), t)
	assertTrue(isPassChallenge([]byte{0, 0}, 2), t)

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
