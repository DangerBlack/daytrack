package models

import "time"

type SignIn struct {
	Email           string `json:"email"`
	Challenge       string `json:"challenge"`
	SignedChallenge string `json:"signed_challenge"`
}

type Token struct {
	Token   string    `json:"token"`
	ExpDate time.Time `json:"exp_date"`
}
