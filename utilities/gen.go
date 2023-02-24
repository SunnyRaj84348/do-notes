package utilities

import (
	"crypto/rand"
	"log"
)

// Generate 64 secure random numbers
func RandToken() []byte {
	b := make([]byte, 64)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}
