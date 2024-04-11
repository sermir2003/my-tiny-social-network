package models

type Credentials struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
