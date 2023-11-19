package middlewares

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/khallihub/godoc/service"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""
		receivedToken := c.Query("token")
		if receivedToken != "" {
			tokenString = receivedToken
		} else {
			const BEARER_SCHEMA = "Bearer "
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || len(authHeader) < len(BEARER_SCHEMA) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			tokenString = authHeader[len(BEARER_SCHEMA):]
		}

		token, err := service.NewJWTService().ValidateToken(tokenString)
		print(token)
		print(err)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claims[Name]: ", claims["name"])
			log.Println("Claims[Admin]: ", claims["admin"])
			log.Println("Claims[Issuer]: ", claims["iss"])
			log.Println("Claims[IssuedAt]: ", claims["iat"])
			log.Println("Claims[ExpiresAt]: ", claims["exp"])
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
