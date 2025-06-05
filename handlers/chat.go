package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlianZip/web-messenger/database"
	"github.com/AlianZip/web-messenger/utils"
	"github.com/gorilla/mux"
)

func ChatPageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/chat.html")
}

func ChatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		chats, err := database.GetChats()
		if err != nil {
			http.Error(w, `{"error": "Failed to fetch chats"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chats)
		return
	}

	http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
}

func ChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatIDStr := vars["chatID"]
	chatID := utils.StringToInt64(chatIDStr)

	messages, err := database.GetMessagesByChatID(chatID)
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch messages"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
