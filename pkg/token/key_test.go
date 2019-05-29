package token

import (
	"crypto/rsa"

	"github.com/stretchr/testify/assert"

	// "reflect"
	"testing"
)

func TestGenerateRsaKeyPairIfNotExist(t *testing.T) {
	type args struct {
		privKeyFile string
		pubKeyFile  string
	}
	tests := []struct {
		name  string
		args  args
		want  *rsa.PrivateKey
		want1 *rsa.PublicKey
	}{
		// TODO: Add test cases.
		{
			name: "walkthrough",
			args: args{privKeyFile: "/tmp/rsa.priv",
				pubKeyFile: "/tmp/rsa.pub"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GenerateRsaKeyPairIfNotExist(tt.args.privKeyFile, tt.args.pubKeyFile, false)
			assert.True(t, got != nil)
			assert.True(t, got1 != nil)
			/*
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GenerateRsaKeyPairIfNotExist() got = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got1, tt.want1) {
					t.Errorf("GenerateRsaKeyPairIfNotExist() got1 = %v, want %v", got1, tt.want1)
				}
			*/
		})
	}
}
