package main

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	UserID string
	jwt.StandardClaims
}
