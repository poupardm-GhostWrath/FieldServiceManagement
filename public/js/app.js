// app.js - Simple DOM Interaction
document.addEventListener('DOMContentLoaded', function() {
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');

    // Login button clicks
    loginBtn.addEventListener('click', function() {
        window.location.href = 'login.html';
    });

    // Register button clicks
    registerBtn.addEventListener('click', function() {
        window.location.href = 'register.html';
    });
});
