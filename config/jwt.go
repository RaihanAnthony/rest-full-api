package config

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("ibvoufg9w893ofbvfykbauei9v32v23")

type JWTClaim struct {
	UserName string
	jwt.RegisteredClaims
}