document.addEventListener("DOMContentLoaded", function () {
    const toggleButton = document.getElementById("themeToggle");
    const form = document.getElementById("loginForm");
    const usernameInput = document.getElementById("username");
    const passwordInput = document.getElementById("password");
    const body = document.body;
    const savedTheme = localStorage.getItem("theme") || "light";
    
    
    
    body.setAttribute("data-bs-theme", savedTheme);
    toggleButton.textContent = savedTheme === "dark" ? "Светлая тема" : "Темная тема";

    toggleButton.addEventListener("click", function () {
        const currentTheme = body.getAttribute("data-bs-theme");
        const newTheme = currentTheme === "dark" ? "light" : "dark";
        body.setAttribute("data-bs-theme", newTheme);
        localStorage.setItem("theme", newTheme);
        toggleButton.textContent = newTheme === "dark" ? "Светлая тема" : "Темная тема";
    });


    //clear all errors
    function clearErrors() {
        document.querySelectorAll(".is-invalid").forEach(el => el.classList.remove("is-invalid"));
        document.querySelectorAll(".invalid-feedback").forEach(el => {
            el.textContent = "";
            el.style.display = "none";
        });
    }


    //show error from field
    function showFailedError(fieldId, message) {
        const field = document.getElementById(fieldId);
        if (field) {
            field.classList.add("is-invalid");
            const feedback = field.closest(".mb-3").querySelector(".invalid-feedback");
            if (feedback) {
                feedback.textContent = message;
                feedback.style.display = "block";
            };
        };
    };


    form.addEventListener("submit", async function (e) {
        e.preventDefault();
        
        clearErrors();

        const userData = {
            username: usernameInput.value.trim(),
            password: passwordInput.value.trim(),
        };


        try {
            const response = await fetch("/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(userData)
            });


            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            };


            try {
                result = await response.json();
            } catch (e) {
                console.log("Ошибка парсинга JSON:", e);
                console.log(response);
                alert("Ошибка: сервер вернул некорректный ответ");
                return false;
            };


            if (!result.success) {
                console.log(result);
                if (result.field) {
                    showFailedError(result.field, result.message);
                } else {
                    alert(result.message || "Ошибка входа");
                };
                return false;
            };

            window.location.href = result.redirect;
        } catch (error) {
            console.log("Ошибка сети:", error);
            alert("Не удалось подключиться к серверу");
            return false
        }


        
    });
});