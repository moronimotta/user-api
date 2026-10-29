package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
)

func LoadRSAPrivateKey() (*rsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing RSA private key")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}

func ValidateJWTToken(tokenString string) bool {
	_, err := ParseJWTTokenClaims(tokenString)
	return err == nil
}

// ParseJWTTokenClaims validates a JWT signature and returns the claims from a
// token signed with the configured RSA public key.
func ParseJWTTokenClaims(tokenString string) (jwt.MapClaims, error) {
	publicKeyData, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing RSA public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		if err == nil {
			err = errors.New("invalid token")
		}
		return nil, err
	}
	return claims, nil
}
