package models

import "time"

type SignIn struct {
	Email           string `json:"email" binding:"required" example:"test@dailytrack.it"`
	Challenge       string `json:"challenge" binding:"required" example:"a1b2c3d4e5f6g7h8i9j0"`
	SignedChallenge string `json:"signed_challenge" binding:"required" example:"a1b2c3d4e5f6g7h8i9j0"`
}

type Token struct {
	Token    string    `json:"token" example:"a1b2c3d4e5f6g7h8i9j0"`
	ExpDate  time.Time `json:"exp_date" example:"2021-01-01T00:00:00Z"`
	Username string    `json:"username" example:"default"`
}

type Challenge struct {
	Challenge string `json:"challenge" example:"a1b2c3d4e5f6g7h8i9j0"`
	Salt      string `json:"salt" example:"a1b2c3d4e5f6g7h8i9j0"`
}
