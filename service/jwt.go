package service

import (
	"advanced-webapp-project/utils"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type IJWTService interface {
	GenerateToken(userId string, email string) string
	ValidateToken(token string) (*jwt.Token, error)
	ExtractToken(token string) (map[string]any, bool)
}

type jwtCustomClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
	logger    *utils.Logger
}

func NewJWTService(logger *utils.Logger) *jwtService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "2022-19-11-19KTPM3-Advanced-Web-App",
		logger:    logger,
	}
}

func getSecretKey() string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "Golang"
	}
	return secretKey
}

func (svc *jwtService) GenerateToken(userId string, email string) string {
	claims := &jwtCustomClaims{
		userId,
		email,
		jwt.StandardClaims{
			Issuer:    svc.issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	encoded, err := token.SignedString([]byte(svc.secretKey))
	if err != nil {
		svc.logger.Error(err)
		os.Exit(1)
	}
	return encoded
}

func (svc *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (any, error) {
		if _, isValid := t_.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %+v", t_.Header["alg"])
		}
		return []byte(svc.secretKey), nil
	})
}

func (svc *jwtService) ExtractToken(tokenStr string) (map[string]any, bool) {
	token, _ := jwt.Parse(tokenStr, func(t_ *jwt.Token) (any, error) {
		return []byte(svc.secretKey), nil
	})
	return token.Claims.(jwt.MapClaims), true
}
