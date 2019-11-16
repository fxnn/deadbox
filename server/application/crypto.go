package application

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fxnn/deadbox/server/crypto"
	"github.com/fxnn/deadbox/server/rest"
)

func ReadOrCreateTLSCertFiles(
	certPath string,
	privateKeyPath string,
	name string,
	privateKeySize int,
	certificateValidFor time.Duration,
	hosts []string,
	fingerprintLength uint,
	fingerprintChallengeLevel uint,
) (rest.TLS, error) {
	privateKeyFile, privateKey, err := ReadOrCreatePrivateKey(privateKeyPath, name, privateKeySize, fingerprintLength, fingerprintChallengeLevel)
	if err != nil {
		return nil, err
	}

	certFile := certFileName(certPath, name)
	bytes, err := ioutil.ReadFile(certFile)
	if os.IsNotExist(err) {
		log.Printf("'%s' has no certificate, generating one", name)

		bytes, err = crypto.GenerateCertificateBytes(privateKey, certificateValidFor, hosts)
		if err != nil {
			err = fmt.Errorf("couldn't generate certificate: %s", err)
			return nil, err
		}

		err = ioutil.WriteFile(certFile, bytes, filePermOnlyUserCanReadOrWrite)
		if err != nil {
			err = fmt.Errorf("couldn't write generated certificate to file %s: %s", certFile, err)
			return nil, err
		}
	} else if err != nil {
		return nil, fmt.Errorf("couldn't read certificate file %s: %s", certFile, err)
	}

	if cert, err := crypto.UnmarshalCertificateFromPEMBytes(bytes); err != nil {
		return nil, fmt.Errorf("couldn't parse certificate from file %s: %s", certFile, err)
	} else if certRsaPublicKey, ok := cert.PublicKey.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("certificate from %s has unsupported public key (only RSA supported)", certFile)
	} else if certRsaPublicKey.N.Cmp(privateKey.N) != 0 || certRsaPublicKey.E != privateKey.E {
		return nil, fmt.Errorf("certificate from %s has different public key than private key from %s", certFile, privateKeyFile)
	}

	return rest.NewFileBasedTLS(privateKeyFile, certFile), nil
}

func ReadOrCreatePrivateKey(
	privateKeyPath string,
	name string,
	privateKeySize int,
	fingerprintLength uint,
	fingerprintChallengeLevel uint,
) (fileName string, privateKey *rsa.PrivateKey, err error) {
	var bytes []byte

	fileName = privateKeyFileName(privateKeyPath, name)
	bytes, err = ioutil.ReadFile(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			err = fmt.Errorf("couldn't read file %s: %s", fileName, err)
			return
		}

		log.Printf("'%s' has no private key, generating one", name)
		privateKey, err = crypto.GeneratePrivateKey(privateKeySize)
		if err != nil {
			err = fmt.Errorf("couldn't generate private key: %s", err)
			return
		}

		bytes = crypto.MarshalPrivateKeyToPEMBytes(privateKey)

		err = ioutil.WriteFile(fileName, bytes, filePermOnlyUserCanReadOrWrite)
		if err != nil {
			err = fmt.Errorf("couldn't write generated private key to file %s: %s", fileName, err)
			return
		}
	} else if privateKey, err = crypto.UnmarshalPrivateKeyFromPEMBytes(bytes); err != nil {
		err = fmt.Errorf("couldn't read private key from file %s: %s", fileName, err)
		return
	}

	if fingerprint, err := crypto.FingerprintPublicKey(&privateKey.PublicKey, fingerprintChallengeLevel, fingerprintLength); err != nil {
		log.Printf("couldn't calculate fingerprint for '%s': %s", name, err)
	} else {
		log.Printf("'%s' has public key fingerprint %s", name, fingerprint)
	}

	if privateKey.N.BitLen() != privateKeySize {
		log.Printf("'%s' has configured key size '%d', but existing key has size '%d'",
			name, privateKeySize, privateKey.N.BitLen())
	}

	return
}
