package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"time"
)

func makeUserJWTToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	})
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func getUserID(ctx echo.Context) string {
	u := ctx.Get("user").(*jwt.Token)
	claims := u.Claims.(*UserClaims)
	return claims.UserID
}

func getCurrentUser(ctx echo.Context) (*User, error) {
	u := &User{ID: getUserID(ctx)}
	if err := db.Select(u); err != nil {
		return nil, err
	}
	return u, nil
}

