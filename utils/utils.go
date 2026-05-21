package utils

import (
	"crypto/rand"
)

func NewSecret(length uint32) ([]byte, error){
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
