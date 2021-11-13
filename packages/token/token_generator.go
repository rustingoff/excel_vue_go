package token

import (
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	_Cost     = 16
	tokenTTL  = 24 * 7 * time.Hour // 7 days for token
	signInKey = "rhTC5@<HD<`(6!F%F8=cWw^`2zBnAxCr}kHP"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

func GenerateToken(userId, email string, active bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
		email,
		active,
	})

	return token.SignedString([]byte(signInKey))
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), _Cost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
