package middleware

import (
	"advanced-webapp-project/helper"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func AuthorizeJWT(svc service.IJWTService, log helper.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "no token found!"})
			return
		}

		token, err := svc.ValidateToken(header)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Info(claims)
		} else {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "invalid token!"})
		}
	}
}
