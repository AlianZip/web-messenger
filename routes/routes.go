package routes

import (
	"github.com/AlianZip/web-messenger/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	//home
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	//auth && delete account
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST", "GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST", "GET")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	// r.HandleFunc("/delete-account", handlers.DeleteAccountHandler).Methods("POST")

	//chat
	protected := r.PathPrefix("/chats").Subrouter()
	protected.Use(handlers.AuthMiddleware)
	protected.HandleFunc("/", handlers.ChatsHandler).Methods("POST", "GET")

	return r
}
