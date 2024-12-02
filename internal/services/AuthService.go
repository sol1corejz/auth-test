package services

import (
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sol1corejz/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"userID"`
	UserIP string `json:"userIP"`
}

const tokenExp = 30 * time.Minute
const secretKey = "supersecretkey"

var key = []byte("your-32-byte-long-secret-key!")

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
		UserID: userID,
		UserIP: userIP,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func generateRefreshToken() (string, string, error) {

	refreshToken := base64.URLEncoding.EncodeToString(key)

	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hashedRefreshToken), nil
}

func IsRefreshValid(refreshToken, hashedToken string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refreshToken))
	return err == nil
}

func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token is not valid")
	}

	if claims.UserID == "" {
		return "", errors.New("user ID is nil")
	}

	return claims.UserID, nil
}
