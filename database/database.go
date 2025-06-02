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

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./messenger.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to SQLite")

	//table for users
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            email TEXT NOT NULL UNIQUE,
            hash TEXT NOT NULL,
			timestamp INT NOT NULL,
			sessionID TEXT NOT NULL
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
}

// /
// /
// / USERS
// /
// addd new user
func CreateUser(user *models.User) error {
	_, err := DB.Exec(
		"INSERT INTO users (username, email, hash, timestamp) VALUES (?, ?, ?, ?)",
		user.Username,
		user.Email,
		user.Hash,
		time.Now().Unix(),
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
		&user.Email,
		&user.Hash,
		&user.Timestamp,
	)
	if err == sql.ErrNoRows {
		return &models.User{}, nil
	}
	return &user, err
}

// get user by it email
func GetUserByEmail(email string) (*models.User, error) {
	row := DB.QueryRow("SELECT * FROM users WHERE email = ?", email)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Hash,
		&user.Timestamp,
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
		&user.Email,
		&user.Hash,
		&user.Timestamp,
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
