package middleware

import (
	"advanced-webapp-project/helper"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizeJWT(svc service.IJWTService, logger *helper.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "no token found!"})
			return
		}

		token, err := svc.ValidateToken(header)
		if !token.Valid {
			logger.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"message": "invalid token!"})
			return
		}
		c.Next()
	}
}
