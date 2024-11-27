package services

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sol1corejz/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	userID string
	userIP string
}

const tokenExp = 30 * time.Minute
const secretKey = "supersecretkey"

func GenerateTokens(userID string, userIP string) (models.AuthTokens, error) {

	tokenString, err := generateAccessToken(userID, userIP)
	if err != nil {
		return models.AuthTokens{}, err
	}
	refreshToken, hashedRefreshToken, err := generateRefreshToken()
	if err != nil {
		return models.AuthTokens{}, err
	}

	return models.AuthTokens{
		AccessToken:        tokenString,
		RefreshToken:       refreshToken,
		HashedRefreshToken: hashedRefreshToken,
	}, nil
}

func generateAccessToken(userID string, userIP string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		userID: userID,
		userIP: userIP,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateRefreshToken() (string, string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", "", err
	}

	refreshToken := base64.URLEncoding.EncodeToString(randomBytes)

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hashedRefreshToken), nil
}
