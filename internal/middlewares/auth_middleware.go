package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/excel_vue_go/internal/repositories"
	"github.com/rustingoff/excel_vue_go/packages/token"
	"log"
	"net/http"
	"strings"
)

func CheckToken(userRepo repositories.UserRepository) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		value, isCookie := ctx.Request.Header["Authorization"]

		if !isCookie {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("user is not authorized"))
			return
		}

		headerValue, err := stripBearerPrefixFromToken(value[0])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
			return
		}

		claims, err := token.ParseToken(headerValue)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
			return
		}

		userId := claims.UserId

		user, err := userRepo.GetUserById(userId)
		if err != nil {
			log.Println("[ERR]: user not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "user not found")
			return
		}

		if user.Email == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "invalid user")
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

//Bearer test_token12313131313131
func stripBearerPrefixFromToken(token string) (string, error) {

	bearer := "BEARER"

	if len(token) > len(bearer) && strings.ToUpper(token[0:len(bearer)]) == bearer {
		return token[len(bearer)+1:], nil
	}

	return token, nil
}
