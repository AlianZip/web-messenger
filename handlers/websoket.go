package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AlianZip/web-messenger/models"
	"github.com/AlianZip/web-messenger/utils"

	"github.com/AlianZip/web-messenger/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true //only for localhost or test
	},
}

// websoket client
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

// chat room
type ChatRoom struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var chatRooms = make(map[string]*ChatRoom)

func GetChatRoom(roomName string) *ChatRoom {
	if room, ok := chatRooms[roomName]; ok {
		return room
	}

	// create new chat
	room := &ChatRoom{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}

	chatRooms[roomName] = room

	go room.Run()

	return room
}

func (r *ChatRoom) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(r.Clients, client)
				}
			}
		}
	}
}

// websocket handler
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// chatid from url
	vars := mux.Vars(r)
	chatIDStr := vars["chatID"]
	chatID := utils.StringToInt64(chatIDStr)

	sessionID := utils.GetSessionCookie(r)
	if sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := database.GetSessionBySessionID(sessionID)
	if err != nil || session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.UserID

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to open WebSocket: %v\n", err)
		http.Error(w, "Failed to open WebSocket", http.StatusInternalServerError)
		return
	}

	client := &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	room := GetChatRoom(chatIDStr)
	room.Register <- client

	go sendHistory(client, chatID)

	go client.WritePump()
	go client.ReadPump(userID, chatIDStr)
}

func (c *Client) ReadPump(userID int64, chatID string) {
	defer func() {
		room := GetChatRoom(chatID)
		room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		message := models.Message{
			ChatID:    utils.StringToInt64(chatID),
			UserID:    userID,
			Content:   string(msg),
			Timestamp: time.Now().Unix(),
		}

		err = database.CreateMessage(&message)
		if err != nil {
			log.Printf("Failed to save message: %v\n", err)
			continue
		}

		room := GetChatRoom(chatID)
		room.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	for {
		msg, ok := <-c.Send
		if !ok {
			return
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func sendHistory(client *Client, chatID int64) {
	messages, err := database.GetMessagesByChatID(chatID)
	if err != nil {
		log.Printf("Failed to get message history: %v\n", err)
		return
	}

	for _, msg := range messages {
		client.Send <- []byte(fmt.Sprintf("[%d] %s", msg.UserID, msg.Content))
	}
}
