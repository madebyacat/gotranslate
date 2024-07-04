package models

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthConfig struct {
	JwtKey             []byte
	Username           string
	Password           string
	SkipAuthentication bool
}

func NewAuthConfig(skip bool, jwtKey, username, password string) AuthConfig {
	return AuthConfig{
		SkipAuthentication: skip,
		JwtKey:             []byte(jwtKey),
		Username:           username,
		Password:           password,
	}
}
