package main

import (
	"context"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// LoadDefaultPrivateKey loads the default private key from keys/privateKey.pem (relative to app root)
func LoadDefaultPrivateKey() (string, error) {
	keyPath := "keys/privateKey.pem"

	data, err := os.ReadFile(keyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("private key file not found")
		}
		return "", fmt.Errorf("error reading private key file: %w", err)
	}

	return string(data), nil
}

// GenerateHexEncodedSHA256 takes in msg and returns
// the hex coded string after hashing it using SHA256
func GenerateHexEncodedSHA256(msg string) string {
	h := sha256.New()
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

// GetPKCS8RSAPrivateKey takes in an RSA private key string.
// It tries to decode the PEM string and constructs an RSA private key
// object assuming the key is in PKCS8format. The rsaPrivateKey object is returned
// if successful, else an error is returned.
func GetPKCS8RSAPrivateKey(ctx context.Context, pKey string) (rsaPrivateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(pKey))
	if block == nil {
		err = errors.New("ssh: no key found")
		return
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	rsaPrivateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		err = errors.New("cannot type assert as rsa private key pointer")
		return
	}

	return
}

// GenerateBase64EncodedSHA256withRSA takes in a msg as slice of string
// along with an rsaPrivateKey Object. The msg is hashed using SHA256 and
// signed using the RSA private key
func GenerateBase64EncodedSHA256withRSA(ctx context.Context, message []byte, rsaPrivateKey *rsa.PrivateKey) (encryptedMsg string, err error) {
	randomReader := rand.Reader

	// SHA-256 encode the payload
	hashed := sha256.Sum256(message)
	// sign using rsa private key
	byteSign, err := rsa.SignPKCS1v15(randomReader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return
	}

	// base-64 encode the signed message
	encryptedMsg = base64.StdEncoding.EncodeToString([]byte(byteSign))
	return
}

// GenerateBase64EncodedHMAC512WithSecretKey generates HMAC-SHA512 signature
func GenerateBase64EncodedHMAC512WithSecretKey(secretKey string, data []byte) (string, error) {
	hash := hmac.New(sha512.New, []byte(secretKey))
	_, err := hash.Write(data)
	if err != nil {
		return "", err
	}

	signature := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature), nil
}
