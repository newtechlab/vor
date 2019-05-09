package main

import (
	"encoding/base64"

	"golang.org/x/crypto/sha3"
)

func generateID(phone string) string {
	if len(fSalt) < 32 {
		panic("could not read sufficiently large salt")
	}
	hash := sha3.New512()
	hash.Write([]byte(fSalt))
	hash.Write([]byte(phone))
	sum := hash.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}
