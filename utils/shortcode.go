package utils

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

func GenerateShortCode() (string, error) {
	code := make([]byte, codeLength)

	randomBytes := make([]byte, codeLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	for i := range code {
		code[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return string(code), nil
}
