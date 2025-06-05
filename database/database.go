package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AlianZip/web-messenger/models"
	"github.com/AlianZip/web-messenger/utils"
	"github.com/go-playground/validator/v10"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var validate = validator.New()

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/storage/messenger.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to SQLite")

	//table for users
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            hash TEXT NOT NULL,
			timestamp INT NOT NULL,
			premission INT NOT NULL
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	//table for sessions
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS sessions (
        session_id TEXT PRIMARY KEY,
        user_id INTEGER NOT NULL,
        expires_at INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id)
    )`)
	if err != nil {
		log.Fatal(err)
	}

	// table for chats
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS chats (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	// table for messages
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        chat_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        timestamp INTEGER NOT NULL,
        FOREIGN KEY(chat_id) REFERENCES chats(id),
        FOREIGN KEY(user_id) REFERENCES users(id)
    )`)
	if err != nil {
		log.Fatal(err)
	}
}

// /
// /
// / USERS
// /
// addd new user
func CreateUser(user *models.User) error {
	_, err := DB.Exec(
		"INSERT INTO users (username, hash, timestamp, premission) VALUES (?, ?, ?, ?)",
		user.Username,
		user.Hash,
		time.Now().Unix(),
		0,
	)
	return err
}

// get user by it username
func GetUserByUsername(username string) (*models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Hash,
		&user.Timestamp,
		&user.Premission,
	)
	if err == sql.ErrNoRows {
		return &models.User{}, nil
	}
	return &user, err
}

// get user by it id
func GetUserByID(id int64) (*models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Hash,
		&user.Timestamp,
		&user.Premission,
	)
	if err == sql.ErrNoRows {
		return &models.User{}, nil
	}
	return &user, err
}

// /
// /
// / SESSIONS
// /
// create new session
func CreateSession(userID int64) (string, error) {
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()

	_, err = DB.Exec(
		"INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// search session by id
func GetSessionBySessionID(sessionID string) (*models.Session, error) {
	row := DB.QueryRow("SELECT session_id, user_id, expires_at FROM sessions WHERE session_id = ?", sessionID)
	var session models.Session
	err := row.Scan(&session.SessionID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt < time.Now().Unix() {
		DeleteSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// delete session
func DeleteSession(sessionID string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

// /
// /
// / CHATS
// /
// create new message
func CreateMessage(message *models.Message) error {
	_, err := DB.Exec(
		"INSERT INTO messages (chat_id, user_id, content, timestamp) VALUES (?, ?, ?, ?)",
		message.ChatID, message.UserID, message.Content, message.Timestamp,
	)
	return err
}

// get message by id
func GetMessagesByChatID(chatID int64) ([]*models.Message, error) {
	rows, err := DB.Query("SELECT id, chat_id, user_id, content, timestamp FROM messages WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message

	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func GetChats() ([]*models.Chat, error) {
	rows, err := DB.Query("SELECT id, name FROM chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*models.Chat

	for rows.Next() {
		var chat models.Chat
		err := rows.Scan(&chat.ID, &chat.Name)
		if err != nil {
			return nil, err
		}
		chats = append(chats, &chat)
	}

	return chats, nil
}
