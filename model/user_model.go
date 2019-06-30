package model

// User model structure
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	UserName  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Country   string `json:"country"`
}
