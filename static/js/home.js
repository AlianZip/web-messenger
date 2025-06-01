document.addEventListener("DOMContentLoaded", function () {
    const toggleButton = document.getElementById("themeToggle");
    const body = document.body;

    toggleButton.addEventListener("click", function () {
        if (body.getAttribute("data-bs-theme") === "dark") {
            body.setAttribute("data-bs-theme", "light");
            toggleButton.textContent = "Темная тема";
        } else {
            body.setAttribute("data-bs-theme", "dark");
            toggleButton.textContent = "Светлая тема";
        }
    });
});