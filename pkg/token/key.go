package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"

	"log"
)

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

func GenerateRsaKeyPairIfNotExist(privKeyFile string, pubKeyFile string, keyDir string, save bool) (*rsa.PrivateKey, *rsa.PublicKey) {
	if privKeyFile == "" {
		privKeyFile = keyDir + "/rsa.priv"
	}
	if pubKeyFile == "" {
		pubKeyFile = keyDir + "/rsa.pub"
	}
	found := true
	if _, err := os.Stat(privKeyFile); os.IsNotExist(err) {
		found = false
	}
	if _, err := os.Stat(pubKeyFile); os.IsNotExist(err) {
		found = false
	}
	if !found {
		log.Printf("Rsa Key files (%s, %s) not found, regenerating.", privKeyFile, pubKeyFile)
		priv, pub := GenerateRsaKeyPair()
		privStr := ExportRsaPrivateKeyAsPemStr(priv)
		pubStr, _ := ExportRsaPublicKeyAsPemStr(pub)
		if save {
			fpri, err := os.Create(privKeyFile)
			if err != nil {
				panic(err)
			}
			defer fpri.Close()
			fpri.WriteString(privStr)

			fpub, err := os.Create(pubKeyFile)
			if err != nil {
				panic(err)
			}
			defer fpub.Close()
			fpub.WriteString(pubStr)
			log.Printf("Saving RSA key pairs to %s and %s.", privKeyFile, pubKeyFile)
		}
		return priv, pub
	}
	priStr, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		panic(err)
	}
	pubStr, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		panic(err)
	}
	log.Printf("Reading RSA key pairs from %s and %s.", privKeyFile, pubKeyFile)
	priv, err := ParseRsaPrivateKeyFromPemStr(string(priStr))
	if err != nil {
		panic(err)
	}
	pub, err := ParseRsaPublicKeyFromPemStr(string(pubStr))
	if err != nil {
		panic(err)
	}
	return priv, pub

}
