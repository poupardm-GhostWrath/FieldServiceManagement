// Configuration
const API_BASE = '/auth'

// Element References
const loginForm = document.getElementById('loginForm')
const emailInput = document.getElementById('email')
const passwordInput = document.getElementById('password')
const submitBtn = document.getElementById('submitBtn')
const apiError = document.getElementById('apiError')
const togglePassword = document.querySelector('.toggle-password')

// Password Visibility Toggle
togglePassword?.addEventListener('click', () => {
  const type = passwordInput.type === 'password' ? 'text' : 'password';
  passwordInput.type = type;
});

// Form Validation
function validateForm() {
  let isValid = true;

  // Clear previous errors
  clearErrors();

  // Email validation
  if (!emailInput.value.trim()) {
    showError('emailError', 'Email is required');
    isValid = false;
  } else if (!isValidEmail(emailInput.value)) {
    showError('emailError', 'Please enter a valid email address');
    isValid = false;
  }

  // Password validation
  if (!passwordInput.value) {
    showError('passwordError', 'Password is required');
    isValid = false;
  }

  return isValid;
}

function isValidEmail(email) {
  const re = /^[A-Za-z0-9]([A-Za-z0-9._%+-]*[A-Za-z0-9])?@[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?\.[A-Za-z]{2,}$/;
  return re.test(email);
}

function showError(elementId, message) {
  const el = document.getElementById(elementId);
  if (el) {
    el.textContent = message;
  }
}

function clearErrors() {
  document.querySelectorAll('.error-message').forEach(el => el.textContent = '');
}

// Show/Hide API Error
function showApiError(message) {
  apiError.querySelector('.alert-message').textContent = message;
  apiError.style.display = 'flex';
  setTimeout(() => { apiError.style.display = 'none'; }, 5000);
}

// Login Handler
async function handleLogin(e) {
  e.preventDefault();
  if (!validateForm()) return;

  // Disable button during request
  submitBtn.disabled = true;
  submitBtn.querySelector('.btn-text').style.display = 'none';
  submitBtn.querySelector('.btn-loader').style.display = 'inline';

  try {
    const response = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: emailInput.value.trim(),
        password: passwordInput.value,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Login failed');
    }

    // Store token securely
    localStorage.setItem('authToken', data.token);

    // Redirect based on role
    window.location.href = '/dashboard.html';
  
  } catch (error) {
    console.error('Login error:', error);
    showApiError(error.message || 'Unable to connect to server');
  } finally {
    submitBtn.disabled = false;
    submitBtn.querySelector('.btn-text').style.display = 'inline';
    submitBtn.querySelector('.btn-loader').style.display = 'none';
  }
}

// Event Listener
loginForm.addEventListener('submit', handleLogin);

// Clear error on input
[emailInput, passwordInput].forEach(input => {
  input.addEventListener('input', () => {
    clearErrors();
    apiError.style.display = 'none';
  });
});
