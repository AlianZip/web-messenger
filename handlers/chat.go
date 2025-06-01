package handlers

import (
	"net/http"
)

func ChatsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/chat.html")
}
