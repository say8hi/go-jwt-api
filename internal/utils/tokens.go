package utils

import (
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/say8hi/go-jwt-api/internal/models"
)


func CreateAccessToken(userID, ip, tokenGUID string) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &models.Claims{
		UserID: userID,
		IP:     ip,
    TokenGUID: tokenGUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGN")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateRefreshToken(tokenGUID string) (string, string, error) {
  secondsStr := strconv.FormatInt(time.Now().Unix(), 10)
  refreshToken := tokenGUID + secondsStr
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	encodedToken := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return encodedToken, string(hashedToken), nil
}
