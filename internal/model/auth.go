package model

type RegisterRequest struct {
	Login string `binding:"required" json:"login"`
	Pass  string `binding:"required" json:"password"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Login string `binding:"required" json:"login"`
	Pass  string `binding:"required" json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `binding:"required" json:"refresh_token"`
}

