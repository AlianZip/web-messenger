package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

// generate new seesionid of 32 chars
func GenerateSessionID() (string, error) {
	b := make([]byte, 24) // 24 байта = 32 символа в base64
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// set cookie with sessionid
func SetSessionCookie(w http.ResponseWriter, sessionID string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //only for localhost
		SameSite: http.SameSiteStrictMode,
	})
}

// read session from coolie
func GetSessionCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// delete cookie
func DeleteSessionCookie(w http.ResponseWriter) {
	SetSessionCookie(w, "", -1)
}
