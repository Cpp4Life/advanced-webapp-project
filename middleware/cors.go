package middleware

import (
	"github.com/gin-contrib/cors"
	"time"
)

func InitCors() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Content-Length", "Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
