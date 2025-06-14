

## Description

Ctrl+Enter Inc. is a lightweight web-based messaging application built using Go for the backend and HTML/CSS/JavaScript for the frontend.

### Main Features:
- User registration and login
- Chat list and selection
- Sending and receiving messages
- Edit and delete own messages
- Global and per-chat search
- Light and dark theme support
- Responsive design (Bootstrap 5)

---

## Technologies

- **Go** – backend server implementation
- **HTML/CSS/JS** – frontend interface
- **Bootstrap 5** – styling and responsiveness
- **Net/HTTP + Gorilla Mux** – routing

---

## Installation and Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/AlianZip/web-messenger.git
   cd web-message
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the server:

   ```bash
   go run main.go
   ```

4. Open in your browser:

   ```
   http://localhost:8080
   ```

---

## Project Structure

```
.
├── main.go
├── config/
│   └── config.go
├── handlers/
│   ├── auth.go
│   ├── chat.go
│   └── media.go
│   └── home.go
├── models/
│   ├── user.go
│   ├── message.go
│   └── chat.go
├── routes/
│   └── routes.go
├── templates/
│   ├── home.html
│   ├── login.html
│   └── chat.html
│   └── register.html
└── static/
    ├── css/
    └── js/
```

---

## TODO

- Implement WebSocket for real-time messaging
- Add media upload functionality
- Connect a database (SQLite)
- Implement JWT-based authentication
- Set up deployment and CI/CD pipeline

---

## License

MIT License – see [LICENSE](LICENSE) file for details.
