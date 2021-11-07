package models

type User struct {
	ID string `json:"-"`

	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

type SingIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
