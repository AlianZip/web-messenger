document.addEventListener("DOMContentLoaded", function () {
    const toggleButton = document.getElementById("themeToggle");
    const form = document.getElementById("loginForm");
    const nameInput = document.getElementById("username");
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


    form.addEventListener("submit", function (e) {
        e.defaultPrevented();
        
        clearErrors();

        
    });
});