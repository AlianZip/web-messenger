document.addEventListener("DOMContentLoaded", function () {
    const chatList = document.getElementById("chatList");
    const chatArea = document.getElementById("chatArea");
    const messageForm = document.getElementById("messageForm");
    const messageInput = document.getElementById("messageInput");
    const globalSearchInput = document.getElementById("globalSearchInput");
    const chatSearchInput = document.getElementById("chatSearchInput");

    let currentChatId = null;

    //theme
    const toggleButton = document.getElementById("themeToggle");
    const body = document.body;

    toggleButton.addEventListener("click", function () {
        if (body.getAttribute("data-bs-theme") === "dark") {
            body.setAttribute("data-bs-theme", "light");
            toggleButton.textContent = "🌙";
        } else {
            body.setAttribute("data-bs-theme", "dark");
            toggleButton.textContent = "☀️";
        }
    });

    //api: get chats
    async function fetchChats() {
        try {
            const res = await fetch("/api/chats");
            const chats = await res.json();
            renderChats(chats);
        } catch (err) {
            console.error("Error fetching chats:", err);
        }
    }

    function renderChats(chats) {
        chatList.innerHTML = "";
        chats.forEach(chat => {
            const li = document.createElement("li");
            li.className = "list-group-item list-group-item-action chat-item";
            li.dataset.chatId = chat.id;
            li.innerHTML = `
                <strong>${chat.name}</strong><br>
                <small class="text-muted">${chat.lastMessage || 'No messages yet'}</small>
            `;
            li.addEventListener("click", () => loadChatMessages(chat.id));
            chatList.appendChild(li);
        });
    }

    //api: get message
    async function loadChatMessages(chatId) {
        currentChatId = chatId;
        chatArea.innerHTML = "<div>Loading messages...</div>";

        try {
            const res = await fetch(`/api/chats/${chatId}/messages`);
            const messages = await res.json();
            renderMessages(messages);
        } catch (err) {
            chatArea.innerHTML = "<div class='text-danger'>Failed to load messages</div>";
        }
    }

    function renderMessages(messages) {
        chatArea.innerHTML = "";

        if (!messages.length) {
            chatArea.innerHTML = "<div class='text-muted'>This chat is empty.</div>";
            return;
        }

        messages.forEach(msg => {
            const div = document.createElement("div");
            div.className = `message ${msg.isMine ? "own" : "others"}`;
            div.dataset.msgId = msg.id;
            div.innerHTML = `
                <div>${msg.text}</div>
                <small class="text-muted">${new Date(msg.time).toLocaleTimeString()}</small>
                ${msg.isMine ? `
                <div class="mt-1">
                    <button class="btn btn-sm btn-outline-primary edit-btn">Edit</button>
                    <button class="btn btn-sm btn-outline-danger delete-btn">Delete</button>
                </div>` : ""}
            `;
            chatArea.appendChild(div);
        });

        attachMessageActions();
        chatArea.scrollTop = chatArea.scrollHeight;
    }

    function attachMessageActions() {
        document.querySelectorAll(".edit-btn").forEach(btn => {
            btn.addEventListener("click", function () {
                const msgDiv = this.closest(".message");
                const msgId = msgDiv.dataset.msgId;
                const originalText = msgDiv.querySelector("div").textContent;

                const modal = new bootstrap.Modal(document.getElementById("editMessageModal"));
                const textarea = document.getElementById("editMessageText");
                const saveBtn = document.getElementById("saveEditButton");

                textarea.value = originalText;
                modal.show();

                saveBtn.onclick = async () => {
                    const newText = textarea.value.trim();
                    if (!newText) return;

                    //patch request to api
                    try {
                        const res = await fetch(`/api/messages/${msgId}`, {
                            method: "PATCH",
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({ text: newText })
                        });

                        if (res.ok) {
                            msgDiv.querySelector("div").textContent = newText;
                            modal.hide();
                        }
                    } catch (err) {
                        alert("Failed to update message");
                    }
                };
            });
        });

        document.querySelectorAll(".delete-btn").forEach(btn => {
            btn.addEventListener("click", function () {
                const msgDiv = this.closest(".message");
                const msgId = msgDiv.dataset.msgId;

                if (confirm("Delete this message?")) {
                    //delete msg
                    fetch(`/api/messages/${msgId}`, {
                        method: "DELETE"
                    }).then(() => {
                        msgDiv.remove();
                    });
                }
            });
        });
    }

    //send msg
    messageForm.addEventListener("submit", async function (e) {
        e.preventDefault();
        const text = messageInput.value.trim();
        if (!text || !currentChatId) return;

        const newMsg = {
            text,
            time: new Date().toISOString(),
            isMine: true
        };

        //local add
        renderMessages([...Array.from(chatArea.getElementsByClassName("message")).map(el => ({
            text: el.querySelector("div").textContent,
            time: el.querySelector("small").textContent,
            isMine: el.classList.contains("own")
        })), newMsg]);

        messageInput.value = "";

        //send to server
        try {
            await fetch("/api/messages", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ chatId: currentChatId, text })
            });
        } catch (err) {
            alert("Failed to send message");
        }
    });

    //search
    globalSearchInput.addEventListener("input", function () {
        const term = this.value.toLowerCase();
        document.querySelectorAll(".chat-item").forEach(el => {
            const name = el.querySelector("strong").textContent.toLowerCase();
            el.style.display = name.includes(term) ? "block" : "none";
        });
    });

    chatSearchInput.addEventListener("input", function () {
        const term = this.value.toLowerCase();
        document.querySelectorAll("#chatArea .message").forEach(el => {
            const text = el.querySelector("div").textContent.toLowerCase();
            el.style.display = text.includes(term) ? "block" : "none";
        });
    });

    //init
    fetchChats();
});