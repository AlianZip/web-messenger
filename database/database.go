package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AlianZip/web-messenger/models"
	"github.com/AlianZip/web-messenger/utils"

	_ "github.com/mattn/go-sqlite3"
)

var DB_USERS *sql.DB
var DB_CHATS *sql.DB

func InitDB() {
	var err error
	DB_USERS, err = sql.Open("sqlite3", "./database/storage/users.db")
	if err != nil {
		log.Fatal(err)
	}

	DB_CHATS, err = sql.Open("sqlite3", "./database/storage/messenger.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to SQLite")

	//table for users
	_, err = DB_USERS.Exec(`
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
	_, err = DB_USERS.Exec(`CREATE TABLE IF NOT EXISTS sessions (
        session_id TEXT PRIMARY KEY,
        user_id INTEGER NOT NULL,
        expires_at INTEGER NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id)
    )`)
	if err != nil {
		log.Fatal(err)
	}

	// table for chats
	_, err = DB_CHATS.Exec(`CREATE TABLE IF NOT EXISTS chats (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	// table for messages
	_, err = DB_CHATS.Exec(`CREATE TABLE IF NOT EXISTS messages (
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
	_, err := DB_USERS.Exec(
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
	row := DB_USERS.QueryRow("SELECT * FROM users WHERE username = ?", username)
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
	row := DB_USERS.QueryRow("SELECT * FROM users WHERE id = ?", id)
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

	_, err = DB_USERS.Exec(
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
	row := DB_USERS.QueryRow("SELECT session_id, user_id, expires_at FROM sessions WHERE session_id = ?", sessionID)
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
	_, err := DB_USERS.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}
