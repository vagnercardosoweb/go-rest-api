package types

import "time"

type UserLoginInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	IpAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
}

type UserLoginOutput struct {
	AccessToken string    `json:"accessToken"`
	ExpiresIn   time.Time `json:"expiresIn"`
	TokenType   string    `json:"tokenType"`
}
