package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
  Email  string `json:"email"`
  TokenGUID  string `json:"token_guid"`
	jwt.RegisteredClaims
}
