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
	r.HandleFunc("/chats", handlers.ChatPageHandler).Methods("GET")
	protectedChat := r.PathPrefix("/api/chats").Subrouter()
	protectedChat.Use(handlers.AuthMiddleware)
	protectedChat.HandleFunc("", handlers.ChatsHandler).Methods("GET")
	protectedChat.HandleFunc("/{chatID}/messages", handlers.ChatMessagesHandler).Methods("GET")

	//api
	protectedAPI := r.PathPrefix("/api").Subrouter()
	protectedAPI.Use(handlers.AuthMiddleware)
	protectedAPI.HandleFunc("/chats", handlers.ChatsHandler).Methods("GET")
	protectedAPI.HandleFunc("/chats/{chatID}/messages", handlers.ChatMessagesHandler).Methods("GET")

	//websocket
	protectedWS := r.PathPrefix("/ws").Subrouter()
	protectedWS.Use(handlers.AuthMiddleware)
	protectedWS.HandleFunc("/{chatID}", handlers.WebSocketHandler)

	return r
}
