package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AlianZip/web-messenger/database"
	"github.com/AlianZip/web-messenger/models"
	"github.com/AlianZip/web-messenger/utils"

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

	password, ok := data["password"]
	if !ok || password == "" {
		return fmt.Errorf("missing password")
	}

	user.Username = username
	user.Password = password

	return nil
}

// secure routes
func AuthMiddleware(next http.Handler) http.Handler {
	middlewareFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := utils.GetSessionCookie(r)

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
		existingUser, err := database.GetUserByUsername(user.Username)
		if existingUser.ID != 0 {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "User exists"}`)
			return
		}

		// hash password
		hash := utils.HashPassword(user.Password)
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

		var user models.User
		err := parseForm(r, &user)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "Bad request"}`)
			return
		}

		// check user existence
		existingUser, _ := database.GetUserByUsername(user.Username)
		if existingUser.ID == 0 {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "User not exist"}`)
			return
		}

		//check password
		fmt.Printf("%v \n", user.Password)
		hash := utils.HashPassword(user.Password)
		if hash != existingUser.Hash {
			fmt.Printf("%v : %v\n", hash, existingUser.Hash)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"success": false, "message": "Invalid pasword"}`)
			return
		}

		// set cookie
		sessionID, err := database.CreateSession(user.ID)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		fmt.Printf("set session cookie")
		utils.SetSessionCookie(w, sessionID, 86400*7) // 7 days
		fmt.Println(" succses")
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
