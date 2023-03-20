package utilities

import (
	"crypto/rand"
	"encoding/hex"
)

func GenVerificationCode() (string, error) {
	b := make([]byte, 2)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
