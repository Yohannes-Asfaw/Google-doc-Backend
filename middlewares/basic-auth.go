package middlewares

import "github.com/gin-gonic/gin"

// BasicAuth takes as argument a map[string]string where
// the key is the user name and the value is the password.
func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"pragmatic": "reviews",
	})
}
