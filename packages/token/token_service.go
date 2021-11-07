package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
)

func ParseToken(accessToken string) (tokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("[ERR]: failed to parse token")
			return nil, errors.New("invalid token")
		}

		return []byte(signInKey), nil
	})

	if err != nil {
		log.Println("[ERR]: provided invalid token, ", err.Error())
		return tokenClaims{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		log.Println("[ERR]: token not are type of token claims")
		return tokenClaims{}, errors.New("token claims are not of type *tokenClaims")
	}

	return *claims, nil
}
