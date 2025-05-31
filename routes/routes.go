package routes

import (
	"github.com/AlianZip/web-messenger/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	//auth
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/delete-account", handlers.DeleteAccountHandler).Methods("POST")

	//chat

	return r
}
