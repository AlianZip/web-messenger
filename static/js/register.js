document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("registerForm");
    const passwordInput = document.getElementById("password");
    const confirmPasswordInput = document.getElementById("confirmPassword");
    const passwordError = document.getElementById("passwordError");
    const toggleButton = document.getElementById("themeToggle");
    const body = document.body;
    const savedTheme = localStorage.getItem("theme") || "light";
    
    
    //theme
    body.setAttribute("data-bs-theme", savedTheme);
    toggleButton.textContent = savedTheme === "dark" ? "Светлая тема" : "Темная тема";

    toggleButton.addEventListener("click", function () {
        const currentTheme = body.getAttribute("data-bs-theme");
        const newTheme = currentTheme === "dark" ? "light" : "dark";
        body.setAttribute("data-bs-theme", newTheme);
        localStorage.setItem("theme", newTheme);
        toggleButton.textContent = newTheme === "dark" ? "Светлая тема" : "Темная тема";
    });


    //check password
    form.addEventListener("submit", function (e) {
        if (passwordInput.value !== confirmPasswordInput.value) {
            e.preventDefault();
            confirmPasswordInput.classList.add("is-invalid");
            passwordError.style.display = "block";
        } else {
            confirmPasswordInput.classList.remove("is-invalid");
            passwordError.style.display = "none";
        }
    });
});