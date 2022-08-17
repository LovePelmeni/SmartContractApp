package middlewares

import (
	"github.com/gin-gonic/gin"
)

func JwtAuthenticationMiddleware() gin.HandlerFunc {
	// Middleware for Validating Jwt Auth Tokens
	return func(context *gin.Context) {
	}
}
