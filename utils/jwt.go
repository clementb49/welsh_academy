// The package 'utils' contains utility functions used throughout the application
package utils

import (
	"strconv"
	"time"

	"github.com/clementb49/welsh_academy/config"
	"github.com/golang-jwt/jwt/v4"
)
// Define the issuer of the jwt claims
const JWT_ISSUER = "welsh_academy.fr"

// Global variable to store a cached version of the signinn key
var signingKey []byte
// getJwtClaim returns the claims associated with the JWT token.
func getJwtClaim(userId uint) *jwt.RegisteredClaims {
	now := time.Now()
	expiresAt := now.Add(30 * time.Minute)
	return &jwt.RegisteredClaims{
		Issuer:    JWT_ISSUER,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ID:        strconv.Itoa(int(userId)),
	}
}

// initSigningKey initializes the signing key used for JWT token signing. It get the signing key from application configuration. 
func initSigningKey() {
	if len(signingKey) == 0 {
		config := config.GetWaConfig()
		signingKey = []byte(config.JwtKey)
	}
}
// GetSignedJwt generates and returns a signed JWT token string for a given user ID.
func GetSignedJwt(userId uint) (string, error) {
	tokenCfg := jwt.NewWithClaims(jwt.SigningMethodHS256, getJwtClaim(userId))
	initSigningKey()
	token, err := tokenCfg.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyToken verifies the signature and validity of a given JWT token.
func VerifyToken(accessToken string) (*jwt.RegisteredClaims, error) {
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(accessToken, &claims, func(t *jwt.Token) (interface{}, error) {
		initSigningKey()
		return signingKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
