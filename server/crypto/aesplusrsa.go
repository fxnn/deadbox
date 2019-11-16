package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
)

const (
	aesKeySize = 32
)

func encryptUsingAESPlusRSA(plaintext []byte, rsaPublicKey *rsa.PublicKey) (ciphertext []byte, aesNonce []byte, aesKeyEncrypted []byte, err error) {
	var rng = rand.Reader

	aesKey := make([]byte, aesKeySize)
	_, err = io.ReadFull(rng, aesKey)
	if err != nil {
		err = fmt.Errorf("generating AES key failed: %s", err)
		return
	}

	aesKeyEncrypted, err = encryptSessionKeyUsingRSA(aesKey, rsaPublicKey)
	if err != nil {
		err = fmt.Errorf("encrypting AES key failed: %s", err)
		return
	}

	aesCipher, err := createAESCipher(aesKey)
	if err != nil {
		err = fmt.Errorf("initializing AES cipher failed: %s", err)
	}

	aesNonce = make([]byte, aesCipher.NonceSize())
	_, err = io.ReadFull(rng, aesNonce)
	if err != nil {
		err = fmt.Errorf("generating AES nonce failed: %s", err)
		return
	}

	ciphertext = aesCipher.Seal(nil, aesNonce, plaintext, nil)
	return
}

func decryptUsingAESPlusRSA(ciphertext []byte, aesNonce []byte, aesKeyEncrypted []byte, rsaPrivateKey *rsa.PrivateKey) ([]byte, error) {
	aesKey, err := decryptSessionKeyUsingRSA(aesKeyEncrypted, rsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("decryption of AES key failed: %s", err)
	}

	aesCipher, err := createAESCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("initializing AES cipher failed: %s", err)
	}

	if len(aesNonce) != aesCipher.NonceSize() {
		return nil, fmt.Errorf("given AES nonce of size %d does not match requested size %d", len(aesNonce), aesCipher.NonceSize())
	}

	plaintext, err := aesCipher.Open(nil, aesNonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("running block cipher failed: %s", err)
	}

	return plaintext, nil
}

func createAESCipher(aesKey []byte) (cipher.AEAD, error) {
	blockCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("AES Cipher: %s", err)
	}

	blockCipherMode, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("GCM: %s", err)
	}

	return blockCipherMode, nil
}

func encryptSessionKeyUsingRSA(plaintext []byte, rsaPublicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, plaintext)
}

func decryptSessionKeyUsingRSA(rsaCyphertext []byte, rsaPrivateKey *rsa.PrivateKey) ([]byte, error) {
	rng := rand.Reader

	// NOTE, that this random key will be returned in case of a wrong rsaPrivateKey, making it harder for attackers
	// to guess the key
	resultKey := make([]byte, aesKeySize)
	if _, err := io.ReadFull(rng, resultKey); err != nil {
		return nil, fmt.Errorf("random key creation failed: %s", err)
	}

	if err := rsa.DecryptPKCS1v15SessionKey(rng, rsaPrivateKey, rsaCyphertext, resultKey); err != nil {
		// NOTE, that returned errors may be disclosed, i.e. they don't contain sensible information
		return nil, fmt.Errorf("RSA decryption failed: %s", err)
	}

	return resultKey, nil
}
