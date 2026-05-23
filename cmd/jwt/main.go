package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateSecretKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func main() {
	secret, err := GenerateSecretKey()
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	fmt.Println("Generated Secret Key:", secret)
}
