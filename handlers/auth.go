package handlers

import (
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "Login logic here")
		return
	}
	http.ServeFile(w, r, "./templates/login.html")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "Register logic here")
		return
	}
	http.ServeFile(w, r, "./templates/register.html")
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	return
}
