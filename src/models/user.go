package models

type User struct {
	ID        int64  `json:"id" example:"1"`
	Username  string `json:"username" example:"default"`
	Email     string `json:"email" example:"test@dailytrack.it"`
	PublicKey string `json:"public_key" example:"a1b2c3d4e5f6g7h8i9j0"`
	Salt      string `json:"salt" example:"a1b2c3d4e5f6g7h8i9j0"`
}
