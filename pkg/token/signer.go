package token

import (
	"crypto/rsa"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Signer interface {
	// Sign a token and return the serialized cryptographic token.
	Sign(claims *jwt.Claims) (*string, error)
}

type joseSigner struct {
	signer jose.Signer
}

func NewSigner(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) (Signer, error) {
	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	//publicKey := &privateKey.PublicKey
	//key := []byte("secret")
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privKey}, (&jose.SignerOptions{}).WithType("JWT"))

	if err != nil {
		panic(err)
	}

	return &joseSigner{signer: sig}, nil
}

func (j *joseSigner) Sign(claims *jwt.Claims) (*string, error) {
	raw, err := jwt.Signed(j.signer).Claims(claims).CompactSerialize()
	if err != nil {
		panic(err)
	}
	return &raw, nil
}
