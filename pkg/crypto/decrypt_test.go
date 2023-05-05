package crypto_test

import (
	"fmt"
	"testing"

	"github.com/CafeKetab/book/pkg/crypto"
)

func TestDecrypt(t *testing.T) {
	plainText, secret := "plainText", "A?D(G-KaPdSgVkYp"

	encrypted, err := crypto.Encrypt(plainText, secret)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(encrypted)

	decrypted, err := crypto.Decrypt(encrypted, secret)
	if err != nil {
		t.Error(err)
	}

	if plainText != decrypted {
		t.Errorf("expected: %s, recieved: %s\n", plainText, decrypted)
	}
}
