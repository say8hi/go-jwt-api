package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/say8hi/go-jwt-api/internal/database"
	"github.com/say8hi/go-jwt-api/internal/models"
	"github.com/say8hi/go-jwt-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
  var requestUser models.CreateUserRequest
  err := json.NewDecoder(r.Body).Decode(&requestUser)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
  requestUser.IP = ip
  
  user, err := database.CreateUser(requestUser)
  if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
      http.Error(w, "This username is already taken.", http.StatusBadRequest)
      return
  } else if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  
  fmt.Println(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func CreateTokensPair(w http.ResponseWriter, r *http.Request) {
  userID := r.URL.Query().Get("user_id")
	ip := r.RemoteAddr
  tokenGUID := uuid.New().String()

	accessToken, err := utils.CreateAccessToken(userID, ip, tokenGUID)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  
  refreshToken, hashedRefreshToken, err := utils.CreateRefreshToken(tokenGUID)
  if err != nil {
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		return
	}
  
  userRequest := models.UserUpdateRequest{Username: "", RefreshHash: hashedRefreshToken, IP: ""}
  err = database.UpdateUser(userID, userRequest)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
    
  w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"access_token": "%s", "refresh_token": "%s"}`, accessToken, refreshToken)))
}

func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var request struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := jwt.ParseWithClaims(request.AccessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGN")), nil
	})
  if err != nil{
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		userID := claims.UserID
		user, err := database.GetUserByID(userID)
		if err == sql.ErrNoRows {
		  http.Error(w, "Wrong refresh token", http.StatusConflict)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		refreshTokenBytes, err := base64.StdEncoding.DecodeString(request.RefreshToken)
		if err != nil {
			http.Error(w, "Invalid refresh token format", http.StatusUnauthorized)
			return
		}

    refreshString := string(refreshTokenBytes)
    if refreshString[:36] != claims.TokenGUID{
      http.Error(w, "Wrong refresh/access token.", http.StatusBadRequest)
      return
    }
    
		if err := bcrypt.CompareHashAndPassword([]byte(*user.RefreshHash), refreshTokenBytes); err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		if user.IP != claims.IP {
			utils.SendEmailWarning(user.Email)
		}

    tokenGUID := uuid.New().String()
		newAccessToken, err := utils.CreateAccessToken(userID, r.RemoteAddr, tokenGUID)
		if err != nil {
			http.Error(w, "Error generating new access token", http.StatusInternalServerError)
			return
		}

		newRefreshToken, hashedNewRefreshToken, err := utils.CreateRefreshToken(tokenGUID)
		if err != nil {
			http.Error(w, "Error generating new refresh token", http.StatusInternalServerError)
			return
		}

    userRequest := models.UserUpdateRequest{Username: "", RefreshHash: hashedNewRefreshToken, IP: ""}
    err = database.UpdateUser(userID, userRequest)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"access_token": "%s", "refresh_token": "%s"}`, newAccessToken, newRefreshToken)))
	} else {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
	}
}
