package models

type Session struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
	Ts          int    `json:"Ts"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Ts       int    `json:"Ts"`
}