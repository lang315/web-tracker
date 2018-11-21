package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/thanhps42/tlib/rand"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type AuthRouter struct {
}

func (this *AuthRouter) Install() {
	r := app.Group("/auth")
	r.POST("/login", this.login)
	r.GET("/user", this.user, jwtAuth)
	r.POST("/register", this.register)
	r.GET("/refresh", this.refresh, jwtAuth)
	r.PUT("/password", this.changePassword, jwtAuth)
}

func (*AuthRouter) changePassword(ctx echo.Context) error {
	f := &struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}{}
	if err := ctx.Bind(f); err != nil {
		return err
	}

	u, err := getCurrentUser(ctx)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(f.CurrentPassword)); err != nil {
		return ctx.JSON(http.StatusOK, echo.Map{"error": "Mật khẩu hiện tại không chính xác."})
	}

	buf, err := bcrypt.GenerateFromPassword([]byte(f.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(buf)
	if _, err = db.Model(u).WherePK().Column("password").Update(); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (*AuthRouter) refresh(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*UserClaims)
	token, err := makeUserJWTToken(claims.UserID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, echo.Map{"token": token})
}

func (*AuthRouter) login(ctx echo.Context) error {
	f := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := ctx.Bind(f); err != nil {
		return err
	}
	f.Email = strings.ToLower(f.Email)

	u := &User{}
	if err := db.Model(u).Where("email = ?", f.Email).Select(u); err != nil {
		return echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(f.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	token, err := makeUserJWTToken(u.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{"token": token})
}

func (*AuthRouter) user(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*UserClaims)
	u := &User{ID: claims.UserID}
	if err := db.Select(u); err != nil {
		return err
	}


	return ctx.JSON(http.StatusOK, echo.Map{
		"user": echo.Map{
			"id": claims.UserID,
		},
	})
}

func (*AuthRouter) register(ctx echo.Context) error {
	f := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := ctx.Bind(f); err != nil {
		return err
	}
	f.Email = strings.ToLower(f.Email)

	u := &User{}
	exists, err := db.Model(u).Where("email = ?", f.Email).Exists()
	if err != nil {
		return err
	}
	if exists {
		return ctx.JSON(http.StatusOK, echo.Map{
			"error": "Email này đã được đăng ký.",
		})
	}

	buf, err := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.ID = trand.ULID()
	u.Email = f.Email
	u.Password = string(buf)
	u.Membership = "member"
	u.CreatedAt = time.Now().Unix()
	if err = db.Insert(u); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{})
}

