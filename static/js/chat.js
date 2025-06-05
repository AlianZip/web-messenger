document.addEventListener("DOMContentLoaded", function () {
    const profileButton = document.getElementById("profileButton");
    const profileModal = new bootstrap.Modal(document.getElementById('profileModal'));
    const sendMessageButton = document.getElementById('sendMessageButton');
    const profileUsername = document.getElementById('profileUsername');
    const profileEmail = document.getElementById('profileEmail');
    const chatList = document.getElementById("chatList");
    const chatArea = document.getElementById("chatArea");
    const messageForm = document.getElementById("messageForm");
    const messageInput = document.getElementById("messageInput");
    const globalSearchInput = document.getElementById("globalSearchInput");
    const chatSearchInput = document.getElementById("chatSearchInput");
    const toggleButton = document.getElementById("themeToggle");
    const savedTheme = localStorage.getItem("theme") || "light";
    const body = document.body;

    let socket = null;
    let currentChatId = null;

    //theme
    body.setAttribute("data-bs-theme", savedTheme);
    toggleButton.textContent = savedTheme === "dark" ? "☀️" : "🌙";

    toggleButton.addEventListener("click", function () {
        const currentTheme = body.getAttribute("data-bs-theme");
        const newTheme = currentTheme === "dark" ? "light" : "dark";
        body.setAttribute("data-bs-theme", newTheme);
        localStorage.setItem("theme", newTheme);
        toggleButton.textContent = newTheme === "dark" ? "☀️" : "🌙";
    });



    function connectWebSocket(chatId) {
        if (socket && socket.readyState !== WebSocket.CLOSED) {
            socket.close();
        }

        socket = new WebSocket(`ws://localhost:8080/ws/${chatId}`);

        socket.onopen = function () {
            console.log("Connected to WebSocket server");
        };

        socket.onmessage = function (event) {
            const message = JSON.parse(event.data);
            appendMessage(message);
        };

        socket.onerror = function (error) {
            console.error("WebSocket error:", error);
        };

        socket.onclose = function () {
            console.log("WebSocket connection closed");
        };
    }


    function appendMessage(message) {
        const messageElement = document.createElement("div");
        messageElement.classList.add("message");
        messageElement.innerHTML = `<strong>${message.username}</strong>: ${message.content}`;
        chatArea.appendChild(messageElement);
        chatArea.scrollTop = chatArea.scrollHeight;
    }


    chatList.addEventListener("click", function (event) {
        if (event.target.tagName === "LI") {
            currentChatId = event.target.dataset.chatId;
            connectWebSocket(currentChatId);
            chatArea.innerHTML = "";
            fetchChatHistory(currentChatId);
        }
    });


    async function fetchChatHistory(chatId) {
        try {
            const response = await fetch(`/api/chats/${chatId}/messages`);
            const messages = await response.json();
            messages.forEach(message => appendMessage(message));
        } catch (error) {
            console.error("Error fetching chat history:", error);
        }
    }


    fetchChats();
    

    async function fetchChats() {
        try {
            const response = await fetch("/api/chats");
            const chats = await response.json();
            renderChats(chats);
        } catch (error) {
            console.error("Error fetching chats:", error);
        }
    }

    function renderChats(chats) {
        chatList.innerHTML = "";
        chats.forEach(chat => {
            const li = document.createElement("li");
            li.className = "list-group-item list-group-item-action";
            li.dataset.chatId = chat.id;
            li.textContent = chat.name;
            chatList.appendChild(li);
        });
    }
});