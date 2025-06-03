package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AlianZip/web-messenger/database"
	"github.com/AlianZip/web-messenger/models"
	"github.com/AlianZip/web-messenger/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// parse from request to type User from github.com/AlianZip/web-messenger/models(user.go)
func parseForm(r *http.Request, user *models.User) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("invalid content-type, expected application/json")
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	// check field exists
	username, ok := data["username"]
	if !ok || username == "" {
		return fmt.Errorf("missing username")
	}

	email, ok := data["email"]
	if !ok || email == "" {
		return fmt.Errorf("missing email")
	}

	password, ok := data["password"]
	if !ok || password == "" {
		return fmt.Errorf("missing password")
	}

	user.Username = username
	user.Email = email
	user.Password = password

	return nil
}

// secure routes
func AuthMiddleware(next http.Handler) http.Handler {
	middlewareFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdStr := utils.GetCookie(r, "user_id")
		sessionID := utils.GetSessionCookie(r)

		if userIdStr == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, _ := database.GetUserByID(userId)
		if user.ID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if sessionID == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, err := database.GetSessionBySessionID(sessionID)
		if err != nil || session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})

	return middlewareFunc
}

// register new user
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		var user models.User
		err := parseForm(r, &user)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "Bad request"}`)
			return
		}

		// check valid
		if err := validate.Struct(user); err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "Invalid structure"}`)
			return
		}

		// check user existence
		existingUser, _ := database.GetUserByUsername(user.Username)
		if existingUser.ID != 0 {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "User existence"}`)
			return
		}

		existingEmail, _ := database.GetUserByEmail(user.Email)
		if existingEmail.ID != 0 {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "User existence"}`)
			return
		}

		// hash password
		hash, err := utils.HashPassword(user.Password)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "Invalid hashing"}`)
			return
		}

		// save user
		user.Hash = hash
		user.Timestamp = time.Now().Unix()
		if err := database.CreateUser(&user); err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "Invalid add user"}`)
			return
		}

		// go to login
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": true, "redirect": "/login"}`)
		return
	}

	http.ServeFile(w, r, "templates/register.html")
}

// login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := database.GetUserByUsername(username)
		if err != nil || user.ID == 0 {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// set cookie
		sessionID, err := database.CreateSession(user.ID)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		utils.SetSessionCookie(w, sessionID, 86400*7) // 7 days

		// go to chats
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": true, "redirect": "/chats"}`)
		return
	}

	http.ServeFile(w, r, "templates/login.html")

}

// just logout lol
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := utils.GetSessionCookie(r)
	if sessionID != "" {
		database.DeleteSession(sessionID)
		utils.DeleteSessionCookie(w)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
