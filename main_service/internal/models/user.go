package models

type User struct {
	Uuid      string `form:"uuid"`
	Username  string `username:"uuid"`
	Image_url string
}

type Registration struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Id struct {
	Ids string `json:"id"`
}

type RegistrResponse struct {
	Token string `json:"token"`
	Id    Id     `json:"user"`
}

type AuthResponse struct {
	Token string
	Id    string
}
