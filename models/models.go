package models

type Session struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Duration    int64  `json:"duration"`
	Ts          int64  `json:"Ts"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Ts       int64  `json:"Ts"`
}
