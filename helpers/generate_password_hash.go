package helpers

import (
	"crypto/sha256"
	"fmt"
)

const sasalt = "h0q12hqw124f17ajf3ajs"

func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(sasalt)))
}
