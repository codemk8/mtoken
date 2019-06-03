package token

import (
	"crypto/rsa"

	"gopkg.in/square/go-jose.v2/jwt"
)

type Verifier interface {
	// Sign a token and return the serialized cryptographic token.
	Verify(token *string) (*jwt.Claims, error)
}

type joseVerifier struct {
	pubKey *rsa.PublicKey
}

func NewVerifier(pubKey *rsa.PublicKey) (Verifier, error) {

	return &joseVerifier{pubKey: pubKey}, nil
}

func (j *joseVerifier) Verify(token *string) (*jwt.Claims, error) {
	parsedJWT, err := jwt.ParseSigned(*token)
	if err != nil {
		return nil, err
	}
	resultCl := jwt.Claims{}
	err = parsedJWT.Claims(j.pubKey, &resultCl)
	return &resultCl, err
}
