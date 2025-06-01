document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("registerForm");
    const passwordInput = document.getElementById("password");
    const confirmPasswordInput = document.getElementById("confirmPassword");
    const passwordError = document.getElementById("passwordError");

    const toggleButton = document.getElementById("themeToggle");
    const body = document.body;

    //theme
    toggleButton.addEventListener("click", function () {
        if (body.getAttribute("data-bs-theme") === "dark") {
            body.setAttribute("data-bs-theme", "light");
            toggleButton.textContent = "Темная тема";
        } else {
            body.setAttribute("data-bs-theme", "dark");
            toggleButton.textContent = "Светлая тема";
        }
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