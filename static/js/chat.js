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
    
    
});