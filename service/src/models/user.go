package models

type User struct {
	ID        int64  `json:"id" example:"1"`
	Username  string `json:"username" example:"default"`
	Email     string `json:"email" example:"test@dailytrack.it"`
	PublicKey string `json:"public_key" example:"a1b2c3d4e5f6g7h8i9j0"`
	Salt      string `json:"salt" example:"a1b2c3d4e5f6g7h8i9j0"`
}

type CreateUser struct {
	Username  string `json:"username" binding:"required,gt=2,lt=256" example:"default"`
	Email     string `json:"email" binding:"required,email" example:"test@dailytrack.it"`
	PublicKey string `json:"public_key" binding:"required" example:"a1b2c3d4e5f6g7h8i9j0"`
	Salt      string `json:"salt" binding:"required" example:"a1b2c3d4e5f6g7h8i9j0"`
}

type MagicLink struct {
	Email string `json:"email" binding:"required,email" example:"test@dailytrack.it"`
}
