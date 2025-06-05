package models

type Message struct {
	ID        int64  `json:"id"`
	ChatID    int64  `json:"chat_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}
