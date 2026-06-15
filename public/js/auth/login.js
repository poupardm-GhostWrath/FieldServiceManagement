// public/js/auth/login.js
import { Session } from './session.js';

const ROLE_PRIORITY = ['admin', 'dispatcher', 'technician', 'customer'];

const DASHBOARD_URLS = {
    admin: '/dashboards/admin.html',
    dispatcher: '/dashboards/dispatcher.html',
    technician: '/dashboards/technician.html',
    customer: '/dashboards/customer.html'
};

document.addEventListener('DOMContentLoaded', function () {
    const loginForm = document.getElementById('loginForm');
    const submitBtn = document.getElementById('submitBtn');

    if (!loginForm) {
      console.error("Login form not found!");
      return;
    }

    // Redirect already-logged-in users
    if (Session.isAuthenticated()) {
        const roles = Session.getRoles();
        const url = getRedirectUrl(roles);
        if (url) window.location.href = url;
    }

    // Toggle password visibility
    const togglePassword = document.getElementById('togglePassword');
    const passwordInput = document.getElementById('password');
    if (togglePassword && passwordInput) {
        togglePassword.addEventListener('click', function () {
            const type = passwordInput.type === 'password' ? 'text' : 'password';
            passwordInput.type = type;
            this.textContent = type === 'password' ? 'Show' : 'Hide';
        });
    }

    // Form submission
    loginForm.addEventListener('submit', async function (e) {
        e.preventDefault();
        hideAlert();
        clearFieldErrors(['email', 'password']);

        const email = document.getElementById('email').value.trim();
        const password = document.getElementById('password').value;

        if (!email || !isEmailValid(email)) {
            showError('email', 'Please enter a valid email address');
            return;
        }
        if (!password || password.length < 1) {
            showError('password', 'Password is required');
            return;
        }

        submitBtn.disabled = true;
        submitBtn.textContent = 'Logging in...';

        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });

            if (response.ok) {
                const data = await response.json();

                // Use shared Session module
                Session.login(data.token, data.user);

                // Redirect based on roles
                const redirectUrl = getRedirectUrl(data.user.roles);
                if (redirectUrl) {
                    window.location.href = redirectUrl;
                } else {
                    throw new Error('No valid role found for your account');
                }
            } else {
                const errorData = await response.json().catch(() => ({ message: 'Invalid credentials' }));
                showAlert(errorData.message || 'Invalid credentials', 'error');
                submitBtn.disabled = false;
                submitBtn.textContent = 'Login';
            }
        } catch (error) {
            showAlert('Network error. Please try again.', 'error');
            submitBtn.disabled = false;
            submitBtn.textContent = 'Login';
        }
    });
});

function isEmailValid(email) {
  const emailRegex = /^[A-Za-z0-9]([A-Za-z0-9._%+-]*[A-Za-z0-9])?@[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?\.[A-Za-z]{2,}$/;
  return emailRegex.test(email);
}

function getRedirectUrl(roles) {
    if (!roles || !Array.isArray(roles)) return null;
    for (const role of ROLE_PRIORITY) {
        if (roles.includes(role)) return DASHBOARD_URLS[role];
    }
    return null;
}

function showError(fieldId, message) {
    const errorEl = document.getElementById(`error-${fieldId}`);
    const inputEl = document.getElementById(fieldId);
    
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.style.display = 'block';
    }
    if (inputEl) inputEl.classList.add('error');
}

function clearFieldErrors(fields) {
    fields.forEach(fieldId => {
        const errorEl = document.getElementById(`error-${fieldId}`);
        const inputEl = document.getElementById(fieldId);
        
        if (errorEl) errorEl.style.display = 'none';
        if (inputEl) inputEl.classList.remove('error');
    });
}

function showAlert(message, type) {
    const alertBox = document.getElementById('alertBox');
    const successAlert = document.getElementById('successAlert');
    const alertMsg = document.getElementById('alertMessage');
    const successMsg = document.getElementById('successMessage');
    
    if (type === 'error' && alertBox && alertMsg) {
        alertMsg.textContent = message;
        alertBox.style.display = 'block';
        setTimeout(hideAlert, 5000);
    } else if (type === 'success' && successAlert && successMsg) {
        successMsg.textContent = message;
        successAlert.style.display = 'block';
        setTimeout(hideAlert, 5000);
    }
}

function hideAlert() {
    const alertBox = document.getElementById('alertBox');
    const successAlert = document.getElementById('successAlert');
    if (alertBox) alertBox.style.display = 'none';
    if (successAlert) successAlert.style.display = 'none';
}
