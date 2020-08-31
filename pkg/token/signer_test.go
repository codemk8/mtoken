package token

import (
	"reflect"
	"testing"

	"gopkg.in/square/go-jose.v2/jwt"
)

func TestNewSigner(t *testing.T) {
	type args struct {
		privFile string
		pubFile  string
	}
	tests := []struct {
		name    string
		args    args
		want    Signer
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				privFile: "/tmp/rsa.priv",
				pubFile:  "/tmp/rsa.pub",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv, pub := GenerateRsaKeyPairIfNotExist(tt.args.privFile, tt.args.pubFile, ".", false)
			_, err := NewSigner(priv, pub)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSigner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSignerVerifier(t *testing.T) {
	priv, pub := GenerateRsaKeyPairIfNotExist("/tmp/rsa.priv", "/tmp/rsa.pub", ".", false)
	signer, err := NewSigner(priv, pub)
	if err != nil {
		panic(err)
	}
	verifier, err := NewVerifier(pub)
	if err != nil {
		panic(err)
	}

	content := jwt.Claims{
		Issuer:  "issuer1",
		Subject: "subject1",
		ID:      "id1"}
	j, err := signer.Sign(&content)
	if err != nil {
		panic(err)
	}

	repContent, err := verifier.Verify(j)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(*repContent, content) {
		t.Errorf("Claim not reproduced got = %v, want %v", repContent, content)
	}
}
