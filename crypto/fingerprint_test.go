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

	assertFingerprint(t, "6W:HJ:QI:MS:AF:VL:HD:LB", 0, 8, key, encryptionType)
	assertFingerprint(t, "AX:MZ:N4:UJ:JZ:B7:EF:US", 4, 8, key, encryptionType)
	assertFingerprint(t, "UO:4H:KV:XF:SL:GP:OH:AN", 8, 8, key, encryptionType)
	assertFingerprint(t, "VY:NC:DT:LR:I5:OJ:5B:SC", 16, 8, key, encryptionType)
	assertFingerprint(t, "67:R3:EW:3E:JI:FE:2V:FB", 21, 8, key, encryptionType)

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
