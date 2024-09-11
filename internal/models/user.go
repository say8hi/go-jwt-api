package models

type UserInDatabase struct {
	ID           string    `json:"id"`
	Username     string `json:"username"`
  Email        string `json:"email"`
	RefreshHash  *string `json:"refresh_hash,omitempty"`
	IP           string `json:"ip"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
  Email    string `json:"email"`
  IP       string `json:"ip"`
}

type UserUpdateRequest struct {
	Username    string `json:"username,omitempty"`
	IP          string `json:"ip,omitempty"`
	RefreshHash string `json:"refresh_hash,omitempty"`
}
