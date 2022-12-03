package middleware

import (
	"github.com/gin-contrib/cors"
)

func InitCors() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"http://localhost:7777", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Content-Length", "Origin"},
		AllowCredentials: true,
	}
}
