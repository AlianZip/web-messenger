package models

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Hash      string `json:"-"`
	Timestamp int64  `json:"-"`
	Password  string `json:"-"`
}

type Session struct {
	SessionID string
	UserID    int64
	ExpiresAt int64
}
